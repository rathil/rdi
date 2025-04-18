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
	resolve   resolveContext
}

func (a *dependence) applyOptions(options []rdi.Option) {
	for _, item := range options {
		item(a)
	}
}

func (a *dependence) SetTransient() { a.transient = true }
