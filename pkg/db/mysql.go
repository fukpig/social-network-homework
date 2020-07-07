package db

import (
	"context"
	"database/sql"
	"log"
	"social-network/pkg/schema"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlRepository struct {
	db *sql.DB
}

func NewMysql(url string) (*MysqlRepository, error) {
	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &MysqlRepository{
		db,
	}, nil
}

func (r *MysqlRepository) Close() {
	r.db.Close()
}

func (r *MysqlRepository) GetUserByEmail(ctx context.Context, email string) (*schema.User, error) {
	user := new(schema.User)
	err := r.db.QueryRowContext(ctx, "SELECT id, email, name, surname, sex, city, interests, password FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Surname, &user.Sex, &user.City, &user.Interests, &user.Password)
	return user, err
}

func (r *MysqlRepository) InsertUser(ctx context.Context, user schema.User) error {
	stmt, err := r.db.Prepare("INSERT INTO users(id, email, name, surname, sex, city, interests, password, created_at) VALUES(?,?,?,?,?,?,?,?,?)")
	_, err = stmt.Exec(user.ID, user.Email, user.Name, user.Surname, user.Sex, user.City, user.Interests, user.Password, user.CreatedAt)
	return err
}

func (r *MysqlRepository) ListUsers(ctx context.Context, userID string, offset int, limit int) ([]schema.User, error) {
	users := []schema.User{}
	rows, err := r.db.QueryContext(ctx, "SELECT id, email, name, surname, sex, city, interests, password FROM users WHERE id <> ? LIMIT ? OFFSET ?", userID, limit, offset)
	defer rows.Close()

	if err != nil {
		return users, err
	}

	for rows.Next() {
		user := schema.User{}
		if err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Surname, &user.Sex, &user.City, &user.Interests, &user.Password); err == nil {
			users = append(users, user)
		}
	}
	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (r *MysqlRepository) GetUsersByIDS(ctx context.Context, userIDS []string, offset int, limit int) ([]schema.User, error) {
	users := []schema.User{}

	args := make([]interface{}, len(userIDS))
	for i, id := range userIDS {
		args[i] = id
	}

	stmt := ""

	if len(args) > 0 {
		stmt = `SELECT id, email, name, surname, sex, city, interests, password FROM users WHERE id IN (?` + strings.Repeat(",?", len(args)-1) + `)`
		rows, err := r.db.Query(stmt, args...)

		defer rows.Close()

		if err != nil {
			log.Println("Get users by ids error:", err)
			return users, err
		}

		for rows.Next() {
			user := schema.User{}
			if err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Surname, &user.Sex, &user.City, &user.Interests, &user.Password); err == nil {
				users = append(users, user)
			}
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return users, nil
}

func (r *MysqlRepository) ListFriends(ctx context.Context, userID string) ([]string, error) {
	var users []string
	rows, err := r.db.QueryContext(ctx, "SELECT f1.friend from friendship f1 inner join friendship f2 on f1.user = f2.friend and f1.friend = f2.user WHERE f1.user = ?", userID)
	defer rows.Close()

	if err != nil {
		log.Println("Get friends error:", err)
		return users, err
	}

	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err == nil {
			users = append(users, id)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *MysqlRepository) GetFriendship(ctx context.Context, user string, friend string) (*schema.Friendship, error) {
	friendship := new(schema.Friendship)
	err := r.db.QueryRowContext(ctx, "SELECT user, friend FROM friendship WHERE user = ? AND friend = ?", user, friend).Scan(&friendship.User, &friendship.Friend)
	return friendship, err
}

func (r *MysqlRepository) InsertFriendship(ctx context.Context, user string, friend string) error {
	stmt, err := r.db.Prepare("INSERT INTO friendship(user, friend) VALUES(?,?)")
	_, err = stmt.Exec(user, friend)
	return err
}

func (r *MysqlRepository) DeleteFriendship(ctx context.Context, user string, friend string) error {
	stmt, err := r.db.Prepare("DELETE FROM friendship WHERE user = ? and friend = ?")
	_, err = stmt.Exec(user, friend)
	return err
}

func (r *MysqlRepository) ListFriendship(ctx context.Context, userID string) ([]schema.Friendship, error) {
	var friendships []schema.Friendship
	rows, err := r.db.QueryContext(ctx, "SELECT f1.user, f1.friend from friendship f1 inner join friendship f2 on f1.friend = f2.user WHERE f1.user = ?", userID)
	defer rows.Close()

	if err != nil {
		log.Println("Get friends error:", err)
		return friendships, err
	}

	for rows.Next() {
		var friendship schema.Friendship
		if err = rows.Scan(&friendship.User, &friendship.Friend); err == nil {
			friendships = append(friendships, friendship)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friendships, nil
}
