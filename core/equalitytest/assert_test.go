package equalitytest_test

import (
	"fmt"
	"testing"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/core/equalitytest"
)

type money struct {
	Pence    int
	Currency string
}

// recorder stands in for a test, recording a failure rather than ending one.
type recorder struct {
	testing.TB
	failure string
}

func (r *recorder) Helper() {}

func (r *recorder) Fatalf(format string, args ...any) { r.failure = fmt.Sprintf(format, args...) }

func TestTheAssertionPassesWhenTheValuesAreEqual(t *testing.T) {
	r := &recorder{TB: t}
	equalitytest.AssertEqual(r, money{100, "GBP"}, money{100, "GBP"})
	if r.failure != "" {
		t.Fatalf("got failure %q; want none", r.failure)
	}
}

func TestTheAssertionFailsWithBothValuesWhenTheyDiffer(t *testing.T) {
	r := &recorder{TB: t}
	equalitytest.AssertEqual(r, money{100, "GBP"}, money{101, "GBP"})
	want := "got equalitytest_test.money{Pence:101, Currency:\"GBP\"}; want equalitytest_test.money{Pence:100, Currency:\"GBP\"}"
	if r.failure != want {
		t.Fatalf("got failure %q; want %q", r.failure, want)
	}
}
