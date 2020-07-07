package db

import (
	"context"
	"social-network/pkg/schema"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := MysqlRepository{
		db: db,
	}

	now := time.Now().UTC()
	user := schema.User{
		ID:        "1234567",
		Name:      "test",
		Email:     "test@mail.ru",
		Surname:   "test",
		Sex:       "male",
		Interests: "",
		Password:  "1q2w3e4rt5",
		CreatedAt: now,
	}

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.ID, user.Name, user.Email, user.Surname, user.Sex, user.Interests, user.Password, user.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.InsertUser(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
