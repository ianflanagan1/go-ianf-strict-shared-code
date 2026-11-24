package timer_test

import (
	"context"
	"testing"
	"time"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/adapters/driving/timer"
)

func TestRunCallsTheFunctionEachTimeTheDurationPasses(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	calls := 0
	timer.System{}.Run(ctx, time.Millisecond, func(context.Context) {
		calls++
		if calls == 3 {
			cancel()
		}
	})
	if calls != 3 {
		t.Fatalf("got %d calls; want 3", calls)
	}
}

func TestRunWaitsTheDurationBeforeTheFirstCall(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	start := time.Now()
	var first time.Duration
	timer.System{}.Run(ctx, 20*time.Millisecond, func(context.Context) {
		first = time.Since(start)
		cancel()
	})
	if first < 20*time.Millisecond {
		t.Fatalf("got the first call after %v; want at least 20ms", first)
	}
}

func TestRunReturnsWithoutCallingOnceStopped(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	timer.System{}.Run(ctx, time.Millisecond, func(context.Context) {
		t.Fatal("got a call; want none once stopped")
	})
}

func TestRunPassesItsContextToTheFunction(t *testing.T) {
	type key struct{}
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), key{}, "given"))
	defer cancel()
	var got any
	timer.System{}.Run(ctx, time.Millisecond, func(ctx context.Context) {
		got = ctx.Value(key{})
		cancel()
	})
	if got != "given" {
		t.Fatalf("got %v; want the context Run was given", got)
	}
}
