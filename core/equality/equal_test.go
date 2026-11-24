package equality_test

import (
	"math"
	"testing"
	"time"

	"github.com/ianflanagan1/go-ianf-strict-shared-code/core/equality"
)

type Currency string

const (
	GBP Currency = "GBP"
	EUR Currency = "EUR"
)

type Money struct {
	Pence    int
	Currency Currency
}

type Price struct {
	Amount any
}

type Reference struct {
	Code string
}

type Ticket struct {
	Code string
}

type LineItem struct {
	Name  string
	Price Money
}

type Order struct {
	Items []LineItem
	Notes map[string]Money
	Tags  []string
}

type Shape interface{ isShape() }

type Circle struct{ Size int }

type Square struct{ Size int }

func (Circle) isShape() {}
func (Square) isShape() {}

type Drawing struct{ Shapes []Shape }

type Appointment struct{ At time.Time }

type hidden struct{ Visible, secret int }

type Handler struct{ Run func() }

type Pointing struct{ To *Money }

type Inbox struct{ Messages chan string }

type pair struct {
	a, b  any
	equal bool
}

func check(t *testing.T, pairs ...pair) {
	t.Helper()
	for _, p := range pairs {
		if got := equality.Equal(p.a, p.b); got != p.equal {
			t.Errorf("Equal(%#v, %#v) = %v; want %v", p.a, p.b, got, p.equal)
		}
		if got := equality.Equal(p.b, p.a); got != p.equal {
			t.Errorf("Equal(%#v, %#v) = %v; want %v", p.b, p.a, got, p.equal)
		}
	}
}

func TestAScalarEqualsTheSameScalar(t *testing.T) {
	check(t, pair{1, 1, true}, pair{1.5, 1.5, true}, pair{true, true, true})
}

func TestAScalarDiffersFromAnotherOfItsType(t *testing.T) {
	check(t, pair{1, 2, false}, pair{1.5, 2.5, false}, pair{true, false, false})
}

func TestAStringEqualsTheSameString(t *testing.T) {
	check(t, pair{"ian", "ian", true})
}

func TestAStringDiffersFromAnotherString(t *testing.T) {
	check(t, pair{"ian", "Ian", false})
}

func TestTwoStringsThatReadAsTheSameNumberDiffer(t *testing.T) {
	check(t, pair{"1", "1.0", false}, pair{"1e3", "1000", false})
}

func TestAnIntegerAFloatAndAStringOfTheSameNumberAllDiffer(t *testing.T) {
	check(t, pair{1, 1.0, false}, pair{1, "1", false}, pair{1.0, "1", false}, pair{int64(1), 1, false})
}

func TestNilDiffersFromEveryEmptyValue(t *testing.T) {
	check(t, pair{nil, "", false}, pair{nil, 0, false}, pair{nil, false, false}, pair{nil, []int{}, false})
}

func TestZeroEqualsNegativeZero(t *testing.T) {
	check(t, pair{0.0, math.Copysign(0, -1), true})
}

func TestAFloatThatIsNotANumberEqualsNothingNotEvenItself(t *testing.T) {
	nan := math.NaN()
	check(t, pair{nan, nan, false}, pair{nan, 0.0, false}, pair{[]float64{nan}, []float64{nan}, false})
}

func TestAnEnumCaseEqualsItself(t *testing.T) {
	check(t, pair{GBP, GBP, true})
}

func TestAnEnumCaseDiffersFromAnotherCase(t *testing.T) {
	check(t, pair{GBP, EUR, false})
}

func TestAnEnumCaseDiffersFromTheValueItHolds(t *testing.T) {
	check(t, pair{GBP, "GBP", false})
}

func TestTwoInstancesOfOneTypeWithEqualFieldsAreEqual(t *testing.T) {
	check(t, pair{Money{100, GBP}, Money{100, GBP}, true})
}

func TestTwoInstancesThatDifferInOneFieldDiffer(t *testing.T) {
	check(t, pair{Money{100, GBP}, Money{101, GBP}, false}, pair{Money{100, GBP}, Money{100, EUR}, false})
}

func TestTwoTypesWithTheSameFieldsDiffer(t *testing.T) {
	check(t, pair{Reference{"A1"}, Ticket{"A1"}, false})
}

func TestTwoInstancesWhoseStringFieldsReadAsTheSameNumberDiffer(t *testing.T) {
	check(t, pair{Reference{"1"}, Reference{"1.0"}, false})
}

func TestTwoInstancesWhoseFieldsHoldTheSameNumberAsDifferentTypesDiffer(t *testing.T) {
	check(t, pair{Price{1}, Price{1.0}, false}, pair{Price{1}, Price{"1"}, false})
}

func TestTwoInstancesHoldingEqualValuesAreEqual(t *testing.T) {
	check(t, pair{LineItem{"tea", Money{100, GBP}}, LineItem{"tea", Money{100, GBP}}, true})
}

func TestTwoInstancesHoldingValuesThatDifferDiffer(t *testing.T) {
	check(t, pair{LineItem{"tea", Money{100, GBP}}, LineItem{"tea", Money{100, EUR}}, false})
}

func TestTwoInstancesHoldingTheSameListInAnotherOrderDiffer(t *testing.T) {
	check(t, pair{Order{Tags: []string{"new", "sale"}}, Order{Tags: []string{"sale", "new"}}, false})
}

func TestAStructDiffersFromAValueThatIsNoStruct(t *testing.T) {
	check(t, pair{Money{}, 0, false}, pair{Money{}, nil, false}, pair{Money{}, []int{}, false})
}

func TestTwoValuesOfOneVariantOfASealedUnionWithEqualFieldsAreEqual(t *testing.T) {
	check(t, pair{Shape(Circle{1}), Shape(Circle{1}), true}, pair{Drawing{[]Shape{Circle{1}}}, Drawing{[]Shape{Circle{1}}}, true})
}

func TestTwoValuesOfOneVariantThatDifferInAFieldDiffer(t *testing.T) {
	check(t, pair{Shape(Circle{1}), Shape(Circle{2}), false})
}

func TestTwoVariantsOfASealedUnionWithTheSameFieldsDiffer(t *testing.T) {
	check(t, pair{Shape(Circle{1}), Shape(Square{1}), false}, pair{Drawing{[]Shape{Circle{1}}}, Drawing{[]Shape{Square{1}}}, false})
}

func TestTwoListsWithEqualElementsInTheSameOrderAreEqual(t *testing.T) {
	check(t, pair{[]int{1, 2, 3}, []int{1, 2, 3}, true}, pair{[3]int{1, 2, 3}, [3]int{1, 2, 3}, true})
}

func TestTwoListsWithTheSameElementsInAnotherOrderDiffer(t *testing.T) {
	check(t, pair{[]int{1, 2, 3}, []int{3, 2, 1}, false}, pair{[3]int{1, 2, 3}, [3]int{3, 2, 1}, false})
}

func TestTwoListsOfDifferentLengthsDiffer(t *testing.T) {
	check(t, pair{[]int{1, 2}, []int{1, 2, 3}, false}, pair{[2]int{1, 2}, [3]int{1, 2, 3}, false})
}

func TestTwoListsOfEqualValuesAreEqual(t *testing.T) {
	check(t, pair{[]Money{{1, GBP}, {2, EUR}}, []Money{{1, GBP}, {2, EUR}}, true})
}

func TestTwoListsOfValuesThatDifferInOneElementDiffer(t *testing.T) {
	check(t, pair{[]Money{{1, GBP}, {2, EUR}}, []Money{{1, GBP}, {3, EUR}}, false})
}

func TestTwoListsHoldingTheSameNumberAsDifferentTypesDiffer(t *testing.T) {
	check(t, pair{[]any{1}, []any{1.0}, false}, pair{[]any{1}, []any{"1"}, false})
}

func TestTwoMapsWithTheSameKeysAndEqualValuesAreEqual(t *testing.T) {
	check(t, pair{map[string]Money{"a": {1, GBP}}, map[string]Money{"a": {1, GBP}}, true})
}

func TestTwoMapsThatDifferInOneValueDiffer(t *testing.T) {
	check(t, pair{map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 3}, false})
}

func TestTwoMapsThatDifferInOneKeyDiffer(t *testing.T) {
	check(t, pair{map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "c": 2}, false})
}

func TestAMapDiffersFromOneWithAKeyMore(t *testing.T) {
	check(t, pair{map[string]int{"a": 1}, map[string]int{"a": 1, "b": 2}, false})
}

func TestAKeyThatHoldsNilDiffersFromAKeyThatIsMissing(t *testing.T) {
	check(t, pair{map[string]any{"a": nil}, map[string]any{}, false}, pair{map[string]any{"a": nil}, map[string]any{"b": nil}, false})
}

func TestTwoInstancesHoldingEqualListsAndMapsOfValuesAreEqual(t *testing.T) {
	order := func() Order {
		return Order{
			Items: []LineItem{{"tea", Money{100, GBP}}, {"cake", Money{250, GBP}}},
			Notes: map[string]Money{"tip": {50, GBP}},
			Tags:  []string{"new"},
		}
	}
	check(t, pair{order(), order(), true})
}

func TestTwoInstancesThatDifferInOneValueDeepInsideDiffer(t *testing.T) {
	one := Order{Items: []LineItem{{"tea", Money{100, GBP}}}, Notes: map[string]Money{"tip": {50, GBP}}}
	other := Order{Items: []LineItem{{"tea", Money{100, GBP}}}, Notes: map[string]Money{"tip": {50, EUR}}}
	check(t, pair{one, other, false})
}

func TestTwoListsOfMapsOfListsAreComparedAllTheWayDown(t *testing.T) {
	check(t,
		pair{[]map[string][]string{{"tags": {"new", "sale"}}}, []map[string][]string{{"tags": {"new", "sale"}}}, true},
		pair{[]map[string][]string{{"tags": {"new", "sale"}}}, []map[string][]string{{"tags": {"new", "sold"}}}, false},
	)
}

// Go has no instances of a struct, only copies, so a pointer is its nearest
// to an object that is no value: it equals only the same address.
func TestAPointerEqualsOnlyTheSameAddress(t *testing.T) {
	money := &Money{100, GBP}
	check(t,
		pair{[]*Money{money}, []*Money{money}, true},
		pair{[]*Money{money}, []*Money{{100, GBP}}, false},
		pair{Pointing{money}, Pointing{money}, true},
		pair{Pointing{money}, Pointing{&Money{100, GBP}}, false},
		pair{Pointing{nil}, Pointing{nil}, true},
	)
}

func TestAChannelEqualsOnlyTheSameChannel(t *testing.T) {
	messages := make(chan string)
	check(t, pair{Inbox{messages}, Inbox{messages}, true}, pair{Inbox{messages}, Inbox{make(chan string)}, false})
}

// A struct with an unexported field keeps data the function may not read, so
// it equals nothing, itself included: two are never found equal for want of
// anything to compare.
func TestAStructWithAnUnexportedFieldEqualsNothingNotEvenItself(t *testing.T) {
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	h := hidden{Visible: 1, secret: 2}
	check(t,
		pair{at, at, false},
		pair{h, h, false},
		pair{Appointment{at}, Appointment{at}, false},
		pair{Appointment{at}, Appointment{at.AddDate(1, 0, 0)}, false},
	)
}

// Go can tell two functions apart only from nil, so a function that isn't
// nil equals nothing, itself included.
func TestAFunctionEqualsNothingNotEvenItself(t *testing.T) {
	double := func() {}
	check(t,
		pair{Handler{double}, Handler{double}, false},
		pair{Handler{double}, Handler{func() {}}, false},
		pair{Handler{nil}, Handler{nil}, true},
		pair{Handler{double}, Handler{nil}, false},
	)
}

func TestANilListEqualsAnEmptyList(t *testing.T) {
	check(t, pair{[]int(nil), []int{}, true}, pair{Order{}, Order{Items: []LineItem{}, Notes: map[string]Money{}, Tags: []string{}}, true})
}

func TestNilEqualsNil(t *testing.T) {
	check(t, pair{nil, nil, true}, pair{Shape(nil), Shape(nil), true}, pair{(*Money)(nil), (*Money)(nil), true})
}

func TestATypedNilPointerDiffersFromNil(t *testing.T) {
	check(t, pair{(*Money)(nil), nil, false})
}
