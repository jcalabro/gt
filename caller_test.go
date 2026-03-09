package gt

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCallerName(t *testing.T) {
	name := CallerName(1)
	require.Equal(t, "TestCallerName", name)
}

func TestCallerNameNested(t *testing.T) {
	name := helper()
	require.Equal(t, "helper", name)
}

func helper() string {
	return CallerName(1)
}

func TestCallerNameCached(t *testing.T) {
	// Call twice to exercise the cache path
	name1 := CallerName(1)
	name2 := CallerName(1)
	require.Equal(t, name1, name2)
	require.Equal(t, "TestCallerNameCached", name1)
}
