package ids_test

import (
	"testing"

	"github.com/ianflanagan1/go-adapters/adapters/driven/ids"
	"github.com/ianflanagan1/go-adapters/values/id"
)

const (
	aRowID          = "01890a5d-ac96-774b-bcce-b302099a8057"
	anotherRowID    = "01890a5d-ac96-774b-bcce-b302099a8058"
	aPublicID       = "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	anotherPublicID = "f47ac10b-58cc-4372-a567-0e02b2c3d47a"
)

func mustRowID(t *testing.T, text string) id.RowID {
	t.Helper()
	r, err := id.ParseRowID(text)
	if err != nil {
		t.Fatalf("parsing %q: %v", text, err)
	}
	return r
}

func mustPublicID(t *testing.T, text string) id.PublicID {
	t.Helper()
	p, err := id.ParsePublicID(text)
	if err != nil {
		t.Fatalf("parsing %q: %v", text, err)
	}
	return p
}

func TestTheFakeHandsOutItsRowIDsInOrder(t *testing.T) {
	f := ids.NewFake([]id.RowID{mustRowID(t, aRowID), mustRowID(t, anotherRowID)}, nil)
	for _, want := range []string{aRowID, anotherRowID} {
		if got := f.NewRowID().Text(); got != want {
			t.Fatalf("got %s; want %s", got, want)
		}
	}
}

func TestTheFakeHandsOutItsPublicIDsInOrder(t *testing.T) {
	f := ids.NewFake(nil, []id.PublicID{mustPublicID(t, aPublicID), mustPublicID(t, anotherPublicID)})
	for _, want := range []string{aPublicID, anotherPublicID} {
		if got := f.NewPublicID().String(); got != want {
			t.Fatalf("got %s; want %s", got, want)
		}
	}
}

func TestTheFakePanicsWhenAskedForMoreRowIDsThanItWasGiven(t *testing.T) {
	f := ids.NewFake([]id.RowID{mustRowID(t, aRowID)}, nil)
	f.NewRowID()
	defer func() {
		if recover() == nil {
			t.Fatal("got a row ID; want a panic")
		}
	}()
	f.NewRowID()
}

func TestTheFakePanicsWhenAskedForMorePublicIDsThanItWasGiven(t *testing.T) {
	f := ids.NewFake(nil, []id.PublicID{mustPublicID(t, aPublicID)})
	f.NewPublicID()
	defer func() {
		if recover() == nil {
			t.Fatal("got a public ID; want a panic")
		}
	}()
	f.NewPublicID()
}
