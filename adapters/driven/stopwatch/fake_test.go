package stopwatch_test

import (
	"testing"
	"time"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/adapters/driven/stopwatch"
)

func TestTheFakeReportsItsDurationsInOrder(t *testing.T) {
	f := stopwatch.NewFake([]time.Duration{time.Second, 3 * time.Second})
	elapsed := f.Start()
	for _, want := range []time.Duration{time.Second, 3 * time.Second} {
		if got := elapsed(); got != want {
			t.Fatalf("got %v; want %v", got, want)
		}
	}
}

func TestTheFakeSharesItsDurationsAcrossMeasurements(t *testing.T) {
	f := stopwatch.NewFake([]time.Duration{time.Second, 3 * time.Second})
	first, second := f.Start(), f.Start()
	if got := second(); got != time.Second {
		t.Fatalf("got %v from the second measurement; want 1s", got)
	}
	if got := first(); got != 3*time.Second {
		t.Fatalf("got %v from the first measurement; want 3s", got)
	}
}

func TestTheFakePanicsWhenReadMoreTimesThanItHasDurations(t *testing.T) {
	elapsed := stopwatch.NewFake([]time.Duration{time.Second}).Start()
	elapsed()
	defer func() {
		if recover() == nil {
			t.Fatal("got a duration; want a panic")
		}
	}()
	elapsed()
}
