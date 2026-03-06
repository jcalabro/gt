package gt_test

import (
	"encoding/json"
	"testing"

	"github.com/jcalabro/gt"
	"github.com/stretchr/testify/require"
)

func TestRef_SomeAndNone(t *testing.T) {
	r := gt.SomeRef("hello")
	require.True(t, r.HasVal())
	require.False(t, r.IsNone())
	require.Equal(t, "hello", *r.Val())

	n := gt.NoneRef[string]()
	require.False(t, n.HasVal())
	require.True(t, n.IsNone())
}

func TestRef_ValPanicsWhenNone(t *testing.T) {
	r := gt.NoneRef[int]()
	require.Panics(t, func() { r.Val() })
}

func TestRef_JSON_Some(t *testing.T) {
	r := gt.SomeRef("hello")
	data, err := json.Marshal(r)
	require.NoError(t, err)
	require.Equal(t, `"hello"`, string(data))

	var out gt.Ref[string]
	require.NoError(t, json.Unmarshal(data, &out))
	require.True(t, out.HasVal())
	require.Equal(t, "hello", *out.Val())
}

func TestRef_JSON_None(t *testing.T) {
	r := gt.NoneRef[string]()
	data, err := json.Marshal(r)
	require.NoError(t, err)
	require.Equal(t, "null", string(data))

	var out gt.Ref[string]
	require.NoError(t, json.Unmarshal(data, &out))
	require.True(t, out.IsNone())
}

func TestRef_JSON_RecursiveStruct(t *testing.T) {
	type Node struct {
		Name  string       `json:"name"`
		Child gt.Ref[Node] `json:"child,omitzero"`
	}

	tree := Node{
		Name:  "root",
		Child: gt.SomeRef(Node{Name: "leaf"}),
	}
	data, err := json.Marshal(tree)
	require.NoError(t, err)
	require.Contains(t, string(data), `"name":"root"`)
	require.Contains(t, string(data), `"name":"leaf"`)

	var out Node
	require.NoError(t, json.Unmarshal(data, &out))
	require.Equal(t, "root", out.Name)
	require.True(t, out.Child.HasVal())
	require.Equal(t, "leaf", out.Child.Val().Name)
	require.True(t, out.Child.Val().Child.IsNone())
}

func TestRef_JSON_UnmarshalError(t *testing.T) {
	var r gt.Ref[int64]
	err := json.Unmarshal([]byte(`"not a number"`), &r)
	require.Error(t, err)
	require.True(t, r.IsNone())
}

func TestRef_JSON_NullClearsValue(t *testing.T) {
	r := gt.SomeRef("hello")
	require.NoError(t, json.Unmarshal([]byte("null"), &r))
	require.True(t, r.IsNone())
}
