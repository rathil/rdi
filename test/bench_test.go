package test

import (
	"testing"

	"github.com/rathil/rdi"
	"github.com/rathil/rdi/standard"
)

func BenchmarkStandardGet(b *testing.B) {
	standard.SetTraceLevel(standard.TraceNone)
	type data1 struct {
		A int
	}
	di := standard.New().
		MustProvide(data1{11})
	for i := 0; i < b.N; i++ {
		standard.MustGet[data1](di)
	}
}

func BenchmarkInterfaceGet(b *testing.B) {
	standard.SetTraceLevel(standard.TraceNone)
	type data1 struct {
		A int
	}
	di := standard.New().
		MustProvide(data1{11})
	for i := 0; i < b.N; i++ {
		rdi.MustGet[data1](di)
	}
}
