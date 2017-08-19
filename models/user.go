package models

// User represents the structure of our resource
type User struct {
	UserName string `json:"UserName"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}
