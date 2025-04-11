package test

import (
	"github.com/rathil/rdi"
)

type testDi struct {
	d rdi.DI
}

func (a *testDi) Override(p any, o ...rdi.Option) (rdi.DI, error) { return a.d.Override(p, o...) }
func (a *testDi) MustOverride(p any, o ...rdi.Option) rdi.DI      { a.d.MustOverride(p, o...); return a }
func (a *testDi) Provide(p any, o ...rdi.Option) error            { return a.d.Provide(p, o...) }
func (a *testDi) MustProvide(p any, o ...rdi.Option) rdi.DI       { a.d.MustProvide(p, o...); return a }
func (a *testDi) Invoke(f ...any) error                           { return a.d.Invoke(f...) }
func (a *testDi) InvokeWithDI(di rdi.DI, f ...any) error          { return a.d.InvokeWithDI(di, f...) }
func (a *testDi) MustInvoke(f ...any) rdi.DI                      { a.d.MustInvoke(f...); return a }
