// Package ids makes new IDs: a UUIDv7 for each new row and a UUIDv4 for each
// new public ID. It is the only maker of new IDs; values/id only checks and
// holds the IDs it is given.
package ids

import (
	"fmt"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/adapters/driven/clock"
	"github.com/ianflanagan1/go-ianf-strict-shared-code/adapters/driven/random"
	"github.com/ianflanagan1/go-ianf-strict-shared-code/values/id"
)

const (
	rowIDVersion    = 0x70
	publicIDVersion = 0x40
	rfc9562Variant  = 0x80
	// rowIDRandomBytes holds a UUIDv7's 74 random bits: 4 beside the version,
	// 6 beside the variant and 64 more.
	rowIDRandomBytes = 10
)

// System makes new IDs from the system clock and the system random source.
type System struct {
	clock  clock.System
	random random.Secure
}

// NewSystem returns the system IDs, which reads the time from clock and the
// random bits from random.
func NewSystem(clock clock.System, random random.Secure) System {
	return System{clock: clock, random: random}
}

// NewRowID returns a new row ID: a UUIDv7 of the current millisecond.
func (s System) NewRowID() id.RowID {
	var r [rowIDRandomBytes]byte
	s.read(r[:])
	b := layOutRowID(s.clock.Now().UnixMilli(), r)
	row, err := id.RowIDFromBytes(b[:])
	if err != nil {
		panic(fmt.Sprintf("ids: laid out a row ID that isn't a UUIDv7: %v", err))
	}
	return row
}

// NewPublicID returns a new public ID: a UUIDv4.
func (s System) NewPublicID() id.PublicID {
	var r [16]byte
	s.read(r[:])
	b := layOutPublicID(r)
	public, err := id.PublicIDFromBytes(b[:])
	if err != nil {
		panic(fmt.Sprintf("ids: laid out a public ID that isn't a UUIDv4: %v", err))
	}
	return public
}

// read fills b from the random source, which crypto/rand documents never to
// fail.
func (s System) read(b []byte) {
	if _, err := s.random.Read(b); err != nil {
		panic(fmt.Sprintf("ids: the system random source failed: %v", err))
	}
}

// layOutRowID lays out a UUIDv7: the low 48 bits of ms, the version, the 74
// random bits r holds, and the variant.
func layOutRowID(ms int64, r [rowIDRandomBytes]byte) [16]byte {
	var b [16]byte
	for i := range 6 {
		b[i] = byte(ms >> (8 * (5 - i)) & 0xff)
	}
	b[6] = rowIDVersion | r[0]&0x0f
	b[7] = r[1]
	b[8] = rfc9562Variant | r[2]&0x3f
	copy(b[9:], r[3:])
	return b
}

// layOutPublicID lays out a UUIDv4: the version, 122 of r's random bits, and
// the variant.
func layOutPublicID(r [16]byte) [16]byte {
	b := r
	b[6] = publicIDVersion | b[6]&0x0f
	b[8] = rfc9562Variant | b[8]&0x3f
	return b
}
