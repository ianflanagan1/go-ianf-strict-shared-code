package equality_test

import (
	"slices"
	"testing"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/core/equality"
)

func flat() (Money, Money) { return Money{100, GBP}, Money{100, GBP} }

func nested() (Order, Order) {
	o := func() Order {
		return Order{
			Items: []LineItem{{"tea", Money{100, GBP}}, {"cake", Money{250, GBP}}, {"jam", Money{300, GBP}}},
			Notes: map[string]Money{"tip": {50, GBP}, "fee": {20, GBP}},
			Tags:  []string{"new", "sale"},
		}
	}
	return o(), o()
}

func handEqual(a, b Order) bool {
	if !slices.Equal(a.Items, b.Items) || !slices.Equal(a.Tags, b.Tags) || len(a.Notes) != len(b.Notes) {
		return false
	}
	for k, v := range a.Notes {
		w, ok := b.Notes[k]
		if !ok || v != w {
			return false
		}
	}
	return true
}

func BenchmarkFlatOperator(b *testing.B) {
	x, y := flat()
	for b.Loop() {
		_ = x == y
	}
}

func BenchmarkFlatEqual(b *testing.B) {
	x, y := flat()
	for b.Loop() {
		_ = equality.Equal(x, y)
	}
}

func BenchmarkNestedHandWritten(b *testing.B) {
	x, y := nested()
	for b.Loop() {
		_ = handEqual(x, y)
	}
}

func BenchmarkNestedEqual(b *testing.B) {
	x, y := nested()
	for b.Loop() {
		_ = equality.Equal(x, y)
	}
}
