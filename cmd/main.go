package main

import (
	"fmt"
	"log"
	"net/http"
	"social-network/pkg/db"
	"social-network/pkg/handlers"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/kelseyhightower/envconfig"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tinrab/retry"
)

type Config struct {
	MysqlHost     string `envconfig:"MYSQL_HOST"`
	MysqlDB       string `envconfig:"MYSQL_DB"`
	MysqlUser     string `envconfig:"MYSQL_USER"`
	MysqlPassword string `envconfig:"MYSQL_PASSWORD"`
}

var (
	storeKey = []byte("store-key")
)

func newRouter(sessionsStore *sessions.CookieStore) (router *mux.Router) {
	router = mux.NewRouter()
	handler := &handlers.Handler{SessionsStore: sessionsStore}

	router.HandleFunc("/", handler.IndexHandler).Methods("GET")
	router.HandleFunc("/logout", handler.LogoutHandler).
		Methods("POST")
	router.HandleFunc("/signup", handler.SignupHandler).
		Methods("GET")
	router.HandleFunc("/signup", handler.SignupPostHandler).
		Methods("POST")
	router.HandleFunc("/login", handler.LoginHandler).
		Methods("GET")
	router.HandleFunc("/login", handler.LoginPostHandler).
		Methods("POST")
	router.HandleFunc("/users", handler.UserListHandler).
		Methods("GET")
	router.HandleFunc("/friends", handler.FriendsListHandler).
		Methods("GET")
	router.HandleFunc("/friend-requests", handler.FriendRequestListHandler).
		Methods("GET")
	router.HandleFunc("/friend-request", handler.FriendRequestPostHandler).
		Methods("POST")
	router.HandleFunc("/friend-request-accept", handler.FriendRequestAcceptHandler).
		Methods("POST")
	router.HandleFunc("/friend-request-decline", handler.FriendRequestDeclineHandler).
		Methods("POST")
	router.HandleFunc("/profile", handler.MyProfileHandler).
		Methods("GET")
	router.HandleFunc("/profile", handler.MyProfileUpdateHandler).
		Methods("PUT")
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	return
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	/*cfg.MysqlHost = "localhost:3306"
	cfg.MysqlUser = "user"
	cfg.MysqlPassword = "123456"
	cfg.MysqlDB = "app"*/

	retry.ForeverSleep(2*time.Second, func(attempt int) error {
		addr := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.MysqlUser, cfg.MysqlPassword, cfg.MysqlHost, cfg.MysqlDB)
		repo, err := db.NewMysql(addr)

		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	defer db.Close()

	// Run HTTP server

	store := sessions.NewCookieStore(storeKey)
	router := newRouter(store)
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}

}
