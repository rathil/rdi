package standard

import (
	"reflect"
	"sync"

	"github.com/rathil/rdi"
)

type di struct {
	storage sync.Map
	parent  rdi.DI
}

func (a *di) errDependencyAlreadyExists(rt reflect.Type) error {
	return &wrapError{
		msg: rt.String(),
		err: rdi.ErrDependencyAlreadyExists,
	}
}
