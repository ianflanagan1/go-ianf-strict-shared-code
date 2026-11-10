package id

import "log/slog"

// placeholder is what every printed form of an Internal value shows.
const placeholder = "[internal]"

// Internal holds a value that must stay inside the program. Its text, debug,
// serialised and logged forms are a placeholder that holds nothing of the
// value; only Reveal gives the value back.
//
// The value is held in a closure, not a field, because fmt prints a struct's
// unexported fields by reflection without calling their methods: a core type
// with an unexported Internal field, printed with %v, would otherwise show the
// value. Reflection can't see into a closure, so fmt shows the placeholder
// field and the closure's code address. As a result, Internal values can't be
// compared with == or reflect.DeepEqual; compare what Reveal returns.
type Internal[T any] struct {
	shown  string
	reveal func() T
}

// NewInternal wraps value.
func NewInternal[T any](value T) Internal[T] {
	return Internal[T]{shown: placeholder, reveal: func() T { return value }}
}

// Reveal returns the value i wraps.
func (i Internal[T]) Reveal() T { return i.reveal() }

// String returns the placeholder.
func (Internal[T]) String() string { return placeholder }

// GoString returns the placeholder.
func (Internal[T]) GoString() string { return placeholder }

// MarshalJSON returns the placeholder as a JSON string.
func (Internal[T]) MarshalJSON() ([]byte, error) { return []byte(`"` + placeholder + `"`), nil }

// LogValue returns the placeholder.
func (Internal[T]) LogValue() slog.Value { return slog.StringValue(placeholder) }
