package mock

import (
	"time"

	"orinayooyelade.com/snippetbox/pkg/models"
)

var mockUser = &models.User{
	ID:      1,
	Name:    "Andy",
	Email:   "andy@example.com",
	Created: time.Now(),
	Active:  true,
}

type UserModel struct{}

func (model *UserModel) Insert(name, email, password string) error {
	switch email {
	case "andy@example.com":
		return models.ErrDuplicateEmail
	default:
		return nil
	}
}

func (model *UserModel) Authenticate(email, password string) (int, error) {
	switch email {
	case "andy@example.com":
		return 1, nil
	default:
		return 0, models.ErrInvalidCredentials
	}
}

func (model *UserModel) Get(id int) (*models.User, error) {
	switch id {
	case 1:
		return mockUser, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (model *UserModel) ChangePassword(int, string, string) error {
	return nil
}
