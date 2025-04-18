package standard

import (
	"reflect"

	"github.com/rathil/rdi"
)

// Get a single dependency
func Get[T any](bdi rdi.DI) (res_ T, _ error) {
	return get[T](bdi)
}

// MustGet is like Get but panics if an error has occurred
func MustGet[T any](bdi rdi.DI) T {
	result, err := get[T](bdi)
	if err != nil {
		panic(err)
	}
	return result
}

func get[T any](bdi rdi.DI) (res_ T, _ error) {
	if d, ok := bdi.(*di); ok {
		resolve := makeResolveContext(3)
		resolve.dep = reflect.TypeFor[T]()
		result, err := d.resolveDependence(d, resolve)
		if err != nil {
			return res_, err
		}
		return result.Interface().(T), nil
	}
	return rdi.Get[T](bdi)
}
