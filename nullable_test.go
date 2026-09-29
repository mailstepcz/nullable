package nullable

import (
	"encoding/json"
	"reflect"
	"testing"
	stduuid "uuid"

	"github.com/google/uuid"
	"github.com/mailstepcz/maybe"
	"github.com/stretchr/testify/require"
)

func TestMaybeNullable(t *testing.T) {
	req := require.New(t)

	typ, ok := Type(reflect.TypeFor[maybe.Maybe[string]]())
	req.True(ok)
	req.Equal(reflect.TypeFor[Nullable[string]](), typ)
}

func TestUUIDTypes(t *testing.T) {
	t.Run("google uuid", func(t *testing.T) {
		req := require.New(t)

		typ, ok := Type(reflect.TypeFor[uuid.UUID]())
		req.True(ok)
		req.Equal(reflect.TypeFor[Nullable[uuid.UUID]](), typ)

		typ, ok = Type(reflect.TypeFor[[]uuid.UUID]())
		req.True(ok)
		req.Equal(reflect.TypeFor[Slice[uuid.UUID]](), typ)
	})

	t.Run("stdlib uuid", func(t *testing.T) {
		req := require.New(t)

		typ, ok := Type(reflect.TypeFor[stduuid.UUID]())
		req.True(ok)
		req.Equal(reflect.TypeFor[Nullable[stduuid.UUID]](), typ)

		typ, ok = Type(reflect.TypeFor[maybe.Maybe[stduuid.UUID]]())
		req.True(ok)
		req.Equal(reflect.TypeFor[Nullable[stduuid.UUID]](), typ)

		typ, ok = Type(reflect.TypeFor[[]stduuid.UUID]())
		req.True(ok)
		req.Equal(reflect.TypeFor[Slice[stduuid.UUID]](), typ)
	})

	t.Run("stdlib uuid unmarshal", func(t *testing.T) {
		req := require.New(t)

		var v struct {
			ID  Nullable[stduuid.UUID] `json:"id"`
			IDs Slice[stduuid.UUID]    `json:"ids"`
		}
		err := json.Unmarshal([]byte(`{"id": "f81d4fae-7dec-11d0-a765-00a0c91e6bf6", "ids": ["f81d4fae-7dec-11d0-a765-00a0c91e6bf6"]}`), &v)
		req.NoError(err)
		req.True(v.ID.IsNonNull())
		req.Equal(stduuid.MustParse("f81d4fae-7dec-11d0-a765-00a0c91e6bf6"), v.ID.Value())
		req.Equal(1, v.IDs.Len())
	})
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
