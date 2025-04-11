package standard

import (
	"github.com/rathil/rdi"
)

// Base returns the default global DI container.
// It can be used as a shared root for registering common dependencies.
func Base() rdi.DI { return base }

// New creates a new DI container with the Base container as its parent.
// Dependencies not found in the new container will be resolved from the Base container.
func New() rdi.DI {
	return NewWithParent(base)
}

// NewWithParent creates a new DI container with the given parent container.
// The new container can override or extend dependencies from the parent.
func NewWithParent(parent rdi.DI) rdi.DI {
	d := &di{
		parent: parent,
	}
	return d.MustProvide(func() rdi.DI { return d })
}

var base = NewWithParent(nil)
