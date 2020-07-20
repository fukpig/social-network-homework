package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/bgadrian/fastfaker/faker"
	"github.com/kelseyhightower/envconfig"

	_ "github.com/go-sql-driver/mysql"
)

const (
	gophers = 20
	entries = 50000
)

type Config struct {
	MysqlHost     string `envconfig:"MYSQL_HOST"`
	MysqlDB       string `envconfig:"MYSQL_DB"`
	MysqlUser     string `envconfig:"MYSQL_USER"`
	MysqlPassword string `envconfig:"MYSQL_PASSWORD"`
}

/*func main() {
	count := flag.Int("usersCount", 0, "an int")
	flag.Parse()

	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	cfg.MysqlHost = "localhost:3306"
	cfg.MysqlUser = "user"
	cfg.MysqlPassword = "123456"
	cfg.MysqlDB = "app"

	addr := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.MysqlUser, cfg.MysqlPassword, cfg.MysqlHost, cfg.MysqlDB)
	db, err := sql.Open("mysql", addr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	log.Println("Faking started")

	createdAt := time.Now().UTC()

	stmt, err := db.Prepare("INSERT INTO users(email, name, surname, sex, city, interests, password, created_at) VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		fmt.Println("ERROR", err)
	}

	defer stmt.Close()
	defer db.Close()

	for i := 1; i <= *count; i++ {
		wg.Add(1)
		go func() {
			routineFaker := faker.NewFastFaker()

			fullname := routineFaker.Name()
			nameParts := strings.Split(fullname, " ")

			_, err = stmt.Exec(routineFaker.Email(), nameParts[0], nameParts[1], routineFaker.Gender(), routineFaker.City(), routineFaker.HackerPhrase(), routineFaker.PasswordFull(), createdAt)

			if err != nil {
				log.Println("INSERT ERROR", err)
			}

			wg.Done()
		}()
	}

	wg.Wait()
	log.Println("Faking finished. Done:", *count)
}*/

func main() {

	var sStmt string = "INSERT INTO users(email, name, surname, sex, city, interests, password, created_at) VALUES(?,?,?,?,?,?,?,?)"

	var wg sync.WaitGroup
	for i := 0; i <= gophers; i++ {
		wg.Add(1)
		go inserter(wg, sStmt)
	}

	wg.Wait()

}

func inserter(wg sync.WaitGroup, sStmt string) {
	createdAt := time.Now().UTC()
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	cfg.MysqlHost = "localhost:3306"
	cfg.MysqlUser = "user"
	cfg.MysqlPassword = "123456"
	cfg.MysqlDB = "app"

	addr := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.MysqlUser, cfg.MysqlPassword, cfg.MysqlHost, cfg.MysqlDB)
	db, err := sql.Open("mysql", addr)

	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.Prepare(sStmt)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < entries; i++ {
		routineFaker := faker.NewFastFaker()
		fullname := routineFaker.Name()
		nameParts := strings.Split(fullname, " ")
		res, err := stmt.Exec(routineFaker.Email(), nameParts[0], nameParts[1], routineFaker.Gender(), routineFaker.City(), routineFaker.HackerPhrase(), routineFaker.PasswordFull(), createdAt)
		if err != nil || res == nil {
			log.Fatal(err)
		}
	}
	stmt.Close()
	db.Close()
	wg.Done()
}
