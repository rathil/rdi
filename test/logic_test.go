package test

import (
	"errors"
	"strings"
	"testing"

	"github.com/rathil/rdi"
	"github.com/rathil/rdi/standard"
)

func TestSimple(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	if err := standard.New().
		MustProvide(data1{11}).
		MustProvide(data2{22}).
		MustProvide(&data2{33}).
		Invoke(func(d data1, d2 data2, d2p *data2) {
			if d.A != 11 {
				t.Errorf("got %d, want %d", d.A, 11)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
			if d2p.A != 33 {
				t.Errorf("got %d, want %d", d2p.A, 33)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleWithParent(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	if err := standard.NewWithParent(
		standard.New().
			MustProvide(data1{11}).
			MustProvide(data2{22}),
	).
		MustProvide(data1{33}).
		Invoke(func(d1 data1, d2 data2) {
			if d1.A != 33 {
				t.Errorf("got %d, want %d", d1.A, 33)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleOverride(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	if err := standard.New().
		MustOverride(data1{11}).
		MustOverride(data2{22}).
		MustOverride(data1{33}).
		Invoke(func(d1 data1, d2 data2) {
			if d1.A != 33 {
				t.Errorf("got %d, want %d", d1.A, 33)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestSimpleOverrideFail(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("expected panic, got nil")
		}
		if !errors.Is(err.(error), rdi.ErrDependencyAlreadyExists) {
			t.Errorf("expected ErrDependencyNotFound panic, got %v", err)
		}
	}()

	standard.New().
		MustOverride(data1{11}).
		MustOverride(data2{22}).
		MustProvide(data1{33})
}

func TestFunction(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	if err := standard.New().
		MustProvide(func() data1 { return data1{11} }).
		MustProvide(func() data2 { return data2{22} }).
		MustProvide(func(d data2) *data2 { return &data2{d.A * 10} }).
		Invoke(func(d data1, d2 data2, d2p *data2) {
			if d.A != 11 {
				t.Errorf("got %d, want %d", d.A, 11)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
			if d2p.A != 220 {
				t.Errorf("got %d, want %d", d2p.A, 220)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestFunctionWithParent(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	if err := standard.NewWithParent(
		standard.New().
			MustProvide(func() data1 { return data1{11} }).
			MustProvide(func() data2 { return data2{22} }),
	).
		MustProvide(func() data1 { return data1{33} }).
		Invoke(func(d1 data1, d2 data2) {
			if d1.A != 33 {
				t.Errorf("got %d, want %d", d1.A, 33)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestFunctionOverride(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	if err := standard.New().
		MustOverride(func() data1 { return data1{11} }).
		MustOverride(func() data2 { return data2{22} }).
		MustOverride(func() data1 { return data1{33} }).
		Invoke(func(d1 data1, d2 data2) {
			if d1.A != 33 {
				t.Errorf("got %d, want %d", d1.A, 33)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestMixed(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	if err := standard.New().
		MustProvide(data1{11}).
		MustProvide(func() data2 { return data2{22} }).
		Invoke(func(d data1, d2 data2) {
			if d.A != 11 {
				t.Errorf("got %d, want %d", d.A, 11)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestMixedWithParent(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	if err := standard.NewWithParent(
		standard.New().
			MustProvide(data1{11}).
			MustProvide(func() data2 { return data2{22} }),
	).
		MustProvide(func() data1 { return data1{33} }).
		MustProvide(func() *data2 { return &data2{44} }).
		Invoke(func(d2p *data2, d1 data1, d2 data2) {
			if d2p.A != 44 {
				t.Errorf("got %d, want %d", d2p.A, 44)
			}
			if d1.A != 33 {
				t.Errorf("got %d, want %d", d1.A, 33)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestMixedOverride(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	if err := standard.New().
		MustOverride(data1{11}).
		MustOverride(func() data2 { return data2{22} }).
		MustOverride(func() data1 { return data1{33} }).
		MustOverride(func() *data2 { return &data2{44} }).
		Invoke(func(d2p *data2, d1 data1, d2 data2) {
			if d2p.A != 44 {
				t.Errorf("got %d, want %d", d2p.A, 44)
			}
			if d1.A != 33 {
				t.Errorf("got %d, want %d", d1.A, 33)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestMixedWithParentWithNested(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
		D data1
	}
	if err := standard.NewWithParent(
		standard.New().
			MustProvide(data1{11}).
			MustProvide(func(d data1) data2 {
				return data2{
					A: 22,
					D: d,
				}
			}),
	).
		MustProvide(func() data1 { return data1{33} }).
		Invoke(func(d1 data1, d2 data2) {
			if d1.A != 33 {
				t.Errorf("got %d, want %d", d1.A, 33)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
			if d2.D.A != 33 {
				t.Errorf("got %d, want %d", d2.D.A, 33)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestMixedOverrideWithNested(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
		D data1
	}
	if err := standard.New().
		MustOverride(data1{11}).
		MustOverride(func(d data1) data2 {
			return data2{
				A: 22,
				D: d,
			}
		}).
		MustOverride(func() data1 { return data1{33} }).
		Invoke(func(d1 data1, d2 data2) {
			if d1.A != 33 {
				t.Errorf("got %d, want %d", d1.A, 33)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
			if d2.D.A != 33 {
				t.Errorf("got %d, want %d", d2.D.A, 33)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestMixedWithParentWithNestedGreedy(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
		D data1
	}
	if err := standard.NewWithParent(
		standard.New().
			MustProvide(data1{11}).
			MustProvide(func(d data1) *data2 {
				return &data2{
					A: 22,
					D: d,
				}
			}).
			MustInvoke(func(d2 *data2) {
				if d2.A != 22 {
					t.Errorf("got %d, want %d", d2.A, 22)
				}
				if d2.D.A != 11 {
					t.Errorf("got %d, want %d", d2.D.A, 11)
				}
			}),
	).
		MustProvide(func() data1 { return data1{33} }).
		Invoke(func(d1 data1, d2 *data2) {
			if d1.A != 33 {
				t.Errorf("got %d, want %d", d1.A, 33)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
			if d2.D.A != 11 {
				t.Errorf("got %d, want %d", d2.D.A, 11)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestMixedOverrideWithNestedGreedy(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
		D data1
	}
	if err := standard.New().
		MustOverride(data1{11}).
		MustOverride(func(d data1) *data2 {
			return &data2{
				A: 22,
				D: d,
			}
		}).
		MustInvoke(func(d2 *data2) {
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
			if d2.D.A != 11 {
				t.Errorf("got %d, want %d", d2.D.A, 11)
			}
		}).
		MustOverride(func() data1 { return data1{33} }).
		Invoke(func(d1 data1, d2 *data2) {
			if d1.A != 33 {
				t.Errorf("got %d, want %d", d1.A, 33)
			}
			if d2.A != 22 {
				t.Errorf("got %d, want %d", d2.A, 22)
			}
			if d2.D.A != 11 {
				t.Errorf("got %d, want %d", d2.D.A, 11)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestMixedWithParentWithNestedDeep(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	type data3 struct {
		A int
	}
	type data4 struct {
		A int
	}
	type data5 struct {
		A int
	}
	if err := standard.NewWithParent(
		standard.New().
			MustProvide(data1{55}).
			MustProvide(func(d data1) data2 { return data2{d.A} }).
			MustProvide(func(d data3) data4 { return data4{d.A} }).
			MustProvide(func(d data4) data5 { return data5{d.A} }),
	).
		MustProvide(func(d data2) data3 { return data3{d.A * 10} }).
		Invoke(func(d5 data5) {
			if d5.A != 550 {
				t.Errorf("got %d, want %d", d5.A, 550)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestMixedOverrideWithNestedDeep(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	type data3 struct {
		A int
	}
	type data4 struct {
		A int
	}
	type data5 struct {
		A int
	}
	if err := standard.New().
		MustOverride(data1{55}).
		MustOverride(func(d data1) data2 { return data2{d.A} }).
		MustOverride(func(d data3) data4 { return data4{d.A} }).
		MustOverride(func(d data4) data5 { return data5{d.A} }).
		MustOverride(func(d data2) data3 { return data3{d.A * 10} }).
		Invoke(func(d5 data5) {
			if d5.A != 550 {
				t.Errorf("got %d, want %d", d5.A, 550)
			}
		}); err != nil {
		t.Fatal(err)
	}
}

func TestErrorNotFound(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	err := standard.New().
		MustProvide(data1{11}).
		Invoke(func(d2 data2) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrDependencyNotFound) {
		t.Errorf("expected ErrDependencyNotFound, got %v", err)
	}
}

func TestErrorNotFoundDeep(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	err := standard.NewWithParent(
		standard.NewWithParent(
			standard.NewWithParent(
				standard.NewWithParent(
					standard.New(),
				),
			).
				MustProvide(data1{11}),
		),
	).
		Invoke(func(d2 data2) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrDependencyNotFound) {
		t.Errorf("expected ErrDependencyNotFound, got %v", err)
	}
}

func TestErrorAlreadyExistsSimple(t *testing.T) {
	type data1 struct {
		A int
	}
	err := standard.New().
		MustProvide(data1{11}).
		Provide(data1{22})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrDependencyAlreadyExists) {
		t.Errorf("expected ErrDependencyAlreadyExists, got %v", err)
	}
	if err.Error() == "'data1', err: "+rdi.ErrDependencyAlreadyExists.Error() {
		t.Error(err)
	}
}

func TestErrorProviderWithoutOutputs(t *testing.T) {
	type data1 struct {
		A int
	}
	err := standard.New().
		MustProvide(data1{11}).
		Provide(func(d data1) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrProviderWithoutOutputs) {
		t.Errorf("expected ErrProviderWithoutOutputs, got %v", err)
	}
}

func TestErrorProviderWithoutOutputsSimple(t *testing.T) {
	type data1 struct {
		A int
	}
	err := standard.New().
		MustProvide(data1{11}).
		Provide(func() {})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrProviderWithoutOutputs) {
		t.Errorf("expected ErrProviderWithoutOutputs, got %v", err)
	}
}

func TestErrorAlreadyExistsFunction(t *testing.T) {
	type data1 struct {
		A int
	}
	err := standard.New().
		MustProvide(func() data1 { return data1{11} }).
		Provide(func() data1 { return data1{22} })
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrDependencyAlreadyExists) {
		t.Errorf("expected ErrDependencyAlreadyExists, got %v", err)
	}
}

func TestErrorAlreadyExistsMixed1(t *testing.T) {
	type data1 struct {
		A int
	}
	err := standard.New().
		MustProvide(func() data1 { return data1{11} }).
		Provide(data1{22})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrDependencyAlreadyExists) {
		t.Errorf("expected ErrDependencyAlreadyExists, got %v", err)
	}
}

func TestErrorAlreadyExistsMixed2(t *testing.T) {
	type data1 struct {
		A int
	}
	err := standard.New().
		MustProvide(data1{22}).
		Provide(func() data1 { return data1{11} })
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrDependencyAlreadyExists) {
		t.Errorf("expected ErrDependencyAlreadyExists, got %v", err)
	}
}

func TestErrorNilPointerProvided(t *testing.T) {
	type data1 struct {
		A int
	}
	var d *data1
	err := standard.New().
		Provide(d)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrNilValueProvided) {
		t.Errorf("expected ErrNilValueProvided, got %v", err)
	}
}

func TestErrorNilPointerProvidedFn(t *testing.T) {
	type fn func() int
	err := standard.New().
		Provide((fn)(nil))
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrNilValueProvided) {
		t.Errorf("expected ErrNilValueProvided, got %v", err)
	}
}

func TestErrorInvalidValueProvided(t *testing.T) {
	err := standard.New().
		Provide(nil)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrInvalidValueProvided) {
		t.Errorf("expected ErrInvalidValueProvided, got %v", err)
	}
}

func TestErrorNotAFunction(t *testing.T) {
	err := standard.New().
		Invoke(55)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, rdi.ErrNotAFunction) {
		t.Errorf("expected ErrNotAFunction, got %v", err)
	}
}

func TestMustInvokePanic(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
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

	standard.New().
		MustProvide(data1{11}).
		MustInvoke(func(d data2) {})
}

func TestMustOverridePanic(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("expected panic, got nil")
		}
		if !errors.Is(err.(error), rdi.ErrProviderWithoutOutputs) {
			t.Errorf("expected ErrProviderWithoutOutputs panic, got %v", err)
		}
	}()

	standard.New().
		MustOverride(func() {})
}

func TestErrorDeep(t *testing.T) {
	type data1 struct{ A int }
	type data2 struct{ A int }
	type data3 struct{ A int }
	type data4 struct{ A int }

	err := standard.New().
		MustProvide(func(d data2) data1 { return data1{d.A} }).
		MustProvide(func(d data3, _ data4) data2 { return data2{d.A} }).
		MustProvide(func() data3 { return data3{11} }).
		Invoke(func(data1) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	} else {
		msg := err.Error()
		if !strings.Contains(msg, "required by test.data1") ||
			!strings.Contains(msg, "required by test.data2") ||
			!strings.Contains(msg, "test.data4") {
			t.Errorf("not present required by")
		}
	}
	if !errors.Is(err, rdi.ErrDependencyNotFound) {
		t.Errorf("expected %v, got %v", rdi.ErrDependencyNotFound, err)
	}
}

func TestMustProvidePanic(t *testing.T) {
	type data1 struct {
		A int
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Errorf("expected panic, got nil")
		}
		if !errors.Is(err.(error), rdi.ErrDependencyAlreadyExists) {
			t.Errorf("expected ErrDependencyAlreadyExists panic, got %v", err)
		}
	}()

	standard.New().
		MustProvide(data1{11}).
		MustProvide(func() data1 { return data1{22} })
}

func TestInvokeError(t *testing.T) {
	type data1 struct {
		A int
	}
	errNow := errors.New("my invoke error")

	err := standard.New().
		MustProvide(data1{11}).
		Invoke(func(d data1) (int, error) { return 1, errNow })
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, errNow) {
		t.Errorf("expected errNow, got %v", err)
	}
}

func TestProvideError(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	type data3 struct {
		A int
	}
	errNow := errors.New("my provide error")

	err := standard.New().
		MustProvide(data1{11}).
		MustProvide(func() (data2, error) { return data2{22}, errNow }).
		MustProvide(func(data2) data3 { return data3{333} }).
		Invoke(func(d1 data1, d2 data3) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if !errors.Is(err, errNow) {
		t.Errorf("expected errNow, got %v", err)
	}
}

func TestTransient(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		B int
	}
	var i int
	standard.New().
		MustProvide(
			func() data1 {
				i++
				return data1{i}
			},
			rdi.WithTransient(),
		).
		MustProvide(
			func(d data1) data2 { return data2{d.A} },
			rdi.WithTransient(),
		).
		MustInvoke(
			func(d data2) {
				if d.B != 1 {
					t.Errorf("got %d, want %d", d.B, 1)
				}
			},
			func(d data1) {
				if d.A != 2 {
					t.Errorf("got %d, want %d", d.A, 2)
				}
			},
			func(d data2) {
				if d.B != 3 {
					t.Errorf("got %d, want %d", d.B, 3)
				}
			},
		)
}

func TestGetDi(t *testing.T) {
	di1 := standard.New()
	di1.MustInvoke(func(di2 rdi.DI) {
		if di1 != di2 {
			t.Errorf("expected %v, got %v", di1, di2)
		}
	})
}

func TestBaseStart(t *testing.T) {
	standard.Base().
		MustProvide(5)
}

func TestBaseEnd(t *testing.T) {
	standard.Base().
		MustInvoke(func(int) {})
}
