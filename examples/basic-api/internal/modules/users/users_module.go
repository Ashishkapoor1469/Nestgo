package users

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
)

type UsersModule struct{}

func (m *UsersModule) Module() common.ModuleConfig {
	svc := NewUsersService()
	ctrl := NewUsersController(svc)

	return common.ModuleConfig{
		Name:        "users",
		Controllers: []common.Controller{ctrl},
		Providers:   []di.Provider{{Instance: svc}},
	}
}
