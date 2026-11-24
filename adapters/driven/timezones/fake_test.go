package timezones_test

import (
	"testing"
	"time"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/adapters/driven/timezones"
)

func TestTheFakesLocalZoneIsTheOneItWasGiven(t *testing.T) {
	paris := time.FixedZone("Europe/Paris", 3600)
	if got := timezones.NewFake(paris, nil).Local(); got != paris {
		t.Fatalf("got %s; want Europe/Paris", got)
	}
}

func TestTheFakeLoadsEachZoneItWasGiven(t *testing.T) {
	london, paris := time.FixedZone("Europe/London", 0), time.FixedZone("Europe/Paris", 3600)
	f := timezones.NewFake(time.UTC, []*time.Location{london, paris})
	for _, want := range []*time.Location{london, paris} {
		got, err := f.Load(want.String())
		if err != nil {
			t.Fatalf("loading %s: %v", want, err)
		}
		if got != want {
			t.Fatalf("got %s; want %s", got, want)
		}
	}
}
