package standard

import "sync/atomic"

// TraceLevel defines the verbosity level for capturing trace information
// during dependency resolution. It controls whether file paths and/or
// function names are recorded in the resolve context.
type TraceLevel uint32

const (
	// TraceNone disables all trace information.
	// No function names or file paths will be displayed.
	// This is the default trace level.
	TraceNone TraceLevel = 0
	// TraceFilePath displays the full file path and line number.
	TraceFilePath TraceLevel = 1
	// TraceFunctionName displays only the name of the function and line number.
	TraceFunctionName TraceLevel = 2
)

// SetTraceLevel sets the global trace level to control how much trace
// information is displayed during dependency resolution.
func SetTraceLevel(level TraceLevel) {
	traceLevel.Store(uint32(level))
}

var traceLevel atomic.Uint32
