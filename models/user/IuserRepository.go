package models

// IUserRepositoryInterface interface of the User repo
type IUserRepositoryInterface interface {
	Register(user User) error
	LogIn(userName, password string) (string, error)
	LogOut(token string) error
	ExistsUsername(userName string) (bool, error)
	ExistsEmail(email string) (bool, error)
	CheckToken(token string) (bool, error)
}
