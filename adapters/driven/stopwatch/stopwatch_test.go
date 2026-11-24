package stopwatch_test

import (
	"testing"
	"time"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/adapters/driven/stopwatch"
)

// source is what both the system stopwatch and its fake offer.
type source interface {
	Start() func() time.Duration
}

// contract is the behaviour every stopwatch shares: a measurement's readings
// are never negative and never go back.
func contract(t *testing.T, s source) {
	t.Helper()
	t.Run("a measurement's readings never go back", func(t *testing.T) {
		elapsed := s.Start()
		first := elapsed()
		second := elapsed()
		if first < 0 {
			t.Fatalf("got %v; want a duration no less than 0", first)
		}
		if second < first {
			t.Fatalf("got %v, then %v; want readings that never go back", first, second)
		}
	})
}

func TestTheSystemStopwatchKeepsTheContract(t *testing.T) {
	contract(t, stopwatch.System{})
}

func TestTheFakeKeepsTheContract(t *testing.T) {
	contract(t, stopwatch.NewFake([]time.Duration{time.Second, 2 * time.Second}))
}

func TestTheSystemStopwatchMeasuresTheTimeThatPassed(t *testing.T) {
	elapsed := stopwatch.System{}.Start()
	time.Sleep(10 * time.Millisecond)
	if got := elapsed(); got < 10*time.Millisecond {
		t.Fatalf("got %v; want at least 10ms", got)
	}
}
