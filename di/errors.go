package di

import (
	"fmt"
	"reflect"
)

// ErrProviderNotFound is returned when a provider is not registered for the given type.
type ErrProviderNotFound struct {
	Type reflect.Type
}

func (e *ErrProviderNotFound) Error() string {
	return fmt.Sprintf("di: no provider found for type %s", e.Type)
}

// ErrCircularDependency is returned when a circular dependency is detected.
type ErrCircularDependency struct {
	Type reflect.Type
}

func (e *ErrCircularDependency) Error() string {
	return fmt.Sprintf("di: circular dependency detected for type %s", e.Type)
}

// ErrDuplicateProvider is returned when a provider is already registered.
type ErrDuplicateProvider struct {
	Type reflect.Type
}

func (e *ErrDuplicateProvider) Error() string {
	return fmt.Sprintf("di: duplicate provider for type %s", e.Type)
}

// ErrResolutionFailed is returned when dependency resolution fails.
type ErrResolutionFailed struct {
	Type  reflect.Type
	Cause error
}

func (e *ErrResolutionFailed) Error() string {
	return fmt.Sprintf("di: resolution failed for type %s: %v", e.Type, e.Cause)
}

func (e *ErrResolutionFailed) Unwrap() error {
	return e.Cause
}
