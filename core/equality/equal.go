// Package equality is the shared equality function: whether two values have
// the same type and hold the same data.
package equality

import "reflect"

// Equal reports whether a and b are equal: the same type and equal fields.
//
// Two structs are equal when they have the same type and each exported field
// of one equals the same field of the other. Two slices or arrays are equal
// when they have the same length and an equal element at each index; a nil
// slice equals an empty one. Two maps are equal when they have the same keys
// with an equal value under each; a key is found as Go's map finds it. Two
// interface values are equal when they hold the same type and equal values.
// Numbers, strings and booleans are equal when == says so, so 0 equals -0 and
// NaN equals nothing, itself included.
//
// Equality is defined on values alone. A pointer or a channel equals only the
// same address. A struct with an unexported field, such as time.Time, keeps
// data this function may not read, and a function can't be told from another
// but by nil, so neither equals anything, itself included: two of them are
// never found equal for want of anything to compare.
//
// It reads exported fields as any code outside the type can, and calls
// nothing the values declare.
func Equal[T any](a, b T) bool {
	return equal(reflect.ValueOf(&a).Elem(), reflect.ValueOf(&b).Elem())
}

// equal compares two values found inside the values Equal was given, whose
// types are not known until they are read.
func equal(a, b reflect.Value) bool {
	if a.Type() != b.Type() {
		return false
	}
	switch a.Kind() {
	case reflect.Struct:
		return equalStructs(a, b)
	case reflect.Slice:
		return a.Len() == b.Len() && equalElements(a, b)
	case reflect.Array:
		return equalElements(a, b)
	case reflect.Map:
		return equalMaps(a, b)
	case reflect.Interface:
		return equalInterfaces(a, b)
	case reflect.Func:
		return a.IsNil() && b.IsNil()
	case reflect.Pointer, reflect.Chan, reflect.UnsafePointer:
		return a.Pointer() == b.Pointer()
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return a.Equal(b)
	case reflect.Invalid:
		return false
	}
	return false
}

// equalStructs compares two structs of one type, field by field, when every
// field is exported.
func equalStructs(a, b reflect.Value) bool {
	for field := range a.Type().Fields() {
		if !field.IsExported() {
			return false
		}
	}
	for field, value := range a.Fields() {
		if !equal(value, b.FieldByIndex(field.Index)) {
			return false
		}
	}
	return true
}

// equalElements compares two slices or arrays of one length, index by index.
func equalElements(a, b reflect.Value) bool {
	for i := range a.Len() {
		if !equal(a.Index(i), b.Index(i)) {
			return false
		}
	}
	return true
}

// equalMaps compares two maps by their keys and the value under each.
func equalMaps(a, b reflect.Value) bool {
	if a.Len() != b.Len() {
		return false
	}
	for key, value := range a.Seq2() {
		other := b.MapIndex(key)
		if !other.IsValid() || !equal(value, other) {
			return false
		}
	}
	return true
}

// equalInterfaces compares what two interface values hold.
func equalInterfaces(a, b reflect.Value) bool {
	if a.IsNil() || b.IsNil() {
		return a.IsNil() && b.IsNil()
	}
	return equal(a.Elem(), b.Elem())
}
