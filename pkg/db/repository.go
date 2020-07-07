package db

import (
	"context"
	"social-network/pkg/schema"
)

type Repository interface {
	Close()

	GetUserByEmail(ctx context.Context, email string) (*schema.User, error)
	InsertUser(ctx context.Context, user schema.User) error
	ListUsers(ctx context.Context, userID string, offset int, limit int) ([]schema.User, error)
	GetUsersByIDS(ctx context.Context, userIDS []string, offset int, limit int) ([]schema.User, error)
	ListFriends(ctx context.Context, userID string) ([]string, error)
	GetFriendship(ctx context.Context, userID string, friendID string) (*schema.Friendship, error)
	InsertFriendship(ctx context.Context, userID string, friendID string) error
	DeleteFriendship(ctx context.Context, userID string, friendID string) error
	ListFriendship(ctx context.Context, userID string) ([]schema.Friendship, error)
}

var impl Repository

func SetRepository(repository Repository) {
	impl = repository
}

func Close() {
	impl.Close()
}

func GetUserByEmail(ctx context.Context, email string) (*schema.User, error) {
	return impl.GetUserByEmail(ctx, email)
}

func InsertUser(ctx context.Context, user schema.User) error {
	return impl.InsertUser(ctx, user)
}

func ListUsers(ctx context.Context, userID string, offset int, limit int) ([]schema.User, error) {
	return impl.ListUsers(ctx, userID, offset, limit)
}

func GetUsersByIDS(ctx context.Context, userIDS []string, offset int, limit int) ([]schema.User, error) {
	return impl.GetUsersByIDS(ctx, userIDS, offset, limit)
}

func ListFriends(ctx context.Context, userID string) ([]string, error) {
	return impl.ListFriends(ctx, userID)
}

func GetFriendship(ctx context.Context, userID string, friendID string) (*schema.Friendship, error) {
	return impl.GetFriendship(ctx, userID, friendID)
}

func InsertFriendship(ctx context.Context, userID string, friendID string) error {
	return impl.InsertFriendship(ctx, userID, friendID)
}

func DeleteFriendship(ctx context.Context, userID string, friendID string) error {
	return impl.DeleteFriendship(ctx, userID, friendID)
}

func ListFriendship(ctx context.Context, userID string) ([]schema.Friendship, error) {
	return impl.ListFriendship(ctx, userID)
}
