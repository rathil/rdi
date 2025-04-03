package test

import (
	"errors"
	"testing"

	"github.com/rathil/rdi"
	"github.com/rathil/rdi/standard"
)

func TestInvokeWithDI(t *testing.T) {
	type data1 struct {
		A string
	}
	di1 := standard.New().
		MustProvide("string")
	err := standard.NewWithParent(di1).
		MustProvide(func(v string) data1 { return data1{v} }).
		InvokeWithDI(di1, func(d data1) {
			if d.A != "string" {
				t.Errorf("got %s, want %s", d.A, "string")
			}
		})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestMyDIProvide(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	standard.NewWithParent(
		(&testDi{standard.New()}).
			MustProvide(data1{11}),
	).
		MustProvide(func(d data1) data2 { return data2{d.A} }).
		MustInvoke(func(d data2) {
			if d.A != 11 {
				t.Errorf("got %d, want %d", d.A, 11)
			}
		})
}

func TestMyDIInvoke(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	standard.NewWithParent(
		(&testDi{standard.New()}).
			MustProvide(func() data1 { return data1{11} }),
	).
		MustProvide(func() data2 { return data2{22} }).
		MustInvoke(func(d2 data2, d1 data1) {
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
			if d1.A != 11 {
				t.Errorf("got %d, want %d", d1.A, 11)
			}
		})
}

func TestMyDILast(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	di := &testDi{
		standard.New().
			MustProvide(func() data1 { return data1{11} }),
	}
	err := di.
		MustProvide(func() data2 { return data2{22} }).
		InvokeWithDI(di, func(d2 data2, d1 data1) {
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
			if d1.A != 11 {
				t.Errorf("got %d, want %d", d1.A, 11)
			}
		})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestMyDILastWithErrorProvide(t *testing.T) {
	type data1 struct {
		A int
	}
	errNow := errors.New("my provide error")
	di := &testDi{
		standard.New().
			MustProvide(func() (data1, error) { return data1{11}, errNow }),
	}
	err := di.
		InvokeWithDI(di, func(d data1) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, errNow) {
		t.Errorf("expected errNow, got %v", err)
	}
}

func TestMyDILastWithErrorInvoke(t *testing.T) {
	type data1 struct {
		A int
	}
	errNow := errors.New("my invoke error")
	di := &testDi{
		standard.New().
			MustProvide(func() (data1, error) { return data1{11}, nil }),
	}
	err := di.
		InvokeWithDI(di, func(d data1) error {
			if d.A != 11 {
				t.Errorf("got %d, want %d", d.A, 11)
			}
			return errNow
		})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, errNow) {
		t.Errorf("expected errNow, got %v", err)
	}
}

func TestMyDILastWithErrorNotFound(t *testing.T) {
	type data1 struct {
		A int
	}
	di := &testDi{
		standard.New(),
	}
	err := di.
		InvokeWithDI(di, func(d data1) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrDependencyNotFound) {
		t.Errorf("expected ErrDependencyNotFound, got %v", err)
	}
}

type testDi struct {
	d rdi.DI
}

func (a *testDi) Provide(p any, o ...rdi.Option) error      { return a.d.Provide(p, o...) }
func (a *testDi) MustProvide(p any, o ...rdi.Option) rdi.DI { a.d.MustProvide(p, o...); return a }
func (a *testDi) Invoke(f ...any) error                     { return a.d.Invoke(f...) }
func (a *testDi) InvokeWithDI(di rdi.DI, f ...any) error    { return a.d.InvokeWithDI(di, f...) }
func (a *testDi) MustInvoke(f ...any) rdi.DI                { a.d.MustInvoke(f...); return a }
