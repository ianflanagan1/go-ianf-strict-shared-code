// Package id is the shared ID unit-of-code: the row ID every owned table's ID
// type wraps, the public ID every public ID type wraps, and the Internal
// wrapper for values that must stay inside the program.
//
// Its factory functions only check and hold a UUID they are given; they never
// make a new one. New IDs come from the system IDs adapter, which passes the
// bytes of each one it lays out to RowIDFromBytes or PublicIDFromBytes.
package id

import (
	"errors"
	"fmt"
	"uuid"
)

// ErrNotRowID is the failure of a row ID's factory functions when given
// anything that isn't a well-formed UUIDv7.
var ErrNotRowID = errors.New("not a row ID: want a well-formed UUIDv7")

// ErrNotPublicID is the failure of a public ID's factory functions when given
// anything that isn't a well-formed UUIDv4.
var ErrNotPublicID = errors.New("not a public ID: want a well-formed UUIDv4")

const (
	rowIDVersion    = 7
	publicIDVersion = 4
	// textLength is the length of a UUID's hyphenated text, the only text form
	// the factory functions accept.
	textLength = 36
	byteLength = 16
	// rfc9562Variant is the top two bits of byte 8 in every UUID RFC 9562 lays out.
	rfc9562Variant = 0b10
)

// RowID is the ID of a row in an owned table: a UUIDv7. Each table's ID type
// wraps one. It is an internal value: only decisions and internal adapters
// read its Text or Bytes.
type RowID struct {
	// Bytes is the row ID's 16 bytes.
	Bytes [16]byte
}

// ParseRowID returns the row ID that text, a UUIDv7 in hyphenated form in
// either case, holds, or a failure wrapping ErrNotRowID.
func ParseRowID(text string) (RowID, error) {
	u, err := parse(text, rowIDVersion)
	if err != nil {
		return RowID{}, fmt.Errorf("%w: %w", ErrNotRowID, err)
	}
	return RowID{Bytes: u}, nil
}

// RowIDFromBytes returns the row ID that b, the 16 bytes of a UUIDv7, holds,
// or a failure wrapping ErrNotRowID.
func RowIDFromBytes(b []byte) (RowID, error) {
	u, err := fromBytes(b, rowIDVersion)
	if err != nil {
		return RowID{}, fmt.Errorf("%w: %w", ErrNotRowID, err)
	}
	return RowID{Bytes: u}, nil
}

// Text returns the row ID in lower-case hyphenated form.
func (r RowID) Text() string { return uuid.UUID(r.Bytes).String() }

// PublicID is the ID that things outside the program refer to a row by: a
// UUIDv4. Each table's public ID type wraps one.
type PublicID struct {
	// Bytes is the public ID's 16 bytes.
	Bytes [16]byte
}

// ParsePublicID returns the public ID that text, a UUIDv4 in hyphenated form
// in either case, holds, or a failure wrapping ErrNotPublicID.
func ParsePublicID(text string) (PublicID, error) {
	u, err := parse(text, publicIDVersion)
	if err != nil {
		return PublicID{}, fmt.Errorf("%w: %w", ErrNotPublicID, err)
	}
	return PublicID{Bytes: u}, nil
}

// PublicIDFromBytes returns the public ID that b, the 16 bytes of a UUIDv4,
// holds, or a failure wrapping ErrNotPublicID.
func PublicIDFromBytes(b []byte) (PublicID, error) {
	u, err := fromBytes(b, publicIDVersion)
	if err != nil {
		return PublicID{}, fmt.Errorf("%w: %w", ErrNotPublicID, err)
	}
	return PublicID{Bytes: u}, nil
}

// String returns the public ID in lower-case hyphenated form.
func (p PublicID) String() string { return uuid.UUID(p.Bytes).String() }

// parse reads text as a UUID of version. Its failures hold nothing of text,
// which may be an internal value.
func parse(text string, version byte) ([16]byte, error) {
	if len(text) != textLength {
		return [16]byte{}, fmt.Errorf("text is %d bytes, not the %d of a hyphenated UUID", len(text), textLength)
	}
	u, err := uuid.Parse(text)
	if err != nil {
		return [16]byte{}, errors.New("text isn't a hyphenated UUID")
	}
	return u, checkLayout(u, version)
}

// fromBytes reads b as a UUID of version.
func fromBytes(b []byte, version byte) ([16]byte, error) {
	if len(b) != byteLength {
		return [16]byte{}, fmt.Errorf("got %d bytes, not %d", len(b), byteLength)
	}
	u := [16]byte(b)
	return u, checkLayout(u, version)
}

// checkLayout reports whether u has version and RFC 9562's variant.
func checkLayout(u [16]byte, version byte) error {
	if got := u[6] >> 4; got != version {
		return fmt.Errorf("version is %d", got)
	}
	if u[8]>>6 != rfc9562Variant {
		return errors.New("variant isn't RFC 9562's")
	}
	return nil
}
