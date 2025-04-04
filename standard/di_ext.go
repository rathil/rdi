package standard

import (
	"reflect"

	"github.com/rathil/rdi"
)

// Get a single dependency
func Get[T any](bdi rdi.DI) (res_ T, _ error) {
	if d, ok := bdi.(*di); ok {
		result, err := d.resolveDependence(d, reflect.TypeFor[T]())
		if err != nil {
			return res_, err
		}
		return result.Interface().(T), nil
	}
	return rdi.Get[T](bdi)
}

// MustGet is like Get but panics if an error has occurred
func MustGet[T any](bdi rdi.DI) T {
	result, err := Get[T](bdi)
	if err != nil {
		panic(err)
	}
	return result
}
