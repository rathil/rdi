package standard

import (
	"reflect"
	"unsafe"

	"github.com/rathil/rdi"
	"github.com/rathil/rdi/internal"
)

func (a *di) resolveWire(
	bdi rdi.DI,
	resolve resolveContext,
) (res_ reflect.Value, ok_ bool, _ error) {
	wireType := resolve.dep
	if wireType.Kind() == reflect.Pointer {
		wireType = wireType.Elem()
		defer func() {
			if ok_ {
				res_ = res_.Addr()
			}
		}()
	}
	if wireType.Kind() != reflect.Struct ||
		wireType.NumField() != 2 ||
		wireType.Field(1).Type != internal.WireType {
		return res_, false, nil
	}
	valueType := wireType.Field(0).Type
	if valueType.Kind() == reflect.Pointer {
		valueType = valueType.Elem()
	}
	if valueType.Kind() != reflect.Struct {
		return res_, false, resolve.wrapError(rdi.ErrWireValueNotStruct)
	}

	result := reflect.New(wireType).Elem()
	value := reflect.New(valueType).Elem()

	for i := range valueType.NumField() {
		field := valueType.Field(i)
		if field.Tag.Get("di") != "-" {
			resolve.dep = field.Type
			resolve.wireField = &field
			fieldValue, err := a.resolveDependence(bdi, resolve)
			if err != nil {
				return res_, false, err
			}
			reflect.NewAt(
				field.Type,
				unsafe.Pointer(value.Field(i).UnsafeAddr()),
			).Elem().Set(fieldValue)
		}
	}

	if f := result.Field(0); f.Kind() == reflect.Pointer {
		f.Set(value.Addr())
	} else {
		f.Set(value)
	}

	return result, true, nil
}
