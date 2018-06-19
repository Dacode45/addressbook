package mock

// Hash mocks the Hash interface by not encrypting anything and always passing
type Hash struct{}

// Generate does nothing
func (h *Hash) Generate(s string) (string, error) {
	return s, nil
}

// Compare allways passes
func (h *Hash) Compare(hash string, s string) error {
	return nil
}
