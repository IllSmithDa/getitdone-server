package models

type Todoitem struct {
	ID int `json:"id", binding:"required"`
	Todoitem string `json:"todoitem", binding: "required"`
	Username string   `json:"username" binding:"required"`
}

// User model with Username and password
type User struct {
	Username string   `json:"username" binding:"required"`
	Password string   `json:"password" binding:"required"`
}