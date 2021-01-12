package models

import (
	"errors"
	"time"
)

var (
	// ErrNoRecord is the response for records not found in any DB
	ErrNoRecord = errors.New("models: no matching record")
	// ErrInvalidCredentials is the response for incorrect email or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail is the response for trying to register an email already in use
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

// Snippet struct is the DAO for Snippet Model
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// User struct is the DAO for User Model
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
