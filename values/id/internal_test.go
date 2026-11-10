package id_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/ianflanagan1/go-adapters/values/id"
)

const (
	placeholder = "[internal]"
	secret      = "risk notes: owes money"
)

// customer holds an Internal field unexported, as a core type does.
type customer struct {
	name  string
	notes id.Internal[string]
}

// customerView holds an Internal field exported, as a type encoding/json reads does.
type customerView struct {
	Name  string
	Notes id.Internal[string]
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// assertPlaceholderOnly fails unless got shows the placeholder and nothing of secret.
func assertPlaceholderOnly(t *testing.T, got string) {
	t.Helper()
	if !strings.Contains(got, placeholder) {
		t.Errorf("got %q; want it to show %q", got, placeholder)
	}
	if containsAny(got, secret, "risk", "owes") {
		t.Errorf("got %q; want nothing of %q", got, secret)
	}
}

func TestInternalRevealsTheValueItWraps(t *testing.T) {
	if got := id.NewInternal(secret).Reveal(); got != secret {
		t.Fatalf("got %q; want %q", got, secret)
	}
}

func TestInternalsTextFormIsThePlaceholder(t *testing.T) {
	if got := id.NewInternal(secret).String(); got != placeholder {
		t.Fatalf("got %q; want %q", got, placeholder)
	}
}

func TestInternalsDebugFormIsThePlaceholder(t *testing.T) {
	if got := id.NewInternal(secret).GoString(); got != placeholder {
		t.Fatalf("got %q; want %q", got, placeholder)
	}
}

// structsHoldingInternal holds secret in an Internal field of each kind of struct.
func structsHoldingInternal() map[string]any {
	return map[string]any{
		"an unexported field": customer{name: "Ada", notes: id.NewInternal(secret)},
		"an exported field":   customerView{Name: "Ada", Notes: id.NewInternal(secret)},
		"a pointer":           &customer{name: "Ada", notes: id.NewInternal(secret)},
	}
}

func TestInternalPrintedInsideAStructShowsThePlaceholder(t *testing.T) {
	for name, v := range structsHoldingInternal() {
		for _, verb := range []string{"%v", "%+v", "%#v", "%s", "%q"} {
			t.Run(name+" with "+verb, func(t *testing.T) {
				assertPlaceholderOnly(t, fmt.Sprintf(verb, v))
			})
		}
	}
}

func TestInternalPrintedInsideAStructAsNumbersHoldsNothingOfTheValue(t *testing.T) {
	for name, v := range structsHoldingInternal() {
		for _, verb := range []string{"%x", "%X", "%d"} {
			t.Run(name+" with "+verb, func(t *testing.T) {
				got := strings.ToLower(fmt.Sprintf(verb, v))
				if containsAny(got, secret, fmt.Sprintf("%x", secret), fmt.Sprintf("%x", "risk")) {
					t.Fatalf("got %q; want nothing of %q", got, secret)
				}
			})
		}
	}
}

func TestInternalSerialisedAsJSONInsideAStructIsThePlaceholder(t *testing.T) {
	got, err := json.Marshal(customerView{Name: "Ada", Notes: id.NewInternal(secret)})
	if err != nil {
		t.Fatalf("got %v; want nil", err)
	}
	if want := `{"Name":"Ada","Notes":"[internal]"}`; string(got) != want {
		t.Fatalf("got %s; want %s", got, want)
	}
}

func TestInternalLoggedThroughSlogIsThePlaceholder(t *testing.T) {
	handlers := map[string]func(*bytes.Buffer) slog.Handler{
		"as JSON": func(b *bytes.Buffer) slog.Handler { return slog.NewJSONHandler(b, nil) },
		"as text": func(b *bytes.Buffer) slog.Handler { return slog.NewTextHandler(b, nil) },
	}
	for name, handler := range handlers {
		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			slog.New(handler(&buf)).Info("looked up", "notes", id.NewInternal(secret))
			assertPlaceholderOnly(t, buf.String())
		})
	}
}
