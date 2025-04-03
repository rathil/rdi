package test

import (
	"sync"
	"testing"

	"github.com/rathil/rdi"
	"github.com/rathil/rdi/standard"
)

func TestAsync(t *testing.T) {
	type data1 struct {
		A int
	}
	type data2 struct {
		A int
	}
	value := 10
	di := standard.New().
		MustProvide(
			func() data1 {
				value++
				return data1{value}
			},
			rdi.WithTransient(),
		).
		MustProvide(func(d data1) data2 { return data2{d.A} })

	var wg sync.WaitGroup
	wg.Add(3)

	get := func() {
		di.MustInvoke(func(d data2) {
			if d.A != 11 {
				t.Errorf("got %d, want %d", d.A, 11)
			}
		})
		wg.Done()
	}
	go get()
	go get()
	go get()
	wg.Wait()
}
