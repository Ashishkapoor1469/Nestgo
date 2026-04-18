package auth

import "errors"

// AuthService handles authentication business logic.
// In NestGo, services are typically injected as providers into controllers.
type AuthService struct{}

// NewAuthService creates a new authentication service.
func NewAuthService() *AuthService {
	return &AuthService{}
}

// LoginDTO represents the JSON payload expected for the login route.
// Implementing common.Validatable (optional) allows the framework
// to automatically run validation during binding.
type LoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login verifies credentials and returns a token.
func (s *AuthService) Login(username, password string) (string, error) {
	// Simple mock authentication logic
	if username == "admin" && password == "password123" {
		// Let's pretend we generated a JWT here.
		return "mock-super-secret-token", nil
	}
	return "", errors.New("invalid credentials")
}
