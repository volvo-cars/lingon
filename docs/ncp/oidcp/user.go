package oidcp

import (
	"github.com/zitadel/oidc/v2/example/server/storage"
	"golang.org/x/text/language"
)

type UserStore interface {
	GetUserByID(string) *storage.User
	GetUserByUsername(string) *storage.User
	ExampleClientID() string
}

type userStore struct {
	users map[string]*storage.User
}

type UserStoreOption func(*userStore)

func WithUser(id string, user *storage.User) UserStoreOption {
	return func(s *userStore) {
		s.users[id] = user
	}
}

func WithAdminUser() UserStoreOption {
	return func(s *userStore) {
		s.users["admin"] = &storage.User{
			ID:                "admin",
			Username:          "admin",
			Password:          "admin",
			FirstName:         "Admin",
			LastName:          "Admin",
			Email:             "admin@localhost",
			EmailVerified:     true,
			Phone:             "",
			PhoneVerified:     false,
			PreferredLanguage: language.BritishEnglish,
			IsAdmin:           true,
		}
	}
}

func NewUserStore(opts ...UserStoreOption) UserStore {
	us := userStore{
		users: map[string]*storage.User{},
	}
	for _, opt := range opts {
		opt(&us)
	}
	return us
}

// ExampleClientID is only used in the example server
func (u userStore) ExampleClientID() string {
	return "service"
}

func (u userStore) GetUserByID(id string) *storage.User {
	return u.users[id]
}

func (u userStore) GetUserByUsername(username string) *storage.User {
	for _, user := range u.users {
		if user.Username == username {
			return user
		}
	}
	return nil
}
