package users

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
)

// UsersModule defines the users feature module.
type UsersModule struct{}

func (m *UsersModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name: "users",
		Controllers: []common.Controller{
			&UsersController{},
		},
		Providers: []di.Provider{
			{
				Factory: NewUsersService,
			},
			{
				Factory: NewUsersController,
			},
		},
	}
}
