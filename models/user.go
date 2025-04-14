package models

import (
	"errors"
	"sync"
)

var ErrUserExists = errors.New("user already exists")

type User struct {
	Username string
	Password string 
}

var (
	users   = map[string]string{} // username: hashed password
	userMux = sync.Mutex{}
)

func AddUser(username, hashedPassword string) error {
	userMux.Lock()
	defer userMux.Unlock()

	if _, exists := users[username]; exists {
		return ErrUserExists
	}

	users[username] = hashedPassword
	return nil
}

func Authenticate(username, hashedPassword string) bool {
	userMux.Lock()
	defer userMux.Unlock()

	stored, ok := users[username]
	return ok && stored == hashedPassword
}

