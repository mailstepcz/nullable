package nullable

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/mailstepcz/maybe"
	"github.com/stretchr/testify/require"
)

func TestMaybeNullable(t *testing.T) {
	req := require.New(t)

	typ, ok := Type(reflect.TypeFor[maybe.Maybe[string]]())
	req.True(ok)
	req.Equal(reflect.TypeFor[Nullable[string]](), typ)
}

func TestNullableFields(t *testing.T) {
	type person struct {
		Name Nullable[string] `json:"name"`
		Age  Nullable[int]    `json:"age"`
	}

	t.Run("no fields", func(t *testing.T) {
		req := require.New(t)

		var p person
		err := json.Unmarshal([]byte(`{}`), &p)
		req.NoError(err)
		req.False(p.Name.IsValid())
		req.False(p.Age.IsValid())
	})

	t.Run("null fields", func(t *testing.T) {
		req := require.New(t)

		var p person
		err := json.Unmarshal([]byte(`{"name": null, "age": null}`), &p)
		req.NoError(err)
		req.True(p.Name.IsValid())
		req.True(p.Age.IsValid())
		req.False(p.Name.IsNonNull())
		req.False(p.Age.IsNonNull())
	})

	t.Run("fields with values", func(t *testing.T) {
		req := require.New(t)

		var p person
		err := json.Unmarshal([]byte(`{"name": "NAME", "age": 25}`), &p)
		req.NoError(err)
		req.True(p.Name.IsValid())
		req.True(p.Age.IsValid())
		req.True(p.Name.IsNonNull())
		req.True(p.Age.IsNonNull())
		req.Equal("NAME", p.Name.Value())
		req.Equal(25, p.Age.Value())
	})
}
