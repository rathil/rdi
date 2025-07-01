package standard

import (
	"reflect"
	"runtime"
)

func makeResolveContext(skip int) resolveContext {
	var resolve resolveContext
	tl := TraceLevel(traceLevel.Load())
	if tl == TraceNone {
		return resolve
	}
	if pc, file, line, ok := runtime.Caller(skip + int(traceNesting.Load())); ok {
		resolve.file = file
		resolve.fileLine = line
		if tl == TraceFunctionName {
			if fn := runtime.FuncForPC(pc); fn != nil {
				resolve.function = fn.Name()
			}
		}
	}
	return resolve
}

type resolveContext struct {
	prev                     *resolveContext
	dep                      reflect.Type
	file                     string
	fileLine                 int
	function                 string
	invokeFunctionIndex      int
	invokeFunctionParamIndex int
}

func (a resolveContext) wrapError(parent error) *Error {
	err := a.makeError()
	err.Parent = parent
	for resolve := &a; resolve.prev != nil; resolve = resolve.prev {
		err.RequiredBy = append(err.RequiredBy, resolve.prev.makeError())
	}
	return &err
}

func (a resolveContext) makeError() Error {
	err := Error{
		File:                     a.file,
		FileLine:                 a.fileLine,
		Function:                 a.function,
		InvokeFunctionIndex:      a.invokeFunctionIndex,
		InvokeFunctionParamIndex: a.invokeFunctionParamIndex,
	}
	if a.dep != nil {
		err.Dependence = a.dep.String()
	}
	return err
}
