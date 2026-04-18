package core

import (
	"fmt"

	"github.com/nestgo/nestgo/common"
)

// ModuleRegistry manages module registration and initialization.
type ModuleRegistry struct {
	modules    []common.Module
	registered map[string]bool
	order      []common.Module
}

// NewModuleRegistry creates a new module registry.
func NewModuleRegistry() *ModuleRegistry {
	return &ModuleRegistry{
		registered: make(map[string]bool),
	}
}

// Register adds a module to the registry.
func (r *ModuleRegistry) Register(m common.Module) error {
	cfg := m.Module()
	if cfg.Name == "" {
		return fmt.Errorf("modules: module name cannot be empty")
	}
	if r.registered[cfg.Name] {
		return nil // already registered, skip (idempotent)
	}
	r.registered[cfg.Name] = true
	r.modules = append(r.modules, m)
	return nil
}

// ResolveOrder returns modules in dependency order (topological sort).
func (r *ModuleRegistry) ResolveOrder() ([]common.Module, error) {
	if r.order != nil {
		return r.order, nil
	}

	visited := make(map[string]bool)
	visiting := make(map[string]bool) // for cycle detection
	var order []common.Module

	var visit func(m common.Module) error
	visit = func(m common.Module) error {
		cfg := m.Module()
		name := cfg.Name

		if visiting[name] {
			return fmt.Errorf("modules: circular dependency detected involving module %q", name)
		}
		if visited[name] {
			return nil
		}

		visiting[name] = true
		for _, imp := range cfg.Imports {
			// Register imported module if not already registered.
			if err := r.Register(imp); err != nil {
				return err
			}
			if err := visit(imp); err != nil {
				return err
			}
		}
		delete(visiting, name)

		visited[name] = true
		order = append(order, m)
		return nil
	}

	for _, m := range r.modules {
		if err := visit(m); err != nil {
			return nil, err
		}
	}

	r.order = order
	return order, nil
}

// Modules returns all registered modules.
func (r *ModuleRegistry) Modules() []common.Module {
	return r.modules
}

// GetModule returns a module by name.
func (r *ModuleRegistry) GetModule(name string) (common.Module, bool) {
	for _, m := range r.modules {
		if m.Module().Name == name {
			return m, true
		}
	}
	return nil, false
}

// DependencyGraph returns a map of module name to its dependency names.
func (r *ModuleRegistry) DependencyGraph() map[string][]string {
	graph := make(map[string][]string)
	for _, m := range r.modules {
		cfg := m.Module()
		var deps []string
		for _, imp := range cfg.Imports {
			deps = append(deps, imp.Module().Name)
		}
		graph[cfg.Name] = deps
	}
	return graph
}
