package ids

import "github.com/ianflanagan1/go-adapters/values/id"

// Fake stands in for System in tests, handing out the IDs it was given, in
// order.
type Fake struct {
	rowIDs    []id.RowID
	publicIDs []id.PublicID
}

// NewFake returns a fake that hands out rowIDs and publicIDs, in order.
func NewFake(rowIDs []id.RowID, publicIDs []id.PublicID) *Fake {
	return &Fake{rowIDs: rowIDs, publicIDs: publicIDs}
}

// NewRowID returns the next row ID f was given. It panics when f has none
// left, since the test gave it too few.
func (f *Fake) NewRowID() id.RowID {
	if len(f.rowIDs) == 0 {
		panic("ids: the fake has no row IDs left; give NewFake more")
	}
	next := f.rowIDs[0]
	f.rowIDs = f.rowIDs[1:]
	return next
}

// NewPublicID returns the next public ID f was given. It panics when f has
// none left, since the test gave it too few.
func (f *Fake) NewPublicID() id.PublicID {
	if len(f.publicIDs) == 0 {
		panic("ids: the fake has no public IDs left; give NewFake more")
	}
	next := f.publicIDs[0]
	f.publicIDs = f.publicIDs[1:]
	return next
}
