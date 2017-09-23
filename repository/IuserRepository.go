package repository

import (
	"github.com/44r0n/SessionManager/models"
)

// IUserRepositoryInterface interface of the User repo
type IUserRepositoryInterface interface {
	Register(user models.User) error
	GetIDAndPassword(userName string) (string, string, error)
	CreateToken(userID, token string) error
	DeleteToken(userID, token string) error
	ExistsUsername(userName string) (bool, error)
	ExistsEmail(email string) (bool, error)
	CheckToken(token string) (bool, error)
}
