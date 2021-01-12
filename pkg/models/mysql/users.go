package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"orinayooyelade.com/snippetbox/pkg/models"
)

// UserModel struct which wraps a sql.DB connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert will insert a new user into the database
func (model *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = model.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
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
