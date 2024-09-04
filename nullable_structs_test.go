package nullable

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNullableStructs(t *testing.T) {
	type address struct {
		City string `json:"city"`
	}

	type person struct {
		Name    Nullable[string] `json:"name"`
		Age     Nullable[int]    `json:"age"`
		Address Struct[address]  `json:"address"`
	}

	t.Run("no fields", func(t *testing.T) {
		req := require.New(t)

		var p person
		err := json.Unmarshal([]byte(`{}`), &p)
		req.NoError(err)
		req.False(p.Address.IsValid())
	})

	t.Run("null fields", func(t *testing.T) {
		req := require.New(t)

		var p person
		err := json.Unmarshal([]byte(`{"address": null}`), &p)
		req.NoError(err)
		req.True(p.Address.IsValid())
		req.False(p.Address.IsNonNull())
	})

	t.Run("fields with values", func(t *testing.T) {
		req := require.New(t)

		var p person
		err := json.Unmarshal([]byte(`{"address": {"city": "Baile Ghrífín"}}`), &p)
		req.NoError(err)
		req.True(p.Address.IsValid())
		req.True(p.Address.IsNonNull())
		req.Equal(&address{City: "Baile Ghrífín"}, p.Address.Value())
	})
}
