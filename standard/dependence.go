package standard

import (
	"reflect"
	"sync"
	"sync/atomic"

	"github.com/rathil/rdi"
)

type dependence struct {
	transient bool
	cache     atomic.Pointer[reflect.Value]
	locker    sync.Mutex
	function  reflect.Value
	in        []reflect.Type
	outIndex  int
}

func makeDependence(options []rdi.Option) *dependence {
	dep := &dependence{}
	for _, item := range options {
		item(dep)
	}
	return dep
}

func (a *dependence) SetTransient() { a.transient = true }
