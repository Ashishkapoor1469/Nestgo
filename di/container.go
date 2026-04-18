// Package di provides a compile-time safe, constructor-based dependency
// injection container for the NestGo framework.
package di

import (
	"fmt"
	"reflect"
	"sync"
)

// Scope defines the lifecycle scope of a provider.
type Scope int

const (
	// Singleton providers are created once and shared across the application.
	Singleton Scope = iota
	// RequestScoped providers are created per HTTP request.
	RequestScoped
	// Transient providers are created each time they are resolved.
	Transient
)

// Provider represents a service provider in the DI container.
type Provider struct {
	// Token is the unique name identifying this provider.
	Token string
	// Factory is the constructor function that creates the provider instance.
	// It should accept its dependencies as parameters and return (instance, error).
	Factory any
	// Instance holds the resolved singleton instance.
	Instance any
	// Scope determines the lifecycle of this provider.
	Scope Scope
	// Type is the reflect.Type of the provider's return value.
	Type reflect.Type
	// Deps lists the reflect.Types of the constructor's parameters.
	Deps []reflect.Type
}

// Container is the NestGo dependency injection container.
type Container struct {
	mu        sync.RWMutex
	providers map[reflect.Type]*Provider
	byToken   map[string]*Provider
	resolved  map[reflect.Type]any
	resolving map[reflect.Type]bool // for circular dependency detection
}

// NewContainer creates a new DI container.
func NewContainer() *Container {
	return &Container{
		providers: make(map[reflect.Type]*Provider),
		byToken:   make(map[string]*Provider),
		resolved:  make(map[reflect.Type]any),
		resolving: make(map[reflect.Type]bool),
	}
}

// Provide registers a constructor function as a provider.
// The constructor should be a function that takes its dependencies as parameters
// and returns (interface/struct, error) or just interface/struct.
//
// Example:
//
//	container.Provide(func(repo UserRepository) *UserService {
//	    return &UserService{repo: repo}
//	})
func (c *Container) Provide(factory any, opts ...ProviderOption) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	ft := reflect.TypeOf(factory)
	if ft.Kind() != reflect.Func {
		return fmt.Errorf("di: provider must be a function, got %T", factory)
	}

	if ft.NumOut() == 0 || ft.NumOut() > 2 {
		return fmt.Errorf("di: provider must return (value) or (value, error), got %d return values", ft.NumOut())
	}

	// Check second return value is error if present.
	if ft.NumOut() == 2 {
		errType := reflect.TypeOf((*error)(nil)).Elem()
		if !ft.Out(1).Implements(errType) {
			return fmt.Errorf("di: provider's second return value must be error, got %s", ft.Out(1))
		}
	}

	provider := &Provider{
		Factory: factory,
		Type:    ft.Out(0),
		Scope:   Singleton,
		Deps:    make([]reflect.Type, ft.NumIn()),
	}

	for i := 0; i < ft.NumIn(); i++ {
		provider.Deps[i] = ft.In(i)
	}

	// Apply options.
	for _, opt := range opts {
		opt(provider)
	}

	if provider.Token == "" {
		provider.Token = ft.Out(0).String()
	}

	c.providers[provider.Type] = provider
	c.byToken[provider.Token] = provider

	return nil
}

// ProvideValue registers an already-constructed value in the container.
func (c *Container) ProvideValue(value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	t := reflect.TypeOf(value)
	provider := &Provider{
		Token:    t.String(),
		Instance: value,
		Type:     t,
		Scope:    Singleton,
	}

	c.providers[t] = provider
	c.byToken[provider.Token] = provider
	c.resolved[t] = value

	return nil
}

// Resolve resolves a dependency by type.
func (c *Container) Resolve(target any) error {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr {
		return fmt.Errorf("di: resolve target must be a pointer, got %T", target)
	}

	elemType := targetType.Elem()
	instance, err := c.resolveType(elemType)
	if err != nil {
		return err
	}

	reflect.ValueOf(target).Elem().Set(reflect.ValueOf(instance))
	return nil
}

// MustResolve resolves a dependency or panics.
func (c *Container) MustResolve(target any) {
	if err := c.Resolve(target); err != nil {
		panic(err)
	}
}

// ResolveByType resolves a dependency by its reflect.Type.
func (c *Container) ResolveByType(t reflect.Type) (any, error) {
	return c.resolveType(t)
}

// ResolveAll resolves all registered providers and returns them.
func (c *Container) ResolveAll() ([]any, error) {
	c.mu.RLock()
	providers := make([]*Provider, 0, len(c.providers))
	for _, p := range c.providers {
		providers = append(providers, p)
	}
	c.mu.RUnlock()

	var instances []any
	for _, p := range providers {
		instance, err := c.resolveType(p.Type)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

// Has returns true if the container has a provider for the given type.
func (c *Container) Has(t reflect.Type) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.providers[t]
	return ok
}

// resolveType resolves a dependency by its type, handling circular dependencies.
func (c *Container) resolveType(t reflect.Type) (any, error) {
	c.mu.RLock()
	// Check if already resolved (singleton).
	if instance, ok := c.resolved[t]; ok {
		c.mu.RUnlock()
		return instance, nil
	}

	// Find provider, checking both direct type and interface implementations.
	provider := c.providers[t]
	if provider == nil {
		// Check if any provider's type implements the requested interface.
		if t.Kind() == reflect.Interface {
			for _, p := range c.providers {
				if p.Type.Implements(t) {
					provider = p
					break
				}
				if p.Type.Kind() == reflect.Ptr && p.Type.Implements(t) {
					provider = p
					break
				}
			}
		}
	}
	c.mu.RUnlock()

	if provider == nil {
		return nil, &ErrProviderNotFound{Type: t}
	}

	// Check for pre-constructed instance.
	if provider.Instance != nil {
		return provider.Instance, nil
	}

	// Check circular dependency.
	c.mu.Lock()
	if c.resolving[t] {
		c.mu.Unlock()
		return nil, &ErrCircularDependency{Type: t}
	}
	c.resolving[t] = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		delete(c.resolving, t)
		c.mu.Unlock()
	}()

	// Resolve dependencies.
	ft := reflect.TypeOf(provider.Factory)
	args := make([]reflect.Value, ft.NumIn())
	for i := 0; i < ft.NumIn(); i++ {
		dep, err := c.resolveType(ft.In(i))
		if err != nil {
			return nil, fmt.Errorf("di: resolving dependency %s for %s: %w", ft.In(i), t, err)
		}
		args[i] = reflect.ValueOf(dep)
	}

	// Call constructor.
	results := reflect.ValueOf(provider.Factory).Call(args)

	instance := results[0].Interface()

	// Check for error return.
	if len(results) == 2 && !results[1].IsNil() {
		return nil, fmt.Errorf("di: factory for %s returned error: %w", t, results[1].Interface().(error))
	}

	// Cache singleton.
	if provider.Scope == Singleton {
		c.mu.Lock()
		c.resolved[t] = instance
		provider.Instance = instance
		c.mu.Unlock()
	}

	return instance, nil
}

// Validate checks that all dependencies can be resolved.
func (c *Container) Validate() error {
	c.mu.RLock()
	providers := make([]*Provider, 0, len(c.providers))
	for _, p := range c.providers {
		providers = append(providers, p)
	}
	c.mu.RUnlock()

	for _, p := range providers {
		if p.Instance != nil {
			continue
		}
		for _, dep := range p.Deps {
			if !c.Has(dep) {
				// Check interface implementations.
				found := false
				if dep.Kind() == reflect.Interface {
					c.mu.RLock()
					for _, candidate := range c.providers {
						if candidate.Type.Implements(dep) || (candidate.Type.Kind() == reflect.Ptr && candidate.Type.Implements(dep)) {
							found = true
							break
						}
					}
					c.mu.RUnlock()
				}
				if !found {
					return fmt.Errorf("di: unresolved dependency %s required by %s", dep, p.Type)
				}
			}
		}
	}

	return nil
}

// Clone creates a child container with the same providers.
func (c *Container) Clone() *Container {
	c.mu.RLock()
	defer c.mu.RUnlock()

	child := NewContainer()
	for t, p := range c.providers {
		child.providers[t] = p
	}
	for t, p := range c.byToken {
		child.byToken[t] = p
	}
	// Singletons share resolved instances.
	for t, inst := range c.resolved {
		child.resolved[t] = inst
	}
	return child
}

// ProviderOption configures a provider.
type ProviderOption func(*Provider)

// WithToken sets a custom token for the provider.
func WithToken(token string) ProviderOption {
	return func(p *Provider) {
		p.Token = token
	}
}

// WithScope sets the scope for the provider.
func WithScope(scope Scope) ProviderOption {
	return func(p *Provider) {
		p.Scope = scope
	}
}

// AsRequestScoped marks the provider as request-scoped.
func AsRequestScoped() ProviderOption {
	return func(p *Provider) {
		p.Scope = RequestScoped
	}
}

// AsTransient marks the provider as transient.
func AsTransient() ProviderOption {
	return func(p *Provider) {
		p.Scope = Transient
	}
}
