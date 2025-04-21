package test

import (
	"context"
	"errors"
	"testing"

	"github.com/rathil/rdi"
	"github.com/rathil/rdi/standard"
)

func TestCycle(t *testing.T) {
	type data1 struct{ A int }

	err := standard.New().
		MustProvide(func(c context.Context) data1 { return data1{} }).
		MustProvide(func(data1) context.Context { return context.Background() }).
		Invoke(func(data1) {})
	if !errors.Is(err.(error), rdi.ErrCyclicDependency) {
		t.Errorf("expected ErrCyclicDependency, got %v", err)
	}
}

func TestCycleDirect(t *testing.T) {
	type data1 struct{ A int }

	err := standard.New().
		MustProvide(func(d data1) data1 { return d }).
		Invoke(func(data1) {})
	if !errors.Is(err.(error), rdi.ErrCyclicDependency) {
		t.Errorf("expected ErrCyclicDependency, got %v", err)
	}
}

func TestCycleOverride(t *testing.T) {
	type data1 struct{ A int }
	type data2 struct{ A int }

	err := standard.NewWithParent(
		standard.New().
			MustProvide(func(d data2) data1 { return data1{d.A} }),
	).
		MustProvide(func(d data1) data2 { return data2{d.A} }).
		Invoke(func(data1) {})
	if !errors.Is(err.(error), rdi.ErrCyclicDependency) {
		t.Errorf("expected ErrCyclicDependency, got %v", err)
	}
}
