package config

import "os"

var WikiServiceURL string = "http://127.0.0.1:9454"

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-super-secret-jwt-key-change-this-in-production"
	}
	return secret
}
