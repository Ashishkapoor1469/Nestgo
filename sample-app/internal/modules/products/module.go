package products

import (
	"github.com/Ashishkapoor1469/Nestgo/di"
	"github.com/Ashishkapoor1469/Nestgo/common"
)

// ProductsModule defines the products feature module.
type ProductsModule struct{}

func (m *ProductsModule) Module() common.ModuleConfig {
	service := NewProductsService()
	controller := NewProductsController(service)

	return common.ModuleConfig{
		Name: "products",
		Controllers: []common.Controller{controller},
		Providers: []di.Provider{
			{Instance: service},
		},
	}
}
