package rdi

import "github.com/rathil/rdi/internal"

// Wire requests automatic construction and field injection for T.
//
// Invoke, MustInvoke, and InvokeWithDI recognize this wrapper by its private
// marker field and pass the initialized value through Value. T must be either
// a struct type or a pointer to a struct type.
type Wire[T any] struct {
	Value     T
	_wireType internal.Wire
}
