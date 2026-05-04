package test

import (
	"errors"
	"testing"

	"github.com/rathil/rdi"
	"github.com/rathil/rdi/standard"
)

func TestWire(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	type wire struct {
		data1
		d2   data2
		d2p  *data2
		Int  int
		Int2 int
		Int3 int    `di:"-"`
		S    string `di:""`
	}
	if err := standard.New().
		MustProvide(data1{11}).
		MustProvide(data2{22}).
		MustProvide(&data2{33}).
		MustProvide(44).
		MustProvide("55").
		Invoke(func(v rdi.Wire[wire]) {
			if v.Value.data1.A != 11 {
				t.Errorf("got %d, want %d", v.Value.data1.A, 11)
			}
			if v.Value.d2.A != 22 {
				t.Errorf("got %d, want %d", v.Value.d2.A, 22)
			}
			if v.Value.d2p.A != 33 {
				t.Errorf("got %d, want %d", v.Value.d2p.A, 33)
			}
			if v.Value.Int != 44 {
				t.Errorf("got %d, want %d", v.Value.Int, 44)
			}
			if v.Value.Int2 != 44 {
				t.Errorf("got %d, want %d", v.Value.Int2, 44)
			}
			if v.Value.Int3 != 0 {
				t.Errorf("got %d, want %d", v.Value.Int3, 0)
			}
			if v.Value.S != "55" {
				t.Errorf("got %s, want %s", v.Value.S, "55")
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestWirePointer(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	type wire struct {
		data1
		d2   data2
		d2p  *data2
		Int  int
		Int2 int
		Int3 int    `di:"-"`
		S    string `di:""`
	}
	if err := standard.New().
		MustProvide(data1{11}).
		MustProvide(data2{22}).
		MustProvide(&data2{33}).
		MustProvide(44).
		MustProvide("55").
		Invoke(func(v rdi.Wire[*wire]) {
			if v.Value.data1.A != 11 {
				t.Errorf("got %d, want %d", v.Value.data1.A, 11)
			}
			if v.Value.d2.A != 22 {
				t.Errorf("got %d, want %d", v.Value.d2.A, 22)
			}
			if v.Value.d2p.A != 33 {
				t.Errorf("got %d, want %d", v.Value.d2p.A, 33)
			}
			if v.Value.Int != 44 {
				t.Errorf("got %d, want %d", v.Value.Int, 44)
			}
			if v.Value.Int2 != 44 {
				t.Errorf("got %d, want %d", v.Value.Int2, 44)
			}
			if v.Value.Int3 != 0 {
				t.Errorf("got %d, want %d", v.Value.Int3, 0)
			}
			if v.Value.S != "55" {
				t.Errorf("got %s, want %s", v.Value.S, "55")
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestWireFailValue(t *testing.T) {
	type data1 struct {
		A int
	}
	type wire struct {
		D data1
		*data1
	}
	err := standard.New().
		MustProvide(data1{11}).
		Invoke(func(v rdi.Wire[*wire]) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	} else {
		if !errors.Is(err, rdi.ErrDependencyNotFound) {
			t.Errorf("expected ErrDependencyNotFound, got %v", err)
		}
	}
}

func TestWireFailInterface(t *testing.T) {
	type wire interface{}
	err := standard.New().
		Invoke(func(rdi.Wire[wire]) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	} else {
		if !errors.Is(err, rdi.ErrWireValueNotStruct) {
			t.Errorf("expected wire value error, got %s", err.Error())
		}
	}
}
