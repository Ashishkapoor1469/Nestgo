package health

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
)

type HealthModule struct{}

func (m *HealthModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name:        "health",
		Controllers: []common.Controller{NewHealthController()},
		Providers:   []di.Provider{},
	}
}
