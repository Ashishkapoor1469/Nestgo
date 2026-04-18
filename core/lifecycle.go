// Package core provides the application bootstrap, lifecycle management,
// and request context for the NestGo framework.
package core

// OnInit is implemented by providers that need initialization after DI resolution.
type OnInit interface {
	OnInit() error
}

// OnStart is implemented by providers that need to run logic when the server starts.
type OnStart interface {
	OnStart() error
}

// OnShutdown is implemented by providers that need cleanup during graceful shutdown.
type OnShutdown interface {
	OnShutdown() error
}

// OnDestroy is implemented by providers that need cleanup when the module is destroyed.
type OnDestroy interface {
	OnDestroy() error
}

// LifecycleHost manages lifecycle hooks for all registered providers.
type LifecycleHost struct {
	initHooks     []OnInit
	startHooks    []OnStart
	shutdownHooks []OnShutdown
	destroyHooks  []OnDestroy
}

// NewLifecycleHost creates a new lifecycle host.
func NewLifecycleHost() *LifecycleHost {
	return &LifecycleHost{}
}

// Register checks if the provider implements lifecycle hooks and registers them.
func (lh *LifecycleHost) Register(provider any) {
	if v, ok := provider.(OnInit); ok {
		lh.initHooks = append(lh.initHooks, v)
	}
	if v, ok := provider.(OnStart); ok {
		lh.startHooks = append(lh.startHooks, v)
	}
	if v, ok := provider.(OnShutdown); ok {
		lh.shutdownHooks = append(lh.shutdownHooks, v)
	}
	if v, ok := provider.(OnDestroy); ok {
		lh.destroyHooks = append(lh.destroyHooks, v)
	}
}

// RunInitHooks runs all OnInit hooks in registration order.
func (lh *LifecycleHost) RunInitHooks() error {
	for _, h := range lh.initHooks {
		if err := h.OnInit(); err != nil {
			return err
		}
	}
	return nil
}

// RunStartHooks runs all OnStart hooks in registration order.
func (lh *LifecycleHost) RunStartHooks() error {
	for _, h := range lh.startHooks {
		if err := h.OnStart(); err != nil {
			return err
		}
	}
	return nil
}

// RunShutdownHooks runs all OnShutdown hooks in reverse registration order.
func (lh *LifecycleHost) RunShutdownHooks() error {
	var firstErr error
	for i := len(lh.shutdownHooks) - 1; i >= 0; i-- {
		if err := lh.shutdownHooks[i].OnShutdown(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// RunDestroyHooks runs all OnDestroy hooks in reverse registration order.
func (lh *LifecycleHost) RunDestroyHooks() error {
	var firstErr error
	for i := len(lh.destroyHooks) - 1; i >= 0; i-- {
		if err := lh.destroyHooks[i].OnDestroy(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
