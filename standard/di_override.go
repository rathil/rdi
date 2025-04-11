package standard

import (
	"errors"

	"github.com/rathil/rdi"
)

func (a *di) MustOverride(
	provide any,
	options ...rdi.Option,
) rdi.DI {
	newDi, err := a.Override(provide, options...)
	if err != nil {
		panic(err)
	}
	return newDi
}

func (a *di) Override(
	provide any,
	options ...rdi.Option,
) (rdi.DI, error) {
	if err := a.Provide(provide, options...); err != nil {
		if !errors.Is(err, rdi.ErrDependencyAlreadyExists) {
			return a, err
		}
		return NewWithParent(a).Override(provide, options...)
	}
	return a, nil
}
