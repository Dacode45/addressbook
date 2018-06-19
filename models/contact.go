package models

type Contact struct {
	ID        string `json:"id" csv:"id"`
	FirstName string `json:"first_name" csv:"first_name"`
	LastName  string `json:"last_name" csv:"last_name"`
	Email     string `json:"email" csv:"email"`
	Phone     string `json:"phone" csv:"phone"`
}
