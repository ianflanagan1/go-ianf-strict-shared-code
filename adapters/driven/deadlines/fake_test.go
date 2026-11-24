package deadlines_test

import (
	"context"
	"testing"
	"time"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/adapters/driven/deadlines"
)

func TestTheFakeHandsOutItsDeadlinesInOrder(t *testing.T) {
	f := deadlines.NewFake([]bool{false, true})
	live, cancelLive := f.After(context.Background(), time.Millisecond)
	defer cancelLive()
	expired, cancelExpired := f.After(context.Background(), time.Hour)
	defer cancelExpired()
	if err := live.Err(); err != nil {
		t.Fatalf("got %v from the first context; want nil", err)
	}
	if err := expired.Err(); err == nil {
		t.Fatal("got nil from the second context; want its deadline passed")
	}
}

func TestTheFakePanicsWhenAskedForMoreDeadlinesThanItWasGiven(t *testing.T) {
	f := deadlines.NewFake([]bool{false})
	_, cancel := f.After(context.Background(), time.Hour)
	defer cancel()
	defer func() {
		if recover() == nil {
			t.Fatal("got a context; want a panic")
		}
	}()
	f.After(context.Background(), time.Hour)
}
