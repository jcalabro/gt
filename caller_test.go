package gt_test

import (
	"strings"
	"testing"

	"github.com/jcalabro/gt"
)

func TestCallerBasic(t *testing.T) {
	info := gt.Caller(1)
	require.Equal(t, "TestCallerBasic", info.Func)
	require.Equal(t, "TestCallerBasic", info.FuncQualified)
	require.Equal(t, "github.com/jcalabro/gt_test", info.Package)
	require.Equal(t, "gt_test", info.PackageShort)
	require.True(t, strings.HasSuffix(info.File, "caller_test.go"))
	require.True(t, info.Line > 0)
	require.True(t, info.PC != 0)
}

func TestCallerNested(t *testing.T) {
	info := callerNestedHelper()
	require.Equal(t, "callerNestedHelper", info.Func)
	require.Equal(t, "callerNestedHelper", info.FuncQualified)
	require.Equal(t, "gt_test", info.PackageShort)
}

func callerNestedHelper() gt.CallerInfo {
	return gt.Caller(1)
}

type fooReceiver struct{}

func (f fooReceiver) ValueMethod() gt.CallerInfo {
	return gt.Caller(1)
}

func (f *fooReceiver) PointerMethod() gt.CallerInfo {
	return gt.Caller(1)
}

func TestCallerValueMethod(t *testing.T) {
	info := fooReceiver{}.ValueMethod()
	require.Equal(t, "ValueMethod", info.Func)
	require.Equal(t, "fooReceiver.ValueMethod", info.FuncQualified)
	require.Equal(t, "gt_test", info.PackageShort)
}

func TestCallerPointerMethod(t *testing.T) {
	f := &fooReceiver{}
	info := f.PointerMethod()
	require.Equal(t, "PointerMethod", info.Func)
	require.Equal(t, "(*fooReceiver).PointerMethod", info.FuncQualified)
	require.Equal(t, "gt_test", info.PackageShort)
}

func TestCallerClosure(t *testing.T) {
	fn := func() gt.CallerInfo { return gt.Caller(1) }
	info := fn()
	require.Equal(t, "func1", info.Func)
	require.Contains(t, info.FuncQualified, "TestCallerClosure")
	require.Contains(t, info.FuncQualified, "func1")
	require.Equal(t, "gt_test", info.PackageShort)
}

func TestCallerFileAndLine(t *testing.T) {
	info := gt.Caller(1)
	require.True(t, strings.HasSuffix(info.File, "caller_test.go"))
	require.True(t, info.Line > 0)
}

func TestCallerCacheHit(t *testing.T) {
	// Call gt.Caller from the same lexical site twice. The PC is identical, so
	// the second call must come from the cache and return an equal value.
	var infos [2]gt.CallerInfo
	for i := range infos {
		infos[i] = gt.Caller(1)
	}
	require.Equal(t, infos[0], infos[1])
	require.Equal(t, infos[0].PC, infos[1].PC)
	require.Equal(t, "TestCallerCacheHit", infos[0].Func)
}

//go:noinline
func cachedCallerHelper() gt.CallerInfo {
	return gt.Caller(1)
}

func TestCallerCacheAcrossCalls(t *testing.T) {
	// Two invocations of the same (non-inlined) helper resolve to the same PC
	// inside the helper, so they must produce identical CallerInfo values.
	a := cachedCallerHelper()
	b := cachedCallerHelper()
	require.Equal(t, a, b)
	require.Equal(t, "cachedCallerHelper", a.Func)
}

func TestCallerUnknown(t *testing.T) {
	info := gt.Caller(10000)
	require.Equal(t, "unknown", info.Func)
	require.Equal(t, "unknown", info.FuncQualified)
	require.Equal(t, "unknown", info.Package)
	require.Equal(t, "unknown", info.PackageShort)
	require.Equal(t, "unknown", info.File)
}
