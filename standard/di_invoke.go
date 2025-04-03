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
	// TODO rathil add file and line to err?
	for _, function := range functions {
		params, err := a.resolveDependencies(bdi, function)
		if err != nil {
			return err
		}
		if _, err = a.callFunction(reflect.ValueOf(function), params); err != nil {
			return err
		}
	}
	return nil
}

func (a *di) resolveDependencies(
	bdi rdi.DI,
	function any,
) ([]reflect.Value, error) {
	rt := reflect.TypeOf(function)
	if rt == nil || rt.Kind() != reflect.Func {
		return nil, rdi.ErrNotAFunction
	}
	numIn := rt.NumIn()
	params := make([]reflect.Value, numIn)
	for i := range numIn {
		param, err := a.resolveDependence(bdi, rt.In(i))
		if err != nil {
			return nil, err
		}
		params[i] = param
	}
	return params, nil
}

func (a *di) resolveDependence(
	bdi rdi.DI,
	rt reflect.Type,
) (res_ reflect.Value, _ error) {
	if item, found := a.storage.Load(rt); found {
		return a.getDependence(bdi, item.(*dependence))
	}
	if a.parent == nil {
		return res_, rdi.ErrDependencyNotFound
	}
	if parent, ok := a.parent.(*di); ok {
		return parent.resolveDependence(bdi, rt)
	}
	var result reflect.Value
	err := a.parent.InvokeWithDI(
		bdi,
		reflect.MakeFunc(
			reflect.FuncOf([]reflect.Type{rt}, []reflect.Type{}, false),
			func(in []reflect.Value) []reflect.Value {
				result = in[0]
				return nil
			},
		).Interface(),
	)
	return result, err
}

func (a *di) getDependence(
	bdi rdi.DI,
	dep *dependence,
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

	params, err := a.getDependenceParams(bdi, dep.in)
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
	params []reflect.Type,
) ([]reflect.Value, error) {
	if d, ok := bdi.(*di); ok {
		result := make([]reflect.Value, len(params))
		for i, param := range params {
			item, err := d.resolveDependence(bdi, param)
			if err != nil {
				return nil, err
			}
			result[i] = item
		}
		return result, nil
	}
	var result []reflect.Value
	err := bdi.InvokeWithDI(
		bdi,
		reflect.MakeFunc(
			reflect.FuncOf(params, []reflect.Type{}, false),
			func(out []reflect.Value) []reflect.Value {
				result = out
				return nil
			},
		).Interface(),
	)
	return result, err
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
