// Package timer drives a program on a schedule, calling a function each time
// a duration passes.
package timer

import (
	"context"
	"time"
)

// System is the system timer.
type System struct{}

// Run calls tick with ctx each time every passes, until ctx is cancelled, and
// then returns. It calls tick once at a time: when a call takes longer than
// every, the ticks that passed meanwhile are dropped, not queued.
func (System) Run(ctx context.Context, every time.Duration, tick func(ctx context.Context)) {
	ticker := time.NewTicker(every)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tick(ctx)
		}
	}
}
