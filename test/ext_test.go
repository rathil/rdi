package test

import (
	"errors"
	"testing"

	"github.com/rathil/rdi"
	"github.com/rathil/rdi/standard"
)

func TestInterfaceGet(t *testing.T) {
	type data1 struct {
		A int
	}
	d, err := rdi.Get[data1](
		standard.New().
			MustProvide(data1{11}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if d.A != 11 {
		t.Errorf("got %d, want %d", d.A, 11)
	}
}

func TestStandardGet(t *testing.T) {
	type data1 struct {
		A int
	}
	d, err := standard.Get[data1](
		standard.New().
			MustProvide(data1{11}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if d.A != 11 {
		t.Errorf("got %d, want %d", d.A, 11)
	}
}

func TestStandardMyDIGet(t *testing.T) {
	type data1 struct {
		A int
	}
	d, err := standard.Get[data1](
		(&testDi{standard.New()}).
			MustProvide(data1{11}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if d.A != 11 {
		t.Errorf("got %d, want %d", d.A, 11)
	}
}

func TestInterfaceMustGet(t *testing.T) {
	type data1 struct {
		A int
	}
	d := rdi.MustGet[data1](
		standard.New().
			MustProvide(data1{11}),
	)
	if d.A != 11 {
		t.Errorf("got %d, want %d", d.A, 11)
	}
}

func TestStandardMustGet(t *testing.T) {
	type data1 struct {
		A int
	}
	d := standard.MustGet[data1](
		standard.New().
			MustProvide(data1{11}),
	)
	if d.A != 11 {
		t.Errorf("got %d, want %d", d.A, 11)
	}
}

func TestInterfaceMustGetPanic(t *testing.T) {
	type data1 struct {
		A int
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("expected panic, got nil")
		}
		if !errors.Is(err.(error), rdi.ErrDependencyNotFound) {
			t.Errorf("expected ErrDependencyNotFound panic, got %v", err)
		}
	}()

	rdi.MustGet[data1](standard.New())
}

func TestStandardMustGetPanic(t *testing.T) {
	type data1 struct {
		A int
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("expected panic, got nil")
		}
		if !errors.Is(err.(error), rdi.ErrDependencyNotFound) {
			t.Errorf("expected ErrDependencyNotFound panic, got %v", err)
		}
	}()

	standard.MustGet[data1](standard.New())
}
