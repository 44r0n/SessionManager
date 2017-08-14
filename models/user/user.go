package models


// User represents the structure of our resource
type User struct {
        UserName   string `json:"name"`
        Email string `json:"email"`
        Password    string    `json:"pass"`
        ID     string `json:"id"`
    }
