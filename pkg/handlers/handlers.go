package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"social-network/pkg/auth"
	"social-network/pkg/db"
	"social-network/pkg/schema"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	tokenService "social-network/pkg/token"

	"github.com/gorilla/sessions"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	SessionsStore *sessions.CookieStore
}

func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := h.SessionsStore.Get(r, "session")

	user, err := auth.CheckSession(session)
	if err != nil || user.ID == "" {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}

	http.Redirect(w, r, "/friends", http.StatusMovedPermanently)
}

func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("./template", "layout.html")
	fp := filepath.Join("./template", "signup.html")

	tmpl, _ := template.ParseFiles(lp, fp)
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func (h *Handler) SignupPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {

	}

	ctx := r.Context()

	email := template.HTMLEscapeString(strings.ToLower(r.FormValue("email")))
	if len(email) < 1 || len(email) > 140 {
		log.Println("Create user invalid email:", email)
	}

	name := template.HTMLEscapeString(r.FormValue("name"))
	if len(name) < 1 || len(name) > 140 {
		log.Println("Create user invalid name:", name)
	}

	surname := template.HTMLEscapeString(r.FormValue("surname"))
	if len(surname) < 1 || len(surname) > 140 {
		log.Println("Create user invalid surname:", surname)
	}

	city := template.HTMLEscapeString(r.FormValue("city"))
	if len(city) < 1 || len(city) > 140 {
		log.Println("Create user invalid city:", city)
	}

	password := template.HTMLEscapeString(r.FormValue("password"))
	if len(password) < 1 || len(password) > 140 {
	}

	sex := template.HTMLEscapeString(r.FormValue("sex"))
	interests := template.HTMLEscapeString(r.FormValue("interests"))

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Create user create hash password error:", err)
	}

	_, err = db.GetUserByEmail(ctx, email)
	if err == nil {
		log.Println("Create user error already exist:", err)
		http.Redirect(w, r, "/signup", http.StatusMovedPermanently)
		return
	}

	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		log.Println("Create user ksuid error:", err)
	}

	user := schema.User{
		ID:        id.String(),
		Email:     email,
		Password:  string(hashedPass),
		City:      city,
		Name:      name,
		Surname:   surname,
		Sex:       sex,
		Interests: interests,
		CreatedAt: createdAt,
	}
	if err := db.InsertUser(ctx, user); err != nil {
		log.Println("Create user error:", err)
	}

	token, err := tokenService.Encode(&user, false)
	if err != nil {

	}
	session, _ := h.SessionsStore.Get(r, "session")
	session.Values["token"] = token
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("./template", "layout.html")
	fp := filepath.Join("./template", "signin.html")

	tmpl, _ := template.ParseFiles(lp, fp)
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func (h *Handler) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {

	}

	ctx := r.Context()

	email := template.HTMLEscapeString(strings.ToLower(r.FormValue("email")))
	if len(email) < 1 || len(email) > 140 {
		log.Println("Login user invalid email:", email)
	}

	password := template.HTMLEscapeString(r.FormValue("password"))
	if len(password) < 1 || len(password) > 140 {
		log.Println("Login user invalid password:", password)
	}

	user, err := db.GetUserByEmail(ctx, email)
	if err != nil {
		log.Println("Login user error already exist:", err)
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Println("Login user error cant compare:", err)
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		return
	}

	token, err := tokenService.Encode(user, false)
	if err != nil {
		log.Println("Login user error cant encode:", err)
	}
	session, _ := h.SessionsStore.Get(r, "session")
	session.Values["token"] = token
	session.Save(r, w)

	http.Redirect(w, r, "/friends", http.StatusMovedPermanently)
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := h.SessionsStore.Get(r, "session")
	session.Values["token"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (h *Handler) UserListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, _ := h.SessionsStore.Get(r, "session")

	user, err := auth.CheckSession(session)
	if err != nil || user.ID == "" {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}

	lp := filepath.Join("./template", "layout.html")
	fp := filepath.Join("./template", "user_list.html")

	offsetQuery, ok := r.URL.Query()["offset"]

	offset := 0
	if ok {
		offset, _ = strconv.Atoi(offsetQuery[0])
	}

	limitQuery, ok := r.URL.Query()["limit"]

	limit := 50
	if ok {
		limit, _ = strconv.Atoi(limitQuery[0])
	}

	users, err := db.ListUsers(ctx, user.ID, offset, limit)
	if err != nil {
		log.Println("Get users error", err)
	}

	friends, err := db.ListFriends(ctx, user.ID)

	for uId, user := range users {
		i := sort.Search(len(friends), func(i int) bool { return user.ID == friends[i] })
		if i < len(friends) && friends[i] == user.ID {
			users[uId].IsMineFriend = true
		}
	}

	tmpl, _ := template.ParseFiles(lp, fp)

	data := struct {
		IsAuth bool
		Users  []schema.User
		User   *schema.User
	}{true, users, user}

	tmpl.ExecuteTemplate(w, "layout", data)
}

func (h *Handler) FriendsListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, _ := h.SessionsStore.Get(r, "session")

	user, err := auth.CheckSession(session)
	if err != nil || user.ID == "" {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}

	friends, err := db.ListFriends(ctx, user.ID)

	offsetQuery, ok := r.URL.Query()["offset"]

	offset := 0
	if ok {
		offset, _ = strconv.Atoi(offsetQuery[0])
	}

	limitQuery, ok := r.URL.Query()["limit"]

	limit := 50
	if ok {
		limit, _ = strconv.Atoi(limitQuery[0])
	}

	users, err := db.GetUsersByIDS(ctx, friends, offset, limit)

	lp := filepath.Join("./template", "layout.html")
	fp := filepath.Join("./template", "friends_list.html")

	tmpl, _ := template.ParseFiles(lp, fp)

	data := struct {
		IsAuth bool
		Users  []schema.User
		User   *schema.User
	}{true, users, user}

	tmpl.ExecuteTemplate(w, "layout", data)
}

func (h *Handler) FriendRequestPostHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {

	}

	ctx := r.Context()

	session, _ := h.SessionsStore.Get(r, "session")

	user, err := auth.CheckSession(session)
	if err != nil || user.ID == "" {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}

	friendId := r.FormValue("friend_id")

	friendship, err := db.GetFriendship(ctx, user.ID, friendId)
	if friendship.User != "" {
		http.Redirect(w, r, "/friends", http.StatusMovedPermanently)
	}

	err = db.InsertFriendship(ctx, user.ID, friendId)
	if err != nil {
		log.Println("create friendship", err)
	}
	http.Redirect(w, r, "/friends", http.StatusMovedPermanently)
}

func (h *Handler) FriendRequestAcceptHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {

	}

	ctx := r.Context()

	session, _ := h.SessionsStore.Get(r, "session")

	user, err := auth.CheckSession(session)
	if err != nil || user.ID == "" {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}

	friendId := r.FormValue("friend_id")

	friendship, err := db.GetFriendship(ctx, user.ID, friendId)
	if friendship.User != "" {
		http.Redirect(w, r, "/friends", http.StatusMovedPermanently)
	}

	err = db.InsertFriendship(ctx, friendId, user.ID)
	if err != nil {
		log.Println("create friendship", err)
	}
	http.Redirect(w, r, "/friends", http.StatusMovedPermanently)
}

func (h *Handler) FriendRequestDeclineHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {

	}

	ctx := r.Context()

	session, _ := h.SessionsStore.Get(r, "session")

	user, err := auth.CheckSession(session)
	if err != nil || user.ID == "" {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}

	friendId := r.FormValue("friend_id")

	friendship, err := db.GetFriendship(ctx, user.ID, friendId)
	if friendship.User == "" {
		http.Redirect(w, r, "/friends", http.StatusMovedPermanently)
	}

	err = db.DeleteFriendship(ctx, user.ID, friendId)
	if err != nil {
		log.Println("delete friendship", err)
	}
	http.Redirect(w, r, "/friends", http.StatusMovedPermanently)
}

func (h *Handler) FriendRequestListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	session, _ := h.SessionsStore.Get(r, "session")

	user, err := auth.CheckSession(session)
	if err != nil || user.ID == "" {
		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	}

	offsetQuery, ok := r.URL.Query()["offset"]

	offset := 0
	if ok {
		offset, _ = strconv.Atoi(offsetQuery[0])
	}

	limitQuery, ok := r.URL.Query()["limit"]

	limit := 50
	if ok {
		limit, _ = strconv.Atoi(limitQuery[0])
	}

	friendshipRequests, err := db.ListFriendship(ctx, user.ID)

	var requestsUserIDS []string

	for _, request := range friendshipRequests {
		requestsUserIDS = append(requestsUserIDS, request.Friend)
	}

	users, err := db.GetUsersByIDS(ctx, requestsUserIDS, offset, limit)

	lp := filepath.Join("./template", "layout.html")
	fp := filepath.Join("./template", "friend_requests.html")

	tmpl, _ := template.ParseFiles(lp, fp)

	data := struct {
		IsAuth bool
		Users  []schema.User
		User   *schema.User
	}{true, users, user}
	tmpl.ExecuteTemplate(w, "layout", data)
}

func (h *Handler) MyProfileHandler(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("./template", "layout.html")
	fp := filepath.Join("./template", "my_profile.html")

	tmpl, _ := template.ParseFiles(lp, fp)
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func (h *Handler) MyProfileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("MyProfileUpdateHandler")
}

func (h *Handler) UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	lp := filepath.Join("./template", "layout.html")
	mp := filepath.Join("./template", "menu.html")
	fp := filepath.Join("./template", "profile.html")

	tmpl, _ := template.ParseFiles(lp, mp, fp)
	tmpl.ExecuteTemplate(w, "layout", nil)
}
