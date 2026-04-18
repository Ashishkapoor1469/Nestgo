package modules

import (
	"github.com/Ashishkapoor1469/Nestgo/di"
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/sample-app/internal/modules/auth"
)

// AppModule is the root module of the application.
type AppModule struct{}

func (m *AppModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name:        "app",
		Imports: []common.Module{
			&auth.AuthModule{},
		},
		Controllers: []common.Controller{
			NewAppController(),
		},
		Providers:   []di.Provider{},
	}
}
