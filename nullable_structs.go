package nullable

import (
	"bytes"
	"encoding/json"
)

var (
	_ json.Unmarshaler = (*Struct[struct{}])(nil)
	_ Iface            = Struct[struct{}]{}
)

// Struct is a nullable wrapper type for JSON unmarshalling suited to structures.
type Struct[T any] struct {
	value   *T
	nonNull bool
	valid   bool
}

func (n *Struct[T]) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, null) {
		n.valid = true
		return nil
	}

	if err := json.Unmarshal(b, &n.value); err != nil {
		return err
	}

	n.nonNull = true
	n.valid = true

	return nil
}

func (n Struct[T]) Value() interface{} { return n.value }

func (n Struct[T]) IsNonNull() bool { return n.nonNull }

func (n Struct[T]) IsValid() bool { return n.valid }
