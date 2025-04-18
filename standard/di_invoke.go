package standard

import (
	"reflect"

	"github.com/rathil/rdi"
)

func (a *di) MustInvoke(functions ...any) rdi.DI {
	if err := a.invoke(a, functions...); err != nil {
		panic(err)
	}
	return a
}

func (a *di) Invoke(functions ...any) error {
	return a.invoke(a, functions...)
}

func (a *di) InvokeWithDI(bdi rdi.DI, functions ...any) error {
	return a.invoke(bdi, functions...)
}

func (a *di) invoke(bdi rdi.DI, functions ...any) error {
	resolve := makeResolveContext(3)
	for i, fn := range functions {
		resolve.invokeFunctionIndex = i + 1
		params, err := a.resolveDependencies(bdi, fn, resolve)
		if err != nil {
			return err
		}
		if _, err = a.callFunction(reflect.ValueOf(fn), params); err != nil {
			return err
		}
	}
	return nil
}

func (a *di) resolveDependencies(
	bdi rdi.DI,
	fn any,
	resolve resolveContext,
) ([]reflect.Value, error) {
	rt := reflect.TypeOf(fn)
	if rt == nil || rt.Kind() != reflect.Func {
		return nil, resolve.wrapError(rdi.ErrNotAFunction)
	}
	numIn := rt.NumIn()
	params := make([]reflect.Value, numIn)
	for i := range numIn {
		resolve.dep = rt.In(i)
		resolve.invokeFunctionParamIndex = i + 1
		param, err := a.resolveDependence(bdi, resolve)
		if err != nil {
			return nil, err
		}
		params[i] = param
	}
	return params, nil
}

func (a *di) resolveDependence(
	bdi rdi.DI,
	resolve resolveContext,
) (res_ reflect.Value, _ error) {
	if item, found := a.storage.Load(resolve.dep); found {
		return a.getDependence(bdi, item.(*dependence), resolve)
	}
	if a.parent == nil {
		return res_, resolve.wrapError(rdi.ErrDependencyNotFound)
	}
	if parent, ok := a.parent.(*di); ok {
		return parent.resolveDependence(bdi, resolve)
	}
	var result reflect.Value
	if err := a.parent.InvokeWithDI(
		bdi,
		reflect.MakeFunc(
			reflect.FuncOf([]reflect.Type{resolve.dep}, []reflect.Type{}, false),
			func(in []reflect.Value) []reflect.Value {
				result = in[0]
				return nil
			},
		).Interface(),
	); err != nil {
		return res_, resolve.wrapError(err)
	}
	return result, nil
}

func (a *di) getDependence(
	bdi rdi.DI,
	dep *dependence,
	resolve resolveContext,
) (res_ reflect.Value, _ error) {
	if value := dep.cache.Load(); value != nil {
		return *value, nil
	}
	if !dep.transient {
		dep.locker.Lock()
		defer dep.locker.Unlock()
		if value := dep.cache.Load(); value != nil {
			return *value, nil
		}
	}
	params, err := a.getDependenceParams(bdi, dep, resolve)
	if err != nil {
		return res_, err
	}
	out, err := a.callFunction(dep.function, params)
	if err != nil {
		return res_, err
	}
	value := out[dep.outIndex]
	if !dep.transient {
		dep.cache.Store(&value)
	}
	return value, nil
}

func (a *di) getDependenceParams(
	bdi rdi.DI,
	dep *dependence,
	resolve resolveContext,
) ([]reflect.Value, error) {
	if d, ok := bdi.(*di); ok {
		depResolve := dep.resolve
		depResolve.prev = &resolve
		result := make([]reflect.Value, len(dep.in))
		for i, param := range dep.in {
			depResolve.dep = param
			depResolve.invokeFunctionParamIndex = i + 1
			item, err := d.resolveDependence(bdi, depResolve)
			if err != nil {
				return nil, err
			}
			result[i] = item
		}
		return result, nil
	}
	var result []reflect.Value
	if err := bdi.InvokeWithDI(
		bdi,
		reflect.MakeFunc(
			reflect.FuncOf(dep.in, []reflect.Type{}, false),
			func(out []reflect.Value) []reflect.Value {
				result = out
				return nil
			},
		).Interface(),
	); err != nil {
		return nil, resolve.wrapError(err)
	}
	return result, nil
}

func (a *di) callFunction(
	rv reflect.Value,
	params []reflect.Value,
) ([]reflect.Value, error) {
	results := rv.Call(params)
	for _, result := range results {
		if result.Type().String() == "error" {
			if err, ok := result.Interface().(error); ok {
				return nil, err
			}
		}
	}
	return results, nil
}
