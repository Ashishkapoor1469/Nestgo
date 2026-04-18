// Package guards provides authorization guards for the NestGo framework.
package guards

import (
	"strings"

	"github.com/Ashishkapoor1469/Nestgo/common"
)

// Guard is the authorization interface.
type Guard interface {
	CanActivate(ctx *common.Context) (bool, error)
}

// RBACGuard implements role-based access control.
type RBACGuard struct {
	requiredRoles []string
	roleExtractor func(ctx *common.Context) []string
}

// NewRBACGuard creates a new RBAC guard.
func NewRBACGuard(roles []string, extractor func(ctx *common.Context) []string) *RBACGuard {
	return &RBACGuard{
		requiredRoles: roles,
		roleExtractor: extractor,
	}
}

// CanActivate checks if the user has any of the required roles.
func (g *RBACGuard) CanActivate(ctx *common.Context) (bool, error) {
	userRoles := g.roleExtractor(ctx)
	for _, required := range g.requiredRoles {
		for _, userRole := range userRoles {
			if strings.EqualFold(required, userRole) {
				return true, nil
			}
		}
	}
	return false, common.Forbidden("insufficient role permissions")
}

// PermissionGuard implements permission-based access control.
type PermissionGuard struct {
	requiredPermissions []string
	permExtractor       func(ctx *common.Context) []string
	requireAll          bool
}

// NewPermissionGuard creates a new permission guard.
func NewPermissionGuard(permissions []string, extractor func(ctx *common.Context) []string) *PermissionGuard {
	return &PermissionGuard{
		requiredPermissions: permissions,
		permExtractor:       extractor,
		requireAll:          false,
	}
}

// RequireAll sets whether all permissions are required (AND) vs any (OR).
func (g *PermissionGuard) RequireAll() *PermissionGuard {
	g.requireAll = true
	return g
}

// CanActivate checks if the user has the required permissions.
func (g *PermissionGuard) CanActivate(ctx *common.Context) (bool, error) {
	userPerms := g.permExtractor(ctx)
	permSet := make(map[string]bool, len(userPerms))
	for _, p := range userPerms {
		permSet[strings.ToLower(p)] = true
	}

	if g.requireAll {
		// All permissions required.
		for _, required := range g.requiredPermissions {
			if !permSet[strings.ToLower(required)] {
				return false, common.Forbidden("missing permission: " + required)
			}
		}
		return true, nil
	}

	// Any permission is sufficient.
	for _, required := range g.requiredPermissions {
		if permSet[strings.ToLower(required)] {
			return true, nil
		}
	}
	return false, common.Forbidden("insufficient permissions")
}

// AuthGuard is a simple authentication guard that checks for a bearer token.
type AuthGuard struct {
	validator func(token string) (userID string, err error)
}

// NewAuthGuard creates a new auth guard with a token validator.
func NewAuthGuard(validator func(token string) (string, error)) *AuthGuard {
	return &AuthGuard{validator: validator}
}

// CanActivate checks if the request has a valid bearer token.
func (g *AuthGuard) CanActivate(ctx *common.Context) (bool, error) {
	token := ctx.BearerToken()
	if token == "" {
		return false, common.Unauthorized("missing authentication token")
	}

	userID, err := g.validator(token)
	if err != nil {
		return false, common.Unauthorized("invalid authentication token")
	}

	ctx.Set("user_id", userID)
	return true, nil
}

// CompositeGuard combines multiple guards with AND logic.
type CompositeGuard struct {
	guards []Guard
}

// NewCompositeGuard creates a guard that requires all guards to pass.
func NewCompositeGuard(guards ...Guard) *CompositeGuard {
	return &CompositeGuard{guards: guards}
}

// CanActivate returns true only if all guards pass.
func (g *CompositeGuard) CanActivate(ctx *common.Context) (bool, error) {
	for _, guard := range g.guards {
		allowed, err := guard.CanActivate(ctx)
		if err != nil {
			return false, err
		}
		if !allowed {
			return false, nil
		}
	}
	return true, nil
}
