package auth

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
)

// AuthModule bundles the auth controller, service, and dependencies together.
type AuthModule struct{}

// Module implements common.Module and registers its components.
func (m *AuthModule) Module() common.ModuleConfig {
	authService := NewAuthService()
	authController := NewAuthController(authService)

	return common.ModuleConfig{
		Name: "auth",
		Controllers: []common.Controller{
			authController,
		},
		Providers: []di.Provider{
			{Instance: authService},
		},
	}
}
