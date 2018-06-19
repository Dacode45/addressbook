package models

// User is a wrapper around a list of contacts.
type User struct {
	UserID   string `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Contacts []Contact
}
