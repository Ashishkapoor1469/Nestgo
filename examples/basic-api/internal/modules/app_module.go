package modules

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
	"github.com/Ashishkapoor1469/Nestgo/examples/basic-api/internal/modules/health"
	"github.com/Ashishkapoor1469/Nestgo/examples/basic-api/internal/modules/users"
)

// AppModule is the root application module.
type AppModule struct{}

func (m *AppModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name: "app",
		Imports: []common.Module{
			&health.HealthModule{},
			&users.UsersModule{},
		},
		Controllers: []common.Controller{},
		Providers:   []di.Provider{},
	}
}
