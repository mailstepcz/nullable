// Package nullable provides a generic wrapper type
// allowing to unmarshal JSON data with null values properly.
package nullable

import (
	"bytes"
	"encoding/json"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/mailstepcz/maybe"
	"github.com/oklog/ulid/v2"
	"github.com/shopspring/decimal"
)

var (
	null                   = []byte("null")
	types                  = make(map[reflect.Type]reflect.Type)
	_     json.Unmarshaler = (*Nullable[int])(nil)
	_     Iface            = Nullable[int]{}
)

func init() {
	registerType[int]()
	registerType[string]()
	registerType[bool]()
	registerType[uuid.UUID]()
	registerType[time.Time]()
	registerType[decimal.Decimal]()
	registerType[ulid.ULID]()
}

func registerType[T any]() {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Struct {
		types[reflect.TypeFor[T]()] = reflect.TypeFor[Struct[T]]()
		types[reflect.TypeFor[[]T]()] = reflect.TypeFor[Struct[[]T]]()
	} else {
		types[reflect.TypeFor[T]()] = reflect.TypeFor[Nullable[T]]()
		types[reflect.TypeFor[maybe.Maybe[T]]()] = reflect.TypeFor[Nullable[T]]()
		types[reflect.TypeFor[[]T]()] = reflect.TypeFor[Slice[T]]()
		types[reflect.TypeFor[[]maybe.Maybe[T]]()] = reflect.TypeFor[Slice[T]]()
	}
}

// Type returns the nullable counterpart of the type.
func Type(typ reflect.Type) (reflect.Type, bool) {
	x, ok := types[typ]
	return x, ok
}

// Iface is the non-generic interface for [Nullable].
type Iface interface {
	Value() interface{}
	IsNonNull() bool
	IsValid() bool
}

// SliceIface is the non-generic interface for [NullableSlice].
type SliceIface interface {
	Iface
	Len() int
}

// Nullable is a nullable wrapper type for JSON unmarshalling.
type Nullable[T any] struct {
	value   T
	nonNull bool
	valid   bool
}

func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
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

func (n Nullable[T]) Value() interface{} { return n.value }

func (n Nullable[T]) IsNonNull() bool { return n.nonNull }

func (n Nullable[T]) IsValid() bool { return n.valid }

// Slice is a nullable wrapper type for JSON unmarshalling.
type Slice[T any] struct {
	value   []T
	nonNull bool
	valid   bool
}

func (n *Slice[T]) UnmarshalJSON(b []byte) error {
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

func (n Slice[T]) Value() interface{} { return n.value }

func (n Slice[T]) IsNonNull() bool { return n.nonNull }

func (n Slice[T]) IsValid() bool { return n.valid }

func (n Slice[T]) Len() int { return len(n.value) }
