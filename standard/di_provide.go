package standard

import (
	"reflect"

	"github.com/rathil/rdi"
)

func (a *di) MustProvide(
	provide any,
	options ...rdi.Option,
) rdi.DI {
	if err := a.provide(provide, options...); err != nil {
		panic(err)
	}
	return a
}

func (a *di) Provide(
	provide any,
	options ...rdi.Option,
) error {
	return a.provide(provide, options...)
}

func (a *di) provide(
	provide any,
	options ...rdi.Option,
) error {
	resolve := makeResolveContext(3)
	rv := reflect.ValueOf(provide)
	if !rv.IsValid() {
		return resolve.wrapError(rdi.ErrInvalidValueProvided)
	}
	if a.isNilableKind(rv.Kind()) && rv.IsNil() {
		return resolve.wrapError(rdi.ErrNilValueProvided)
	}
	if rv.Kind() == reflect.Func {
		return a.declareFunction(rv, options, resolve)
	}
	return a.declareValue(rv, options, resolve)
}

func (a *di) declareValue(
	rv reflect.Value,
	options []rdi.Option,
	resolve resolveContext,
) error {
	dep := &dependence{
		resolve: resolve,
	}
	dep.applyOptions(options)
	dep.cache.Store(&rv)
	if _, loaded := a.storage.LoadOrStore(rv.Type(), dep); loaded {
		resolve.dep = rv.Type()
		return resolve.wrapError(rdi.ErrDependencyAlreadyExists)
	}
	return nil
}

func (a *di) declareFunction(
	fn reflect.Value,
	options []rdi.Option,
	resolve resolveContext,
) error {
	rt := fn.Type()

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
		dep := &dependence{
			resolve:  resolve,
			function: fn,
			in:       in,
			outIndex: o,
		}
		dep.applyOptions(options)
		if _, loaded := a.storage.LoadOrStore(out, dep); loaded {
			resolve.dep = out
			return resolve.wrapError(rdi.ErrDependencyAlreadyExists)
		}
		count++
	}
	if count == 0 {
		return resolve.wrapError(rdi.ErrProviderWithoutOutputs)
	}
	return nil
}

func (a *di) isNilableKind(rk reflect.Kind) bool {
	switch rk {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return true
	default:
		return false
	}
}
