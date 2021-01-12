package mysql

import (
	"database/sql"

	"orinayooyelade.com/snippetbox/pkg/models"
)

// UserModel struct which wraps a sql.DB connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert will insert a new user into the database
func (model *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate will verify a user exists with the provided email and password
func (model *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get will return a specific user based on its id
func (model *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
