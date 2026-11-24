// Package equalitytest is the shared equality assertion, which only tests use.
package equalitytest

import (
	"testing"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/core/equality"
)

// AssertEqual passes when equality.Equal finds got equal to want, and
// otherwise ends the test with both values.
func AssertEqual[T any](t testing.TB, want, got T) {
	t.Helper()
	if !equality.Equal(want, got) {
		t.Fatalf("got %#v; want %#v", got, want)
	}
}
