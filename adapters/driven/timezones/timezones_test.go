package timezones_test

import (
	"testing"
	"time"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/adapters/driven/timezones"
)

// source is what both the system time zones and its fake offer.
type source interface {
	Local() *time.Location
	Load(name string) (*time.Location, error)
}

// contract is the behaviour every source of zones shares: there is a local
// zone, a zone loads by its name, and a name that is no zone fails.
func contract(t *testing.T, s source) {
	t.Helper()
	t.Run("there is a local zone", func(t *testing.T) {
		if s.Local() == nil {
			t.Fatal("got no local zone; want one")
		}
	})
	t.Run("a zone loads by its name", func(t *testing.T) {
		loc, err := s.Load("Europe/London")
		if err != nil {
			t.Fatalf("loading Europe/London: %v", err)
		}
		if got := loc.String(); got != "Europe/London" {
			t.Fatalf("got %s; want Europe/London", got)
		}
	})
	t.Run("a name that is no zone fails", func(t *testing.T) {
		if _, err := s.Load("Nowhere/Atlantis"); err == nil {
			t.Fatal("got a zone; want an error")
		}
	})
}

func TestTheSystemTimeZonesKeepTheContract(t *testing.T) {
	contract(t, timezones.System{})
}

func TestTheFakeKeepsTheContract(t *testing.T) {
	contract(t, timezones.NewFake(time.UTC, []*time.Location{time.FixedZone("Europe/London", 0)}))
}

func TestTheSystemLocalZoneIsTheMachines(t *testing.T) {
	got := timezones.System{}.Local()
	if got != time.Local {
		t.Fatalf("got %s; want the machine's local zone, %s", got, time.Local)
	}
}
