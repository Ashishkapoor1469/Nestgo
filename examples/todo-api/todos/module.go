package todos

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
)

// TodosModule is the todos feature module.
type TodosModule struct{}

func (m *TodosModule) Module() common.ModuleConfig {
	service := NewTodoService()
	controller := NewTodoController(service)

	return common.ModuleConfig{
		Name:        "todos",
		Controllers: []common.Controller{controller},
		Providers: []di.Provider{
			{Instance: service},
		},
	}
}
