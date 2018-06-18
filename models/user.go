package models

type User struct {
	UserId   string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
