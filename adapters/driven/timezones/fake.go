package timezones

import (
	"fmt"
	"time"
)

// Fake stands in for System in tests, offering the zones it was given.
type Fake struct {
	local *time.Location
	zones []*time.Location
}

// NewFake returns a fake whose local zone is local and which loads each of
// zones by its name.
func NewFake(local *time.Location, zones []*time.Location) *Fake {
	return &Fake{local: local, zones: zones}
}

// Local returns the local zone f was given.
func (f *Fake) Local() *time.Location {
	return f.local
}

// Load returns the zone f was given with the given name, and fails when it
// was given none.
func (f *Fake) Load(name string) (*time.Location, error) {
	for _, zone := range f.zones {
		if zone.String() == name {
			return zone, nil
		}
	}
	return nil, fmt.Errorf("loading time zone %q: the fake wasn't given it", name)
}
