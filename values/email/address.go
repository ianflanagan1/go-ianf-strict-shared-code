// Package email is the shared email address unit-of-code: the address a person
// types into a form, and why text holds none.
//
// It declares no comparison, and nothing that reads a +tag. Whether two
// addresses that differ only in the case of their local parts are one mailbox
// is a question for a project's decision code, in a function with its own
// name.
package email

import (
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/net/idna"
	"golang.org/x/text/unicode/norm"
)

const (
	localPartMaxBytes = 64
	addressMaxBytes   = 254
	// domainMaxCodePoints is the longest domain that is converted to its ASCII
	// form. The ASCII form has a byte or more for each code point, and a host
	// name is at most 253 bytes, so nothing longer is one, whatever it holds.
	domainMaxCodePoints = 255
	// atomPunctuation is the punctuation an atom may hold. With the ASCII
	// letters and digits it is RFC 5322's atext.
	atomPunctuation = "!#$%&'*+/=?^_`{|}~-"
)

// AddressError is why text isn't an email address: the failure ParseAddress
// returns in place of one. errors.As finds it in an error, and a switch over
// it names every case.
type AddressError int

const (
	// ErrInvalidEncoding is the failure of text that isn't UTF-8.
	ErrInvalidEncoding AddressError = iota + 1
	// ErrMissingAt is the failure of text that has no @.
	ErrMissingAt
	// ErrLocalPartEmpty is the failure of text with nothing before its last @.
	ErrLocalPartEmpty
	// ErrLocalPartTooLong is the failure of a local part that, normalised, is
	// over 64 bytes.
	ErrLocalPartTooLong
	// ErrLocalPartInvalid is the failure of a local part that holds a
	// character an address may not, a second @ among them, or has a dot at
	// either end or two together.
	ErrLocalPartInvalid
	// ErrDomainEmpty is the failure of text with nothing after its last @.
	ErrDomainEmpty
	// ErrDomainInvalid is the failure of a domain that is no host name, or
	// whose last label is all digits, as an IP address's is.
	ErrDomainInvalid
	// ErrDomainNotFullyQualified is the failure of a domain of one label, such
	// as localhost.
	ErrDomainNotFullyQualified
	// ErrTooLong is the failure of an address that, normalised, is over 254
	// bytes.
	ErrTooLong
)

// Error returns why the text holds no address. It holds nothing of the text.
func (e AddressError) Error() string {
	return "not an email address: " + e.reason()
}

func (e AddressError) reason() string {
	switch e {
	case ErrInvalidEncoding:
		return "the text isn't UTF-8"
	case ErrMissingAt:
		return "the text has no @"
	case ErrLocalPartEmpty:
		return "nothing comes before the last @"
	case ErrLocalPartTooLong:
		return "the local part is over 64 bytes"
	case ErrLocalPartInvalid:
		return "the local part holds a character an address may not, or a dot out of place"
	case ErrDomainEmpty:
		return "nothing comes after the last @"
	case ErrDomainInvalid:
		return "the domain is no host name, or its last label is all digits"
	case ErrDomainNotFullyQualified:
		return "the domain is one label"
	case ErrTooLong:
		return "the address is over 254 bytes"
	}
	return "AddressError(" + strconv.Itoa(int(e)) + ")"
}

// Address is an email address of the kind a person types into a form: a local
// part of atoms joined by dots, in any script, and a domain that is a host
// name of two labels or more. It has no quoted local part, no comment and no
// IP address in place of a domain, which the standards allow and people don't
// use.
//
// Two spellings of one address parse to equal values, since
// ParseAddress normalises both parts: the local part to NFC, and the domain to
// its lower-case ASCII form. The local part keeps its case, which the
// standards leave to the mail server that holds the mailbox.
type Address struct {
	// LocalPart is what comes before the @, in NFC and in the case it was
	// typed in: "Ian.Flanagan+news", "用户".
	LocalPart string
	// Domain is what comes after the @, in lower-case ASCII, with a label in
	// another script in its xn-- form: "example.com", "xn--mnchen-3ya.de".
	Domain string
}

// ParseAddress returns the address that text holds, or an AddressError that
// says why it holds none. The text is the address alone, with nothing around
// it: a caller that reads a form field trims it first.
func ParseAddress(text string) (Address, error) {
	if !utf8.ValidString(text) {
		return Address{}, ErrInvalidEncoding
	}
	before, after, ok := strings.CutLast(text, "@")
	if !ok {
		return Address{}, ErrMissingAt
	}
	localPart, err := parseLocalPart(before)
	if err != nil {
		return Address{}, err
	}
	domain, err := parseDomain(after)
	if err != nil {
		return Address{}, err
	}
	if len(localPart)+len("@")+len(domain) > addressMaxBytes {
		return Address{}, ErrTooLong
	}
	return Address{LocalPart: localPart, Domain: domain}, nil
}

// Text returns the address as it is written: the local part, an @ and the
// domain.
func (a Address) Text() string { return a.LocalPart + "@" + a.Domain }

// parseLocalPart returns typed, what the text holds before its last @, in NFC,
// or why it is no local part.
func parseLocalPart(typed string) (string, error) {
	if typed == "" {
		return "", ErrLocalPartEmpty
	}
	normalised := norm.NFC.String(typed)
	if len(normalised) > localPartMaxBytes {
		return "", ErrLocalPartTooLong
	}
	if !isAtomsJoinedByDots(normalised) {
		return "", ErrLocalPartInvalid
	}
	return normalised, nil
}

// isAtomsJoinedByDots reports whether s is atoms with a single dot between
// each and the next. No atom holds a dot, so each piece between two dots, and
// before the first and after the last, must be an atom.
func isAtomsJoinedByDots(s string) bool {
	for piece := range strings.SplitSeq(s, ".") {
		if !isAtom(piece) {
			return false
		}
	}
	return true
}

// isAtom reports whether s is one character or more, each one that an atom may
// hold.
func isAtom(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !isAtomASCII(r) && isAtomRejected(r) {
			return false
		}
	}
	return true
}

// isAtomASCII reports whether r is an ASCII character an atom may hold.
func isAtomASCII(r rune) bool {
	switch {
	case 'A' <= r && r <= 'Z', 'a' <= r && r <= 'z', '0' <= r && r <= '9':
		return true
	}
	return strings.ContainsRune(atomPunctuation, r)
}

// isAtomRejected reports whether r is in one of the ranges an atom may not
// hold.
func isAtomRejected(r rune) bool {
	for _, rejected := range atomRejected() {
		if rejected.first <= r && r <= rejected.last {
			return true
		}
	}
	return false
}

// runeRange is the characters from first to last, both among them.
type runeRange struct {
	first rune
	last  rune
}

// atomRejected lists what an atom may not hold: the rest of ASCII, and beyond
// it every character that shows nothing, shows as a space or changes how the
// text around it is shown, so that no two addresses that look the same differ.
//
// It is a fixed list of ranges and not Unicode's properties, which change with
// the version of Unicode a runtime holds, so that every language's port
// rejects the same characters.
func atomRejected() []runeRange {
	return []runeRange{
		{first: 0x0, last: 0xA0},         // ASCII, the C1 controls and the no-break space
		{first: 0xAD, last: 0xAD},        // the soft hyphen
		{first: 0x34F, last: 0x34F},      // the combining grapheme joiner
		{first: 0x61C, last: 0x61C},      // the Arabic letter mark
		{first: 0x115F, last: 0x115F},    // the Hangul choseong filler
		{first: 0x1160, last: 0x1160},    // the Hangul jungseong filler
		{first: 0x1680, last: 0x1680},    // the Ogham space
		{first: 0x17B4, last: 0x17B4},    // the Khmer inherent vowel aq
		{first: 0x17B5, last: 0x17B5},    // the Khmer inherent vowel aa
		{first: 0x180B, last: 0x180F},    // the Mongolian variation selectors and vowel separator
		{first: 0x2000, last: 0x200F},    // the spaces, zero-width characters and direction marks
		{first: 0x2028, last: 0x202F},    // the line and paragraph separators, direction embeddings and overrides and the narrow no-break space
		{first: 0x205F, last: 0x206F},    // the mathematical space, word joiner, invisible operators and direction isolates
		{first: 0x2800, last: 0x2800},    // the blank Braille pattern
		{first: 0x3000, last: 0x3000},    // the ideographic space
		{first: 0x3164, last: 0x3164},    // the Hangul filler
		{first: 0xE000, last: 0xF8FF},    // the private use area
		{first: 0xFE00, last: 0xFE0F},    // the variation selectors
		{first: 0xFEFF, last: 0xFEFF},    // the byte order mark
		{first: 0xFFA0, last: 0xFFA0},    // the halfwidth Hangul filler
		{first: 0xFFF0, last: 0xFFFF},    // the specials, the replacement character among them
		{first: 0x1BCA0, last: 0x1BCA3},  // the shorthand format controls
		{first: 0x1D173, last: 0x1D17A},  // the musical format controls
		{first: 0xE0000, last: 0xE0FFF},  // the tags and the variation selectors supplement
		{first: 0xF0000, last: 0x10FFFF}, // the private use planes
	}
}

// parseDomain returns typed, what the text holds after its last @, in
// lower-case ASCII, or why it is no domain.
func parseDomain(typed string) (string, error) {
	if typed == "" {
		return "", ErrDomainEmpty
	}
	if utf8.RuneCountInString(typed) > domainMaxCodePoints {
		return "", ErrDomainInvalid
	}
	ascii, err := hostNames().ToASCII(typed)
	if err != nil {
		return "", ErrDomainInvalid
	}
	_, after, ok := strings.CutLast(ascii, ".")
	if !ok {
		return "", ErrDomainNotFullyQualified
	}
	if isDigits(after) {
		return "", ErrDomainInvalid
	}
	return ascii, nil
}

// hostNames returns UTS #46 as host names need it: the mapping to lower case
// and NFC, only letters, digits and hyphens in a label, no hyphen at either
// end of one or in both its third and fourth places, joiners only where a
// script uses them, no mixing of directions within a label, ß and ς kept as
// they are, no empty label, and DNS's limits of 63 bytes for a label and 253
// for the name.
//
// MapForLookup comes first, since it sets several of the options that follow
// it, and a later option overrides an earlier one.
func hostNames() *idna.Profile {
	return idna.New(
		idna.MapForLookup(),
		idna.Transitional(false),
		idna.StrictDomainName(true),
		idna.CheckHyphens(true),
		idna.CheckJoiners(true),
		idna.BidiRule(),
		idna.VerifyDNSLength(true),
		idna.RemoveLeadingDots(false),
	)
}

// isDigits reports whether every character of s is an ASCII digit, as every
// character of no characters is: the last label of a domain that ends with a
// dot is empty, and no top-level domain either.
func isDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
