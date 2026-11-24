package deadlines

import (
	"context"
	"time"
)

// Fake stands in for System in tests, handing out contexts whose deadlines
// have passed or never pass, in the order it was given.
type Fake struct {
	expired []bool
}

// NewFake returns a fake whose next context's deadline has already passed
// for each true in expired, and never passes for each false, in order.
func NewFake(expired []bool) *Fake {
	return &Fake{expired: expired}
}

// After returns a context derived from ctx, ignoring d: its deadline has
// already passed when the next value f was given is true, and never passes
// when it is false. It panics when f has no values left, since the test gave
// it too few.
func (f *Fake) After(ctx context.Context, _ time.Duration) (context.Context, context.CancelFunc) {
	if len(f.expired) == 0 {
		panic("deadlines: the fake has no deadlines left; give NewFake more")
	}
	next := f.expired[0]
	f.expired = f.expired[1:]
	if next {
		return context.WithDeadline(ctx, time.Unix(0, 0))
	}
	return context.WithCancel(ctx)
}
