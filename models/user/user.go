package models

// User represents the structure of our resource
type User struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	ID       string `json:"id"`
}
