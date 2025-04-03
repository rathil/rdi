package standard

import (
	"github.com/rathil/rdi"
)

// Base (default) DI container
func Base() rdi.DI { return base }

// New DI container with Base parent container
func New() rdi.DI {
	return NewWithParent(base)
}

// NewWithParent - new DI container with self parent container
func NewWithParent(parent rdi.DI) rdi.DI {
	d := &di{
		parent: parent,
	}
	return d.MustProvide(func() rdi.DI { return d })
}

var base = NewWithParent(nil)
