package deadlines_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/adapters/driven/deadlines"
)

// source is what both the system deadlines and its fake offer.
type source interface {
	After(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc)
}

// contract is the behaviour every source of deadlines shares. It asks s for a
// deadline that passes, then two that don't.
func contract(t *testing.T, s source) {
	t.Helper()
	t.Run("a context is cancelled once its deadline passes", func(t *testing.T) {
		ctx, cancel := s.After(context.Background(), time.Millisecond)
		defer cancel()
		<-ctx.Done()
		if err := ctx.Err(); !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("got %v; want %v", err, context.DeadlineExceeded)
		}
	})
	t.Run("a context is cancelled with the one it came from", func(t *testing.T) {
		parent, cancelParent := context.WithCancel(context.Background())
		ctx, cancel := s.After(parent, time.Hour)
		defer cancel()
		if err := ctx.Err(); err != nil {
			t.Fatalf("got %v before the parent was cancelled; want nil", err)
		}
		cancelParent()
		<-ctx.Done()
		if err := ctx.Err(); !errors.Is(err, context.Canceled) {
			t.Fatalf("got %v; want %v", err, context.Canceled)
		}
	})
	t.Run("a context is cancelled by its cancel function", func(t *testing.T) {
		ctx, cancel := s.After(context.Background(), time.Hour)
		cancel()
		<-ctx.Done()
		if err := ctx.Err(); !errors.Is(err, context.Canceled) {
			t.Fatalf("got %v; want %v", err, context.Canceled)
		}
	})
}

func TestTheSystemDeadlinesKeepTheContract(t *testing.T) {
	contract(t, deadlines.System{})
}

func TestTheFakeKeepsTheContract(t *testing.T) {
	contract(t, deadlines.NewFake([]bool{true, false, false}))
}

func TestTheSystemDeadlineIsTheDurationAfterTheCall(t *testing.T) {
	before := time.Now()
	ctx, cancel := deadlines.System{}.After(context.Background(), time.Hour)
	defer cancel()
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Fatal("got no deadline; want one")
	}
	if deadline.Before(before.Add(time.Hour)) || deadline.After(time.Now().Add(time.Hour)) {
		t.Fatalf("got %v; want an hour after the call", deadline)
	}
}
