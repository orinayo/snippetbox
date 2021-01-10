package models

import (
	"errors"
	"time"
)

// ErrNoRecord is the response for records not found in any DB ...
var ErrNoRecord = errors.New("models: no matching record")

// Snippet is the DAO for Snippet Model ...
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
