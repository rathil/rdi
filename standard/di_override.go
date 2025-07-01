package standard

import (
	"errors"

	"github.com/rathil/rdi"
)

func (a *di) MustOverride(
	provide any,
	options ...rdi.Option,
) rdi.DI {
	if err := a.provide(provide, options...); err != nil {
		if !errors.Is(err, rdi.ErrDependencyAlreadyExists) {
			panic(err)
		}
		newA := newWithParent(a)
		if err = newA.provide(provide, options...); err != nil {
			panic(err)
		}
		return newA
	}
	return a
}

func (a *di) Override(
	provide any,
	options ...rdi.Option,
) (rdi.DI, error) {
	if err := a.provide(provide, options...); err != nil {
		if !errors.Is(err, rdi.ErrDependencyAlreadyExists) {
			return a, err
		}
		newA := newWithParent(a)
		return newA, newA.provide(provide, options...)
	}
	return a, nil
}
