package server

// ServerConfig is special parameters for the server. In this case it only expects a secret for jwt encoding
type ServerConfig struct {
	JWTSecret string
}
