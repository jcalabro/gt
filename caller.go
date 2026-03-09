package gt

import (
	"runtime"
	"strings"
	"sync"
)

const unknown = "unknown"

var callerCache sync.Map

// CallerName returns the name of the calling function at the given skip level.
// It caches results by program counter for performance.
func CallerName(skip int) string {
	var pcs [1]uintptr
	if n := runtime.Callers(skip+1, pcs[:]); n == 0 {
		return unknown
	}

	pc := pcs[0]

	if v, ok := callerCache.Load(pc); ok {
		s, ok := v.(string)
		if !ok {
			return unknown
		}

		return s
	}

	frames := runtime.CallersFrames(pcs[:])
	frame, _ := frames.Next()

	caller := frame.Function
	if caller == "" {
		caller = unknown
	}

	if ndx := strings.LastIndex(caller, "."); ndx != -1 {
		caller = caller[ndx+1:]
	}

	callerCache.Store(pc, caller)
	return caller
}
