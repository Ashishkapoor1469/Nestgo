package common

import "github.com/Ashishkapoor1469/Nestgo/di"

// Module is implemented by all NestGo modules.
type Module interface {
	Module() ModuleConfig
}

// ModuleConfig describes a module's dependencies, controllers, and providers.
type ModuleConfig struct {
	// Name is the module identifier.
	Name string
	// Imports lists other modules this module depends on.
	Imports []Module
	// Controllers lists the HTTP controllers in this module.
	Controllers []Controller
	// Providers lists the DI providers (services) in this module.
	Providers []di.Provider
	// Exports lists provider types that are available to importing modules.
	Exports []any
}
