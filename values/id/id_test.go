package id_test

import (
	"errors"
	"testing"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/values/id"
)

const (
	aUUIDv7 = "01890a5d-ac96-774b-bcce-b302099a8057"
	aUUIDv4 = "f47ac10b-58cc-4372-a567-0e02b2c3d479"
)

// aUUIDv7Bytes is aUUIDv7 as 16 bytes.
func aUUIDv7Bytes() []byte {
	return []byte{0x01, 0x89, 0x0a, 0x5d, 0xac, 0x96, 0x77, 0x4b, 0xbc, 0xce, 0xb3, 0x02, 0x09, 0x9a, 0x80, 0x57}
}

// aUUIDv4Bytes is aUUIDv4 as 16 bytes.
func aUUIDv4Bytes() []byte {
	return []byte{0xf4, 0x7a, 0xc1, 0x0b, 0x58, 0xcc, 0x43, 0x72, 0xa5, 0x67, 0x0e, 0x02, 0xb2, 0xc3, 0xd4, 0x79}
}

// notWellFormedText is text that is no well-formed UUID of any version.
func notWellFormedText() map[string]string {
	return map[string]string{
		"empty":                     "",
		"not hexadecimal":           "0189zz5d-ac96-774b-bcce-b302099a8057",
		"too short":                 "01890a5d-ac96-774b-bcce-b302099a805",
		"without hyphens":           "01890a5dac96774bbcceb302099a8057",
		"in braces":                 "{01890a5d-ac96-774b-bcce-b302099a8057}",
		"as a URN":                  "urn:uuid:01890a5d-ac96-774b-bcce-b302099a8057",
		"with surrounding space":    " 01890a5d-ac96-774b-bcce-b302099a8057",
		"the nil UUID":              "00000000-0000-0000-0000-000000000000",
		"the max UUID":              "ffffffff-ffff-ffff-ffff-ffffffffffff",
		"a v7 of the wrong variant": "01890a5d-ac96-774b-ccce-b302099a8057",
		"a v4 of the wrong variant": "f47ac10b-58cc-4372-0567-0e02b2c3d479",
	}
}

func TestARowIDParsesAUUIDv7AndGivesItsTextBack(t *testing.T) {
	r, err := id.ParseRowID(aUUIDv7)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if r.Text() != aUUIDv7 {
		t.Fatalf("got %q; want %q", r.Text(), aUUIDv7)
	}
}

func TestARowIDParsesUpperCaseTextAndGivesItBackInLowerCase(t *testing.T) {
	r, err := id.ParseRowID("01890A5D-AC96-774B-BCCE-B302099A8057")
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if r.Text() != aUUIDv7 {
		t.Fatalf("got %q; want %q", r.Text(), aUUIDv7)
	}
}

func TestARowIDParsedFromTextHasTheUUIDsBytes(t *testing.T) {
	r, err := id.ParseRowID(aUUIDv7)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if got := r.Bytes(); string(got[:]) != string(aUUIDv7Bytes()) {
		t.Fatalf("got %x; want %x", got, aUUIDv7Bytes())
	}
}

func TestARowIDIsBuiltFromTheBytesOfAUUIDv7(t *testing.T) {
	r, err := id.RowIDFromBytes(aUUIDv7Bytes())
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if r.Text() != aUUIDv7 {
		t.Fatalf("got %q; want %q", r.Text(), aUUIDv7)
	}
}

func TestARowIDDoesNotKeepTheSliceItWasBuiltFrom(t *testing.T) {
	b := aUUIDv7Bytes()
	r, err := id.RowIDFromBytes(b)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	b[15] = 0
	if r.Text() != aUUIDv7 {
		t.Fatalf("got %q; want %q", r.Text(), aUUIDv7)
	}
}

func TestTwoRowIDsOfTheSameUUIDAreEqual(t *testing.T) {
	a, errA := id.ParseRowID(aUUIDv7)
	b, errB := id.RowIDFromBytes(aUUIDv7Bytes())
	if errA != nil || errB != nil {
		t.Fatalf("got %v, %v; want nil, nil", errA, errB)
	}
	if a != b {
		t.Fatal("got unequal row IDs; want equal")
	}
}

func TestARowIDRejectsTextThatIsNotAWellFormedUUIDv7(t *testing.T) {
	cases := notWellFormedText()
	cases["a UUIDv4"] = aUUIDv4
	for name, text := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := id.ParseRowID(text); !errors.Is(err, id.ErrNotRowID) {
				t.Fatalf("got %v; want %v", err, id.ErrNotRowID)
			}
		})
	}
}

func TestARowIDRejectsBytesThatAreNotAWellFormedUUIDv7(t *testing.T) {
	wrongVariant := aUUIDv7Bytes()
	wrongVariant[8] = 0xcc
	cases := map[string][]byte{
		"none":              nil,
		"15 bytes":          aUUIDv7Bytes()[:15],
		"17 bytes":          append(aUUIDv7Bytes(), 0),
		"16 zero bytes":     make([]byte, 16),
		"a UUIDv4":          aUUIDv4Bytes(),
		"the wrong variant": wrongVariant,
	}
	for name, b := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := id.RowIDFromBytes(b); !errors.Is(err, id.ErrNotRowID) {
				t.Fatalf("got %v; want %v", err, id.ErrNotRowID)
			}
		})
	}
}

func TestARowIDsFailureHoldsNothingOfTheTextItWasGiven(t *testing.T) {
	_, err := id.ParseRowID(aUUIDv4)
	if err == nil {
		t.Fatal("got nil; want an error")
	}
	if containsAny(err.Error(), aUUIDv4, "f47ac10b") {
		t.Fatalf("got %q; want nothing of %q", err, aUUIDv4)
	}
}

func TestAPublicIDParsesAUUIDv4AndGivesItsTextBack(t *testing.T) {
	p, err := id.ParsePublicID(aUUIDv4)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if p.String() != aUUIDv4 {
		t.Fatalf("got %q; want %q", p.String(), aUUIDv4)
	}
}

func TestAPublicIDParsedFromTextHasTheUUIDsBytes(t *testing.T) {
	p, err := id.ParsePublicID(aUUIDv4)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if got := p.Bytes(); string(got[:]) != string(aUUIDv4Bytes()) {
		t.Fatalf("got %x; want %x", got, aUUIDv4Bytes())
	}
}

func TestAPublicIDIsBuiltFromTheBytesOfAUUIDv4(t *testing.T) {
	p, err := id.PublicIDFromBytes(aUUIDv4Bytes())
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if p.String() != aUUIDv4 {
		t.Fatalf("got %q; want %q", p.String(), aUUIDv4)
	}
}

func TestAPublicIDRejectsTextThatIsNotAWellFormedUUIDv4(t *testing.T) {
	cases := notWellFormedText()
	cases["a UUIDv7"] = aUUIDv7
	for name, text := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := id.ParsePublicID(text); !errors.Is(err, id.ErrNotPublicID) {
				t.Fatalf("got %v; want %v", err, id.ErrNotPublicID)
			}
		})
	}
}

func TestAPublicIDRejectsBytesThatAreNotAWellFormedUUIDv4(t *testing.T) {
	wrongVariant := aUUIDv4Bytes()
	wrongVariant[8] = 0x05
	cases := map[string][]byte{
		"none":              nil,
		"15 bytes":          aUUIDv4Bytes()[:15],
		"17 bytes":          append(aUUIDv4Bytes(), 0),
		"16 zero bytes":     make([]byte, 16),
		"a UUIDv7":          aUUIDv7Bytes(),
		"the wrong variant": wrongVariant,
	}
	for name, b := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := id.PublicIDFromBytes(b); !errors.Is(err, id.ErrNotPublicID) {
				t.Fatalf("got %v; want %v", err, id.ErrNotPublicID)
			}
		})
	}
}
