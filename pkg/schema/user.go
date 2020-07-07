package schema

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Surname      string    `json:"surname"`
	Name         string    `json:"name"`
	City         string    `json:"city"`
	Sex          string    `json:"sex"`
	Password     string    `json:"password"`
	Interests    string    `json:"interests"`
	CreatedAt    time.Time `json:"created_at"`
	IsMineFriend bool      `json:"is_mine_friend"`
}
