package auth

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
)

// AuthModule provides authentication and authorization.
type AuthModule struct{}

func (m *AuthModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name: "auth",
		Controllers: []common.Controller{
			&AuthController{},
		},
		Providers: []di.Provider{
			{
				Factory: NewAuthService,
			},
			{
				Factory: NewAuthController,
			},
		},
	}
}
