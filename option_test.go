package gt_test

import (
	"encoding/json"
	"testing"

	"github.com/jcalabro/gt"
	"github.com/stretchr/testify/require"
)

func TestOptionType(t *testing.T) {
	{
		val := 123
		opt := gt.Some(val)
		require.True(t, opt.HasVal())
		require.False(t, opt.IsNone())
		require.Equal(t, val, opt.Val())
	}

	{
		opt := gt.None[any]()
		require.False(t, opt.HasVal())
		require.True(t, opt.IsNone())

		recoverCalled := false
		defer func() { require.True(t, recoverCalled) }()
		defer func() {
			if r := recover(); r != nil {
				recoverCalled = true
			}
		}()
		require.Equal(t, 0, opt.Val())
	}
}

func TestOptionJSON_SomeString(t *testing.T) {
	opt := gt.Some("hello")
	data, err := json.Marshal(opt)
	require.NoError(t, err)
	require.Equal(t, `"hello"`, string(data))

	var out gt.Option[string]
	require.NoError(t, json.Unmarshal(data, &out))
	require.True(t, out.HasVal())
	require.Equal(t, "hello", out.Val())
}

func TestOptionJSON_NoneString(t *testing.T) {
	opt := gt.None[string]()
	data, err := json.Marshal(opt)
	require.NoError(t, err)
	require.Equal(t, "null", string(data))

	var out gt.Option[string]
	require.NoError(t, json.Unmarshal(data, &out))
	require.True(t, out.IsNone())
}

func TestOptionJSON_SomeInt(t *testing.T) {
	opt := gt.Some(int64(42))
	data, err := json.Marshal(opt)
	require.NoError(t, err)
	require.Equal(t, "42", string(data))

	var out gt.Option[int64]
	require.NoError(t, json.Unmarshal(data, &out))
	require.True(t, out.HasVal())
	require.Equal(t, int64(42), out.Val())
}

func TestOptionJSON_SomeBool(t *testing.T) {
	opt := gt.Some(false)
	data, err := json.Marshal(opt)
	require.NoError(t, err)
	require.Equal(t, "false", string(data))

	var out gt.Option[bool]
	require.NoError(t, json.Unmarshal(data, &out))
	require.True(t, out.HasVal())
	require.Equal(t, false, out.Val())
}

func TestOptionJSON_InStruct(t *testing.T) {
	type Profile struct {
		Name   string            `json:"name"`
		Avatar gt.Option[string] `json:"avatar,omitempty"`
		Age    gt.Option[int64]  `json:"age,omitempty"`
	}

	// With values.
	p := Profile{Name: "alice", Avatar: gt.Some("pic.jpg"), Age: gt.Some(int64(30))}
	data, err := json.Marshal(p)
	require.NoError(t, err)
	require.Contains(t, string(data), `"avatar":"pic.jpg"`)
	require.Contains(t, string(data), `"age":30`)

	var out Profile
	require.NoError(t, json.Unmarshal(data, &out))
	require.Equal(t, "alice", out.Name)
	require.True(t, out.Avatar.HasVal())
	require.Equal(t, "pic.jpg", out.Avatar.Val())

	// With null values.
	data2 := []byte(`{"name":"bob","avatar":null}`)
	var out2 Profile
	require.NoError(t, json.Unmarshal(data2, &out2))
	require.Equal(t, "bob", out2.Name)
	require.True(t, out2.Avatar.IsNone())
	require.True(t, out2.Age.IsNone())

	// Missing fields (zero value = None).
	data3 := []byte(`{"name":"charlie"}`)
	var out3 Profile
	require.NoError(t, json.Unmarshal(data3, &out3))
	require.True(t, out3.Avatar.IsNone())
	require.True(t, out3.Age.IsNone())
}

func TestOptionJSON_SomeStruct(t *testing.T) {
	type Address struct {
		City  string `json:"city"`
		State string `json:"state"`
	}

	opt := gt.Some(Address{City: "Portland", State: "OR"})
	data, err := json.Marshal(opt)
	require.NoError(t, err)
	require.Equal(t, `{"city":"Portland","state":"OR"}`, string(data))

	var out gt.Option[Address]
	require.NoError(t, json.Unmarshal(data, &out))
	require.True(t, out.HasVal())
	require.Equal(t, "Portland", out.Val().City)
	require.Equal(t, "OR", out.Val().State)

	// None round-trip.
	none := gt.None[Address]()
	data, err = json.Marshal(none)
	require.NoError(t, err)
	require.Equal(t, "null", string(data))

	var out2 gt.Option[Address]
	require.NoError(t, json.Unmarshal(data, &out2))
	require.True(t, out2.IsNone())
}

func TestOptionJSON_NullUnmarshalClearsValue(t *testing.T) {
	opt := gt.Some("hello")
	require.NoError(t, json.Unmarshal([]byte("null"), &opt))
	require.True(t, opt.IsNone())
}

func TestOptionJSON_UnmarshalTypeMismatch(t *testing.T) {
	// String into int Option.
	var opt gt.Option[int64]
	err := json.Unmarshal([]byte(`"not a number"`), &opt)
	require.Error(t, err)
	require.True(t, opt.IsNone())
}

func TestOptionJSON_UnmarshalInvalidJSON(t *testing.T) {
	var opt gt.Option[string]
	err := json.Unmarshal([]byte(`{broken`), &opt)
	require.Error(t, err)
	require.True(t, opt.IsNone())
}

func TestOptionJSON_UnmarshalWrongTypeForStruct(t *testing.T) {
	type Inner struct {
		X int `json:"x"`
	}
	// Array into struct Option.
	var opt gt.Option[Inner]
	err := json.Unmarshal([]byte(`[1,2,3]`), &opt)
	require.Error(t, err)
	require.True(t, opt.IsNone())
}

func TestOptionJSON_UnmarshalBoolIntoString(t *testing.T) {
	var opt gt.Option[string]
	err := json.Unmarshal([]byte(`true`), &opt)
	require.Error(t, err)
	require.True(t, opt.IsNone())
}

func TestOptionJSON_UnmarshalEmptyBytes(t *testing.T) {
	var opt gt.Option[string]
	err := json.Unmarshal([]byte(``), &opt)
	require.Error(t, err)
	require.True(t, opt.IsNone())
}

func TestOptionJSON_SomeSlice(t *testing.T) {
	opt := gt.Some([]string{"a", "b", "c"})
	data, err := json.Marshal(opt)
	require.NoError(t, err)
	require.Equal(t, `["a","b","c"]`, string(data))

	var out gt.Option[[]string]
	require.NoError(t, json.Unmarshal(data, &out))
	require.True(t, out.HasVal())
	require.Equal(t, []string{"a", "b", "c"}, out.Val())

	// None slice.
	var out2 gt.Option[[]string]
	require.NoError(t, json.Unmarshal([]byte("null"), &out2))
	require.True(t, out2.IsNone())
}

func TestOptionJSON_UnmarshalReplacesValue(t *testing.T) {
	opt := gt.Some("first")
	require.NoError(t, json.Unmarshal([]byte(`"second"`), &opt))
	require.True(t, opt.HasVal())
	require.Equal(t, "second", opt.Val())
}
