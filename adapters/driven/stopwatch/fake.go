package stopwatch

import "time"

// Fake stands in for System in tests, reporting the durations it was given, in
// order.
type Fake struct {
	elapsed []time.Duration
}

// NewFake returns a fake whose readings, from every measurement it starts,
// are elapsed, in order.
func NewFake(elapsed []time.Duration) *Fake {
	return &Fake{elapsed: elapsed}
}

// Start begins a measurement and returns a function that reports the next
// duration f was given each time it is called. That function panics when f
// has none left, since the test gave it too few.
func (f *Fake) Start() func() time.Duration {
	return func() time.Duration {
		if len(f.elapsed) == 0 {
			panic("stopwatch: the fake has no durations left; give NewFake more")
		}
		next := f.elapsed[0]
		f.elapsed = f.elapsed[1:]
		return next
	}
}
