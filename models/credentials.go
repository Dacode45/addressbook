package models

// Credentials are used for logging in
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
