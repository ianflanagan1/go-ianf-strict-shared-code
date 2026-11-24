package email_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/core/equalitytest"
	"github.com/ianflanagan1/go-ianf-strict-shared-code/values/email"
)

// held is text and what one part of the address it holds must be.
type held struct {
	text string
	want string
}

// parsed returns the address text holds, and fails the test when it holds none.
func parsed(t *testing.T, text string) email.Address {
	t.Helper()
	address, err := email.ParseAddress(text)
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	return address
}

// assertLocalParts fails unless each case's text holds an address with the local part wanted.
func assertLocalParts(t *testing.T, cases map[string]held) {
	t.Helper()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := parsed(t, c.text).LocalPart; got != c.want {
				t.Fatalf("got %q; want %q", got, c.want)
			}
		})
	}
}

// assertDomains fails unless each case's text holds an address with the domain wanted.
func assertDomains(t *testing.T, cases map[string]held) {
	t.Helper()
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := parsed(t, c.text).Domain; got != c.want {
				t.Fatalf("got %q; want %q", got, c.want)
			}
		})
	}
}

// assertRejected fails unless each text holds no address, for the reason want.
func assertRejected(t *testing.T, want email.AddressError, texts map[string]string) {
	t.Helper()
	for name, text := range texts {
		t.Run(name, func(t *testing.T) {
			if _, err := email.ParseAddress(text); !errors.Is(err, want) {
				t.Fatalf("got %v; want %v", err, want)
			}
		})
	}
}

func TestAnAddressIsSplitAtItsAtSignIntoALocalPartAndADomain(t *testing.T) {
	address := parsed(t, "ian@example.com")

	if got := address.LocalPart; got != "ian" {
		t.Errorf("got local part %q; want %q", got, "ian")
	}
	if got := address.Domain; got != "example.com" {
		t.Errorf("got domain %q; want %q", got, "example.com")
	}
}

func TestTheTextOfAnAddressIsItsLocalPartAnAtSignAndItsDomain(t *testing.T) {
	cases := map[string]held{
		"an address typed in its normal form": {text: "ian@example.com", want: "ian@example.com"},
		"an address typed in another case":    {text: "Ian@EXAMPLE.com", want: "Ian@example.com"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := parsed(t, c.text).Text(); got != c.want {
				t.Fatalf("got %q; want %q", got, c.want)
			}
		})
	}
}

func TestTheLocalPartKeepsTheCaseItWasTypedIn(t *testing.T) {
	assertLocalParts(t, map[string]held{
		"capital letters": {text: "Ian.Flanagan@example.com", want: "Ian.Flanagan"},
	})
}

func TestTheLocalPartMayHoldDotsBetweenItsAtomsAndEveryPrintableASCIICharacterTheStandardAllowsInOne(t *testing.T) {
	assertLocalParts(t, map[string]held{
		"letters, dots and a plus sign":  {text: "ian.flanagan+news@example.com", want: "ian.flanagan+news"},
		"every punctuation mark allowed": {text: "!#$%&'*+/=?^_`{|}~-@example.com", want: "!#$%&'*+/=?^_`{|}~-"},
		"every digit":                    {text: "0123456789@example.com", want: "0123456789"},
	})
}

func TestTheLocalPartMayHoldLettersOfAnyScript(t *testing.T) {
	assertLocalParts(t, map[string]held{
		"Chinese":             {text: "用户@example.com", want: "用户"},
		"Greek":               {text: "δοκιμή@example.com", want: "δοκιμή"},
		"Latin, with accents": {text: "josé.núñez@example.com", want: "josé.núñez"},
	})
}

func TestTheLocalPartIsNormalisedToNFC(t *testing.T) {
	assertLocalParts(t, map[string]held{
		"a letter followed by a combining accent": {text: "jose\u0301@example.com", want: "jos\u00e9"},
		"a letter typed with its accent":          {text: "jos\u00e9@example.com", want: "jos\u00e9"},
	})
}

func TestTheDomainIsNormalisedToLowerCase(t *testing.T) {
	assertDomains(t, map[string]held{
		"capital letters": {text: "ian@Example.COM", want: "example.com"},
	})
}

func TestADomainInAnotherScriptIsNormalisedToItsASCIIForm(t *testing.T) {
	assertDomains(t, map[string]held{
		"a letter beyond ASCII":                   {text: "ian@münchen.de", want: "xn--mnchen-3ya.de"},
		"a capital letter beyond ASCII":           {text: "ian@MÜNCHEN.de", want: "xn--mnchen-3ya.de"},
		"a letter followed by a combining accent": {text: "ian@mu\u0308nchen.de", want: "xn--mnchen-3ya.de"},
		"every label in another script":           {text: "用户@例子.中国", want: "xn--fsqu00a.xn--fiqs8s"},
		"a sharp s, which is kept":                {text: "ian@straße.de", want: "xn--strae-oqa.de"},
	})
}

func TestADomainAlreadyInItsASCIIFormIsKept(t *testing.T) {
	assertDomains(t, map[string]held{
		"an xn-- label":                  {text: "ian@xn--mnchen-3ya.de", want: "xn--mnchen-3ya.de"},
		"an xn-- label in capitals":      {text: "ian@XN--MNCHEN-3YA.de", want: "xn--mnchen-3ya.de"},
		"a hyphen inside a label":        {text: "ian@a-b.example.co.uk", want: "a-b.example.co.uk"},
		"a first label that is a number": {text: "ian@1.example.com", want: "1.example.com"},
	})
}

func TestTwoSpellingsOfOneAddressParseToEqualValues(t *testing.T) {
	cases := map[string]struct{ one, other string }{
		"a domain in two cases":                  {one: "ian@example.com", other: "ian@EXAMPLE.com"},
		"a domain in its ASCII form and its own": {one: "ian@xn--mnchen-3ya.de", other: "ian@münchen.de"},
		"a local part composed and not composed": {one: "jos\u00e9@example.com", other: "jose\u0301@example.com"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			equalitytest.AssertEqual(t, parsed(t, c.one), parsed(t, c.other))
		})
	}
}

func TestTheLongestLocalPartAndTheLongestAddressAreAccepted(t *testing.T) {
	longestLocalPart := strings.Repeat("a", 64)
	domainOf189Bytes := strings.Repeat(strings.Repeat("b", 62)+".", 2) + strings.Repeat("c", 59) + ".com"

	if got := parsed(t, longestLocalPart+"@example.com").LocalPart; got != longestLocalPart {
		t.Errorf("got local part %q; want %q", got, longestLocalPart)
	}
	if got := len(parsed(t, longestLocalPart+"@"+domainOf189Bytes).Text()); got != 254 {
		t.Errorf("got an address of %d bytes; want 254", got)
	}
}

func TestTextWithNoAtSignIsRejected(t *testing.T) {
	assertRejected(t, email.ErrMissingAt, map[string]string{
		"no text":        "",
		"a name alone":   "ian",
		"a domain alone": "example.com",
	})
}

func TestTextThatIsNotUTF8IsRejected(t *testing.T) {
	assertRejected(t, email.ErrInvalidEncoding, map[string]string{
		"a byte no UTF-8 holds, in the local part": "ia\xffn@example.com",
		"a sequence cut short, in the domain":      "ian@exa\xc3mple.com",
		"a byte no UTF-8 holds, alone":             "\xff",
	})
}

func TestAnEmptyLocalPartIsRejected(t *testing.T) {
	assertRejected(t, email.ErrLocalPartEmpty, map[string]string{
		"nothing before the at sign": "@example.com",
	})
}

func TestALocalPartOver64BytesIsRejected(t *testing.T) {
	assertRejected(t, email.ErrLocalPartTooLong, map[string]string{
		"65 letters of one byte each":    strings.Repeat("a", 65) + "@example.com",
		"22 letters of three bytes each": strings.Repeat("用", 22) + "@example.com",
	})
}

func TestALocalPartWithADotAtEitherEndOrTwoDotsTogetherIsRejected(t *testing.T) {
	assertRejected(t, email.ErrLocalPartInvalid, map[string]string{
		"a dot at the start": ".ian@example.com",
		"a dot at the end":   "ian.@example.com",
		"two dots together":  "ian..flanagan@example.com",
		"a dot alone":        ".@example.com",
	})
}

func TestALocalPartWithAnASCIICharacterOutsideTheAllowedSetIsRejected(t *testing.T) {
	assertRejected(t, email.ErrLocalPartInvalid, map[string]string{
		"a space":          "ian flanagan@example.com",
		"a second at sign": "ian@home@example.com",
		"quotation marks":  `"ian"@example.com`,
		"parentheses":      "ian(work)@example.com",
		"a comma":          "ian,flanagan@example.com",
		"angle brackets":   "<ian>@example.com",
		"a backslash":      `ian\flanagan@example.com`,
		"a colon":          "ian:flanagan@example.com",
		"a semicolon":      "ian;flanagan@example.com",
		"square brackets":  "ian[1]@example.com",
	})
}

func TestALocalPartWithASpaceOrAControlCharacterAtEitherEndIsRejected(t *testing.T) {
	assertRejected(t, email.ErrLocalPartInvalid, map[string]string{
		"a space at the start":                     " ian@example.com",
		"a space at the end":                       "ian @example.com",
		"a newline at the end":                     "ian\n@example.com",
		"a newline at the start":                   "\nian@example.com",
		"a carriage return and newline at the end": "ian\r\n@example.com",
		"a tab at the start":                       "\tian@example.com",
		"a NUL at the start":                       "\x00ian@example.com",
		"a delete character at the end":            "ian\x7f@example.com",
	})
}

func TestALocalPartWithAControlASpaceOrASeparatorBeyondASCIIIsRejected(t *testing.T) {
	assertRejected(t, email.ErrLocalPartInvalid, map[string]string{
		"U+0080, the first C1 control":     "ia\u0080n@example.com",
		"U+0085 next line":                 "ia\u0085n@example.com",
		"U+009F, the last C1 control":      "ia\u009fn@example.com",
		"U+00A0 no-break space":            "ia\u00a0n@example.com",
		"U+1680 Ogham space mark":          "ia\u1680n@example.com",
		"U+2000 en quad":                   "ia\u2000n@example.com",
		"U+2028 line separator":            "ia\u2028n@example.com",
		"U+2029 paragraph separator":       "ia\u2029n@example.com",
		"U+202F narrow no-break space":     "ia\u202fn@example.com",
		"U+205F medium mathematical space": "ia\u205fn@example.com",
		"U+3000 ideographic space":         "ia\u3000n@example.com",
	})
}

func TestALocalPartWithACharacterThatShowsNothingIsRejected(t *testing.T) {
	assertRejected(t, email.ErrLocalPartInvalid, map[string]string{
		"U+00AD soft hyphen":                                 "ia\u00adn@example.com",
		"U+034F combining grapheme joiner":                   "ia\u034fn@example.com",
		"U+115F Hangul choseong filler":                      "ia\u115fn@example.com",
		"U+1160 Hangul jungseong filler":                     "ia\u1160n@example.com",
		"U+17B4 Khmer vowel inherent aq":                     "ia\u17b4n@example.com",
		"U+17B5 Khmer vowel inherent aa":                     "ia\u17b5n@example.com",
		"U+180B Mongolian free variation selector one":       "ia\u180bn@example.com",
		"U+180F Mongolian free variation selector four":      "ia\u180fn@example.com",
		"U+200B zero width space":                            "ia\u200bn@example.com",
		"U+200D zero width joiner":                           "ia\u200dn@example.com",
		"U+2060 word joiner":                                 "ia\u2060n@example.com",
		"U+206F nominal digit shapes":                        "ia\u206fn@example.com",
		"U+2800 Braille pattern blank":                       "ia\u2800n@example.com",
		"U+3164 Hangul filler":                               "ia\u3164n@example.com",
		"U+FE00 variation selector-1":                        "ia\ufe00n@example.com",
		"U+FE0F variation selector-16":                       "ia\ufe0fn@example.com",
		"U+FEFF byte order mark":                             "ia\ufeffn@example.com",
		"U+FFA0 halfwidth Hangul filler":                     "ia\uffa0n@example.com",
		"U+1BCA0 shorthand format letter overlap":            "ia\U0001bca0n@example.com",
		"U+1BCA3 shorthand format up step":                   "ia\U0001bca3n@example.com",
		"U+1D173 musical symbol begin beam":                  "ia\U0001d173n@example.com",
		"U+1D17A musical symbol end phrase":                  "ia\U0001d17an@example.com",
		"U+E0001 language tag":                               "ia\U000e0001n@example.com",
		"U+E0FFF, the last of the tags and their supplement": "ia\U000e0fffn@example.com",
	})
}

func TestALocalPartWithACharacterThatChangesTheDirectionOfTextIsRejected(t *testing.T) {
	assertRejected(t, email.ErrLocalPartInvalid, map[string]string{
		"U+061C Arabic letter mark":      "ia\u061cn@example.com",
		"U+200E left-to-right mark":      "ia\u200en@example.com",
		"U+200F right-to-left mark":      "ia\u200fn@example.com",
		"U+202A left-to-right embedding": "ia\u202an@example.com",
		"U+202E right-to-left override":  "ia\u202en@example.com",
		"U+2066 left-to-right isolate":   "ia\u2066n@example.com",
		"U+2069 pop directional isolate": "ia\u2069n@example.com",
	})
}

func TestALocalPartWithAPrivateUseOrAPlaceholderCharacterIsRejected(t *testing.T) {
	assertRejected(t, email.ErrLocalPartInvalid, map[string]string{
		"U+E000, the first of the private use area":    "ia\ue000n@example.com",
		"U+F8FF, the last of the private use area":     "ia\uf8ffn@example.com",
		"U+FFF0, the first of the specials":            "ia\ufff0n@example.com",
		"U+FFFC object replacement character":          "ia\ufffcn@example.com",
		"U+FFFD replacement character":                 "ia\ufffdn@example.com",
		"U+FFFF, the last of the specials":             "ia\uffffn@example.com",
		"U+F0000, the first of the private use planes": "ia\U000f0000n@example.com",
		"U+10FFFF, the last of the private use planes": "ia\U0010ffffn@example.com",
	})
}

func TestTheCharactersOnEitherSideOfEachRejectedRangeAreAccepted(t *testing.T) {
	assertLocalParts(t, map[string]held{
		"U+00A1, after the no-break space":             {text: "ia\u00a1n@example.com", want: "ia\u00a1n"},
		"U+00AC, before the soft hyphen":               {text: "ia\u00acn@example.com", want: "ia\u00acn"},
		"U+00AE, after the soft hyphen":                {text: "ia\u00aen@example.com", want: "ia\u00aen"},
		"U+1FFE, before the spaces":                    {text: "ia\u1ffen@example.com", want: "ia\u1ffen"},
		"U+2010, after the direction marks":            {text: "ia\u2010n@example.com", want: "ia\u2010n"},
		"U+2027, before the line separator":            {text: "ia\u2027n@example.com", want: "ia\u2027n"},
		"U+2030, after the narrow no-break space":      {text: "ia\u2030n@example.com", want: "ia\u2030n"},
		"U+205E, before the mathematical space":        {text: "ia\u205en@example.com", want: "ia\u205en"},
		"U+2070, after the direction isolates":         {text: "ia\u2070n@example.com", want: "ia\u2070n"},
		"U+3001, after the ideographic space":          {text: "ia\u3001n@example.com", want: "ia\u3001n"},
		"U+D7FF, before the surrogates":                {text: "ia\ud7ffn@example.com", want: "ia\ud7ffn"},
		"U+FDFF, before the variation selectors":       {text: "ia\ufdffn@example.com", want: "ia\ufdffn"},
		"U+FE10, after the variation selectors":        {text: "ia\ufe10n@example.com", want: "ia\ufe10n"},
		"U+FFEF, before the specials":                  {text: "ia\uffefn@example.com", want: "ia\uffefn"},
		"U+10000, after the specials":                  {text: "ia\U00010000n@example.com", want: "ia\U00010000n"},
		"U+1F600, an emoji":                            {text: "ia\U0001f600n@example.com", want: "ia\U0001f600n"},
		"U+E1000, after the tags and their supplement": {text: "ia\U000e1000n@example.com", want: "ia\U000e1000n"},
		"U+EFFFD, before the private use planes":       {text: "ia\U000efffdn@example.com", want: "ia\U000efffdn"},
	})
}

func TestAnEmptyDomainIsRejected(t *testing.T) {
	assertRejected(t, email.ErrDomainEmpty, map[string]string{
		"nothing after the at sign": "ian@",
	})
}

func TestADomainOfOneLabelIsRejected(t *testing.T) {
	assertRejected(t, email.ErrDomainNotFullyQualified, map[string]string{
		"a host name of one label": "ian@localhost",
		"a top-level domain alone": "ian@com",
	})
}

func TestADomainWithAnEmptyLabelOrADotAtEitherEndIsRejected(t *testing.T) {
	assertRejected(t, email.ErrDomainInvalid, map[string]string{
		"two dots together":                   "ian@example..com",
		"a dot at the start":                  "ian@.example.com",
		"a dot at the end":                    "ian@example.com.",
		"three dots at the end":               "ian@example.com...",
		"an ideographic full stop at the end": "ian@example.com。",
		"one label and a dot at the end":      "ian@localhost.",
		"a dot alone":                         "ian@.",
	})
}

func TestADomainWithACharacterNoHostNameHoldsIsRejected(t *testing.T) {
	assertRejected(t, email.ErrDomainInvalid, map[string]string{
		"an underscore":               "ian@ex_ample.com",
		"a space inside a label":      "ian@exam ple.com",
		"a newline inside":            "ian@exa\n.com",
		"a newline at the end":        "ian@example.com\n",
		"a space at the end":          "ian@example.com ",
		"a NUL at the end":            "ian@example.com\x00",
		"a path":                      "ian@example.com/path",
		"a port":                      "ian@example.com:25",
		"an IPv4 address in brackets": "ian@[192.168.0.1]",
		"an IPv6 address in brackets": "ian@[IPv6:2001:db8::1]",
	})
}

func TestADomainWithALabelThatStartsOrEndsWithAHyphenIsRejected(t *testing.T) {
	assertRejected(t, email.ErrDomainInvalid, map[string]string{
		"a hyphen at the start of the first label": "ian@-example.com",
		"a hyphen at the end of the first label":   "ian@example-.com",
		"a hyphen at the start of the last label":  "ian@example.-com",
	})
}

func TestADomainWhoseLastLabelIsAllDigitsIsRejected(t *testing.T) {
	assertRejected(t, email.ErrDomainInvalid, map[string]string{
		"an IPv4 address":               "ian@127.0.0.1",
		"a last label that is a number": "ian@example.123",
	})
}

func TestADomainWithALabelOver63BytesOrOver253BytesInAllIsRejected(t *testing.T) {
	domainOf254Bytes := strings.Repeat(strings.Repeat("b", 62)+".", 3) + strings.Repeat("c", 61) + ".com"

	assertRejected(t, email.ErrDomainInvalid, map[string]string{
		"a label of 64 bytes":   "ian@" + strings.Repeat("b", 64) + ".com",
		"a domain of 254 bytes": "i@" + domainOf254Bytes,
	})
}

func TestADomainWhoseASCIIFormDoesNotDecodeIsRejected(t *testing.T) {
	assertRejected(t, email.ErrDomainInvalid, map[string]string{
		"an xn-- label that is no Punycode": "ian@xn--zzzzzz.de",
	})
}

func TestADomainThatMixesDirectionsOfTextInOneLabelOrHoldsAJoinerOutOfPlaceIsRejected(t *testing.T) {
	assertRejected(t, email.ErrDomainInvalid, map[string]string{
		"a Hebrew letter between a Latin letter and a digit": "ian@a\u05d01.com",
		"a zero width joiner at the start of a label":        "ian@\u200da.com",
	})
}

func TestAnAddressOver254BytesIsRejected(t *testing.T) {
	domainOf190Bytes := strings.Repeat(strings.Repeat("b", 62)+".", 2) + strings.Repeat("c", 60) + ".com"

	assertRejected(t, email.ErrTooLong, map[string]string{
		"a local part of 64 bytes and a domain of 190": strings.Repeat("a", 64) + "@" + domainOf190Bytes,
	})
}

func TestEachFailureSaysWhyTheTextHoldsNoAddress(t *testing.T) {
	cases := map[string]struct {
		failure email.AddressError
		want    string
	}{
		"invalid encoding":           {failure: email.ErrInvalidEncoding, want: "not an email address: the text isn't UTF-8"},
		"missing at sign":            {failure: email.ErrMissingAt, want: "not an email address: the text has no @"},
		"local part empty":           {failure: email.ErrLocalPartEmpty, want: "not an email address: nothing comes before the last @"},
		"local part too long":        {failure: email.ErrLocalPartTooLong, want: "not an email address: the local part is over 64 bytes"},
		"local part invalid":         {failure: email.ErrLocalPartInvalid, want: "not an email address: the local part holds a character an address may not, or a dot out of place"},
		"domain empty":               {failure: email.ErrDomainEmpty, want: "not an email address: nothing comes after the last @"},
		"domain invalid":             {failure: email.ErrDomainInvalid, want: "not an email address: the domain is no host name, or its last label is all digits"},
		"domain not fully qualified": {failure: email.ErrDomainNotFullyQualified, want: "not an email address: the domain is one label"},
		"too long":                   {failure: email.ErrTooLong, want: "not an email address: the address is over 254 bytes"},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := c.failure.Error(); got != c.want {
				t.Fatalf("got %q; want %q", got, c.want)
			}
		})
	}
}
