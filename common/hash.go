package common

// Hash generates and compares encrypted strings
type Hash interface {
	Generate(s string) (string, error)
	Compare(hash string, s string) error
}
