package auth

import (
	"strings"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// AuthGuard is a simple mock authentication guard.
// Guards in NestGo are used to protect routes from unauthorized access.
// They implement the common.Guard interface (CanActivate method).
type AuthGuard struct{}

// NewAuthGuard creates a new instance of the mock AuthGuard.
func NewAuthGuard() *AuthGuard {
	return &AuthGuard{}
}

// CanActivate evaluates if the current request is allowed to proceed.
func (g *AuthGuard) CanActivate(ctx *common.Context) (bool, error) {
	// 1. Extract the Authorization header using standard Context methods.
	authHeader := ctx.Header("Authorization")
	
	// 2. Check if the header exists and starts with "Bearer "
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		// Return common.Unauthorized to send a built-in 401 JSON exception.
		return false, common.Unauthorized("Missing or invalid authorization token")
	}

	// 3. Extract the token itself.
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// 4. Validate the token.
	// In a real application, you would verify a JWT, check database, etc.
	// Here, we just accept "mock-super-secret-token" for demonstration.
	if token != "mock-super-secret-token" {
		return false, common.Unauthorized("Invalid token")
	}

	// 5. Store user information in the context so subsequent handlers can access it.
	ctx.Set("user_id", "user-123")
	
	// Return true to allow the request to proceed.
	return true, nil
}
