package test

import (
	"strings"
	"testing"

	"github.com/rathil/rdi/standard"
)

func TestTraceLevelNone(t *testing.T) {
	type data1 struct{ A int }

	standard.SetTraceLevel(standard.TraceNone)

	err := standard.New().
		Invoke(func(data1) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	} else {
		msg := err.Error()
		if strings.Contains(msg, "TestTraceLevelNone") ||
			strings.Contains(msg, "trace_test.go") {
			t.Errorf("expected no trace information, got %s", msg)
		}
	}
}

func TestTraceLevelFilePath(t *testing.T) {
	type data1 struct{ A int }

	standard.SetTraceLevel(standard.TraceFilePath)

	err := standard.New().
		Invoke(func(data1) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	} else {
		msg := err.Error()
		if !strings.Contains(msg, "trace_test.go") {
			t.Errorf("expected trace information, got %s", msg)
		}
	}
}

func TestTraceLevelFunctionName(t *testing.T) {
	type data1 struct{ A int }

	standard.SetTraceLevel(standard.TraceFunctionName)

	err := standard.New().
		Invoke(func(data1) {})
	if err == nil {
		t.Errorf("expected error, got nil")
	} else {
		msg := err.Error()
		if !strings.Contains(msg, "TestTraceLevelFunctionName") {
			t.Errorf("expected trace information, got %s", msg)
		}
	}
}
