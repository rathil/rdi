package rdi

import "errors"

type (
	DI interface {
		// Override tries to register the given provider in the current container.
		// If registration fails (e.g., due to a conflict), it creates a new child container,
		// registers the provider there, and returns it.
		// If registration succeeds, the current container is returned.
		Override(provide any, options ...Option) (DI, error)
		// MustOverride is like Override, but panics if provider registration fails in both the current and fallback container.
		// Use it when failure to override should be treated as a fatal error.
		MustOverride(provide any, options ...Option) DI
		// Provide registers a new dependency provider in the container.
		// It returns an error if the given value is not a valid provider function or type.
		Provide(provider any, options ...Option) error
		// MustProvide is like Provide, but panics if an error occurs.
		// Use it when failure to register a provider should be considered a fatal error.
		MustProvide(provider any, options ...Option) DI
		// Invoke resolves and calls the given function using dependencies from the container.
		// Returns an error if any dependency is missing or cannot be constructed.
		Invoke(functions ...any) error
		// MustInvoke is like Invoke, but panics if an error occurs during resolution or invocation.
		// Use it when failing to invoke the function should stop execution immediately.
		MustInvoke(functions ...any) DI
		// InvokeWithDI resolves top-level dependencies using the current container,
		// while resolving their nested dependencies using the provided container.
		// Useful when composing containers or overriding lower-level dependencies.
		InvokeWithDI(di DI, functions ...any) error
	}
	Option func(option)
	option interface {
		SetTransient()
	}
)

// WithTransient marks the dependency as transient (non-singleton).
// This means a new instance will be created each time it is requested from the container.
func WithTransient() Option {
	return func(opt option) { opt.SetTransient() }
}

var (
	ErrDependencyAlreadyExists = errors.New("dependency is already registered")
	ErrProviderWithoutOutputs  = errors.New("provider function has no return types")
	ErrInvalidValueProvided    = errors.New("cannot provide an invalid value")
	ErrNilValueProvided        = errors.New("cannot provide a nil value")
	ErrNotAFunction            = errors.New("expected a function, got a non-function value for invocation")
	ErrDependencyNotFound      = errors.New("missing")
	ErrCyclicDependency        = errors.New("cyclic dependency detected")
)
