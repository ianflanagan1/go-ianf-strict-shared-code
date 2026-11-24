// Package deadlines adapts the system's timers to give work a deadline.
package deadlines

import (
	"context"
	"time"
)

// System is the system deadlines.
type System struct{}

// After returns a context derived from ctx that is cancelled once d passes,
// and a function that cancels it sooner. Call that function once the work is
// done, to release the timer.
func (System) After(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, d)
}
