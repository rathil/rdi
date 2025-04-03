package rdi

import "errors"

type (
	DI interface {
		// Provide - add provider in container or return error if the value can't be represented as provider.
		Provide(provider any, options ...Option) error
		// MustProvide is like Provide but panics if method Provide return error.
		MustProvide(provider any, options ...Option) DI
		// Invoke - receive dependencies from container.
		Invoke(functions ...any) error
		// InvokeWithDI - receive top-level dependencies using the current container,
		// but receives their nested dependencies using the provided container.
		InvokeWithDI(di DI, functions ...any) error
		// MustInvoke is like Invoke but panics if method Invoke return error.
		MustInvoke(functions ...any) DI
	}
	Option func(option)
	option interface {
		SetTransient()
	}
)

// WithTransient marks this dependency as transient (non-singleton),
// so it must be resolved from its original source each time it is requested
func WithTransient() Option {
	return func(opt option) { opt.SetTransient() }
}

var (
	ErrDependencyAlreadyExists = errors.New("dependency already registered")
	ErrProviderWithoutOutputs  = errors.New("provider function declared without any provided types")
	ErrNilPointerProvided      = errors.New("cannot provide a nil pointer value")
	ErrNotAFunction            = errors.New("for invoke expected a function but received a non-function value")
	ErrDependencyNotFound      = errors.New("requested dependency not found")
)
