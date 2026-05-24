package gt

import (
	"runtime"
	"strings"
	"sync"
)

const unknown = "unknown"

// CallerInfo describes a single call site, derived from a program counter.
type CallerInfo struct {
	PC            uintptr
	Package       string
	PackageShort  string
	Func          string
	FuncQualified string
	File          string
	Line          int
}

var callerCache sync.Map

// Caller returns information about the calling function at the given skip level.
// Results are cached by program counter; call sites do not move within a single
// invocation of a program, so the cache is safe and effectively permanent.
func Caller(skip int) CallerInfo {
	var pcs [1]uintptr
	if n := runtime.Callers(skip+1, pcs[:]); n == 0 {
		return unknownCaller()
	}

	pc := pcs[0]
	if v, ok := callerCache.Load(pc); ok {
		if info, ok := v.(CallerInfo); ok {
			return info
		}
	}

	frames := runtime.CallersFrames(pcs[:])
	frame, _ := frames.Next()

	info := buildCallerInfo(pc, frame)
	callerCache.Store(pc, info)
	return info
}

func buildCallerInfo(pc uintptr, frame runtime.Frame) CallerInfo {
	info := CallerInfo{
		PC:   pc,
		File: frame.File,
		Line: frame.Line,
	}
	if info.File == "" {
		info.File = unknown
	}

	pkg, qualified := splitPackageAndFunc(frame.Function)
	info.Package = pkg
	info.PackageShort = shortPackage(pkg)
	info.FuncQualified = qualified
	info.Func = bareFuncName(qualified)
	return info
}

// splitPackageAndFunc separates a runtime function name into its package import
// path and the rest (which includes any receiver qualifier and method name).
//
// Examples:
//
//	"github.com/jcalabro/gt.Foo"                -> ("github.com/jcalabro/gt", "Foo")
//	"github.com/jcalabro/gt.(*Foo).Bar"         -> ("github.com/jcalabro/gt", "(*Foo).Bar")
//	"github.com/jcalabro/gt.Foo.func1"          -> ("github.com/jcalabro/gt", "Foo.func1")
//	"main.main"                                 -> ("main", "main")
func splitPackageAndFunc(fn string) (pkg, qualified string) {
	if fn == "" {
		return unknown, unknown
	}

	// The package boundary is the first '.' AFTER the last '/'. Without a slash,
	// search from the start of the string.
	start := strings.LastIndex(fn, "/") + 1
	dot := strings.Index(fn[start:], ".")
	if dot == -1 {
		return unknown, fn
	}
	split := start + dot
	return fn[:split], fn[split+1:]
}

// shortPackage returns the trailing path segment of a package import path.
func shortPackage(pkg string) string {
	if pkg == "" {
		return unknown
	}
	if i := strings.LastIndex(pkg, "/"); i != -1 {
		return pkg[i+1:]
	}
	return pkg
}

// bareFuncName returns the final dot-separated segment of a qualified function
// name, e.g. "(*Foo).Bar" -> "Bar", "Foo.func1" -> "func1", "Foo" -> "Foo".
func bareFuncName(qualified string) string {
	if qualified == "" {
		return unknown
	}
	if i := strings.LastIndex(qualified, "."); i != -1 {
		return qualified[i+1:]
	}
	return qualified
}

func unknownCaller() CallerInfo {
	return CallerInfo{
		Package:       unknown,
		PackageShort:  unknown,
		Func:          unknown,
		FuncQualified: unknown,
		File:          unknown,
	}
}
