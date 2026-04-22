package guards

import (
	"strings"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// JWTGuard is an authentication guard that validates JWT tokens.
// It extracts the bearer token from the Authorization header,
// validates it, and sets user_id and user_email in the context.
//
// Usage:
//
//	guard := guards.NewJWTGuard(func(token string) (string, error) {
//	    claims, err := authService.ValidateToken(token)
//	    if err != nil {
//	        return "", err
//	    }
//	    return claims.UserID, nil
//	})
type JWTGuard struct {
	validator func(token string) (userID string, err error)
}

// NewJWTGuard creates a new JWT guard.
func NewJWTGuard(validator func(token string) (string, error)) *JWTGuard {
	return &JWTGuard{validator: validator}
}

// CanActivate checks if the request has a valid JWT token.
func (g *JWTGuard) CanActivate(ctx *common.Context) (bool, error) {
	auth := ctx.Header("Authorization")
	if auth == "" {
		return false, common.Unauthorized("missing authorization header")
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		return false, common.Unauthorized("invalid authorization format")
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	userID, err := g.validator(token)
	if err != nil {
		return false, common.Unauthorized("invalid or expired token")
	}

	ctx.Set("user_id", userID)
	return true, nil
}
