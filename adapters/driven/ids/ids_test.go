package ids_test

import (
	"testing"
	"time"

	"github.com/ianflanagan1/go-adapters/adapters/driven/clock"
	"github.com/ianflanagan1/go-adapters/adapters/driven/ids"
	"github.com/ianflanagan1/go-adapters/adapters/driven/random"
	"github.com/ianflanagan1/go-adapters/values/id"
)

// source is what both the system IDs and its fake offer.
type source interface {
	NewRowID() id.RowID
	NewPublicID() id.PublicID
}

// contract is the behaviour every source of IDs shares: each call gives an ID
// no earlier call gave.
func contract(t *testing.T, s source) {
	t.Helper()
	t.Run("each row ID is new", func(t *testing.T) {
		first, second := s.NewRowID(), s.NewRowID()
		if first.Bytes() == second.Bytes() {
			t.Fatalf("got %s twice; want two different row IDs", first.Text())
		}
	})
	t.Run("each public ID is new", func(t *testing.T) {
		first, second := s.NewPublicID(), s.NewPublicID()
		if first.Bytes() == second.Bytes() {
			t.Fatalf("got %s twice; want two different public IDs", first)
		}
	})
}

func TestTheSystemIDsKeepTheContract(t *testing.T) {
	contract(t, ids.NewSystem(clock.System{}, random.Secure{}))
}

func TestTheFakeKeepsTheContract(t *testing.T) {
	contract(t, ids.NewFake(
		[]id.RowID{mustRowID(t, aRowID), mustRowID(t, anotherRowID)},
		[]id.PublicID{mustPublicID(t, aPublicID), mustPublicID(t, anotherPublicID)},
	))
}

func TestARowIDStartsWithTheMillisecondItWasMadeIn(t *testing.T) {
	before := time.Now().UnixMilli()
	b := ids.NewSystem(clock.System{}, random.Secure{}).NewRowID().Bytes()
	after := time.Now().UnixMilli()

	var ms int64
	for _, octet := range b[:6] {
		ms = ms<<8 | int64(octet)
	}
	if ms < before || ms > after {
		t.Fatalf("got %d; want a millisecond from %d to %d", ms, before, after)
	}
}

func TestRowIDsMadeInTheSameMillisecondDiffer(t *testing.T) {
	s := ids.NewSystem(clock.System{}, random.Secure{})
	seen := map[[16]byte]bool{}
	for range 1000 {
		b := s.NewRowID().Bytes()
		if seen[b] {
			t.Fatalf("got %x twice; want 1000 different row IDs", b)
		}
		seen[b] = true
	}
}

func TestPublicIDsDiffer(t *testing.T) {
	s := ids.NewSystem(clock.System{}, random.Secure{})
	seen := map[[16]byte]bool{}
	for range 1000 {
		b := s.NewPublicID().Bytes()
		if seen[b] {
			t.Fatalf("got %x twice; want 1000 different public IDs", b)
		}
		seen[b] = true
	}
}
