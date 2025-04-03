package standard

import (
	"reflect"

	"github.com/rathil/rdi"
)

func (a *di) MustProvide(
	value any,
	options ...rdi.Option,
) rdi.DI {
	if err := a.Provide(value, options...); err != nil {
		panic(err)
	}
	return a
}

func (a *di) Provide(
	provide any,
	options ...rdi.Option,
) error {
	rv := reflect.ValueOf(provide)
	if !rv.IsValid() || rv.IsZero() {
		return rdi.ErrNilPointerProvided
	}
	if rv.Kind() == reflect.Func {
		return a.declareFunction(rv, options)
	}
	return a.declareValue(rv, options)
}

func (a *di) declareValue(
	rv reflect.Value,
	options []rdi.Option,
) error {
	dep := makeDependence(options)
	dep.cache.Store(&rv)
	if _, loaded := a.storage.LoadOrStore(rv.Type(), dep); loaded {
		return a.errDependencyAlreadyExists(rv.Type())
	}
	return nil
}

func (a *di) declareFunction(
	function reflect.Value,
	options []rdi.Option,
) error {
	rt := function.Type()

	numIn := rt.NumIn()
	in := make([]reflect.Type, 0, numIn)
	for i := range numIn {
		in = append(in, rt.In(i))
	}

	var count int
	for o := range rt.NumOut() {
		out := rt.Out(o)
		if out.String() == "error" {
			continue
		}
		dep := makeDependence(options)
		dep.function = function
		dep.in = in
		dep.outIndex = o
		if _, loaded := a.storage.LoadOrStore(out, dep); loaded {
			return a.errDependencyAlreadyExists(out)
		}
		count++
	}
	if count == 0 {
		return rdi.ErrProviderWithoutOutputs
	}
	return nil
}
