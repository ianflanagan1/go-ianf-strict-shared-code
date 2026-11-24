// Package timezones adapts the machine's time zone settings and the time zone
// database, giving zones the core converts times with.
package timezones

import (
	"fmt"
	"time"
)

// System is the system time zones.
type System struct{}

// Local returns the machine's local zone.
func (System) Local() *time.Location {
	return time.Local
}

// Load returns the zone with the given name, such as "Europe/London", from the
// time zone database.
func (System) Load(name string) (*time.Location, error) {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return nil, fmt.Errorf("loading time zone %q: %w", name, err)
	}
	return loc, nil
}
