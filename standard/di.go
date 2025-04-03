package standard

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/rathil/rdi"
)

type di struct {
	storage sync.Map
	parent  rdi.DI
}

func (a *di) errDependencyAlreadyExists(rt reflect.Type) error {
	return fmt.Errorf("'%s' %w", rt.String(), rdi.ErrDependencyAlreadyExists)
}
