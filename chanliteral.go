// chanLiteral is a relation with underlying data stored in a channel.
// It is intended to be used for general purpose source data that can only be
// queried as a whole, or that has been preprocessed.  It can also be used as
// an adapter to interface with other sources of data.  Basically anything
// that can produce a chan of structs.
//

package rel

import (
	"reflect"
)

// chanLiteral is an implementation of Relation using a channel
type chanLiteral struct {
	// the channel of tuples in the relation
	rbody reflect.Value

	// set of candidate keys
	cKeys CandKeys

	// the type of the tuples contained within the relation
	zero T

	// sourceDistinct indicates if the source chan was already distinct or if a
	// distinct has to be performed when sending tuples
	sourceDistinct bool

	err error
}

// Tuples sends each tuple in the relation to a channel
// note: this consumes the values of the relation, and when it is finished it
// closes the input channel.
func (r *chanLiteral) Tuples(t chan<- T) chan<- struct{} {
	cancel := make(chan struct{})

	if r.Err() != nil {
		close(t)
		return cancel
	}

	if r.sourceDistinct {
		go func(rbody reflect.Value, res chan<- T) {
			resSel := reflect.SelectCase{reflect.SelectRecv, rbody, reflect.Value{}}
			canSel := reflect.SelectCase{reflect.SelectRecv, reflect.ValueOf(cancel), reflect.Value{}}
			cases := []reflect.SelectCase{resSel, canSel}
			for {
				chosen, rtup, ok := reflect.Select(cases)
				if !ok || chosen == 1 {
					// cancel has been closed, so close the results
					// TODO(jonlawlor): include a cancel channel in the rel.chanLiteral
					// struct so that we can continue the cancellation to the data
					// source.
					break
				}
				res <- rtup.Interface()
			}
			close(res)
		}(r.rbody, t)
		return cancel
	}
	// build up a map where each key is one of the tuples.  This consumes
	// memory.
	mem := map[T]struct{}{}
	go func(rbody reflect.Value, res chan<- T) {
		resSel := reflect.SelectCase{reflect.SelectRecv, rbody, reflect.Value{}}
		canSel := reflect.SelectCase{reflect.SelectRecv, reflect.ValueOf(cancel), reflect.Value{}}
		cases := []reflect.SelectCase{resSel, canSel}
		for {
			chosen, rtup, ok := reflect.Select(cases)
			if !ok || chosen == 1 {
				// cancel has been closed, so close the results
				// TODO(jonlawlor): include a cancel channel in the rel.chanLiteral
				// struct so that we can continue the cancellation to the data
				// source.
				break
			}
			tup := rtup.Interface()
			if _, dup := mem[tup]; !dup {
				res <- tup
				mem[tup] = struct{}{}
			}
		}
		close(res)
	}(r.rbody, t)
	return cancel
}

// Zero returns the zero value of the relation (a blank tuple)
func (r *chanLiteral) Zero() T {
	return r.zero
}

// CKeys is the set of candidate keys in the relation
func (r *chanLiteral) CKeys() CandKeys {
	return r.cKeys
}

// GoString returns a text representation of the Relation
func (r *chanLiteral) GoString() string {
	return goStringTabTable(r)
}

// String returns a text representation of the Relation
func (r *chanLiteral) String() string {
	return "Relation(" + HeadingString(r) + ")"
}

// Project creates a new relation with less than or equal degree
// t2 has to be a new type which is a subdomain of r.
func (r1 *chanLiteral) Project(z2 T) Relation {
	if r1.Err() != nil {
		return r1
	}
	att2 := fieldNames(reflect.TypeOf(z2))
	if Deg(r1) == len(att2) {
		// either projection is an error or a no op
		return r1
	} else {
		return &ProjectExpr{r1, z2, nil}
	}
}

// Restrict creates a new relation with less than or equal cardinality
// p has to be a func(tup T) bool where tup is a subdomain of the input r.
// This is a general purpose restrict - we might want to have specific ones for
// the typical theta comparisons or <= <, =, >, >=, because it will allow much
// better optimization on the source data side.
func (r1 *chanLiteral) Restrict(p Predicate) Relation {
	if r1.Err() != nil {
		return r1
	}
	return &RestrictExpr{r1, p, nil}
}

// Rename creates a new relation with new column names
// z2 has to be a struct with the same number of fields as the input relation
// note: we might want to change this into a projectrename operation?  It will
// be tricky to represent this in go's type system, I think.
func (r1 *chanLiteral) Rename(z2 T) Relation {
	if r1.Err() != nil {
		return r1
	}
	return &RenameExpr{r1, z2, nil}
}

// Union creates a new relation by unioning the bodies of both inputs
//
func (r1 *chanLiteral) Union(r2 Relation) Relation {
	if r1.Err() != nil {
		return r1
	}
	if r2.Err() != nil {
		return r2
	}
	return &UnionExpr{r1, r2, nil}
}

// SetDiff creates a new relation by set minusing the two inputs
//
func (r1 *chanLiteral) SetDiff(r2 Relation) Relation {
	if r1.Err() != nil {
		return r1
	}
	if r2.Err() != nil {
		return r2
	}
	return &SetDiffExpr{r1, r2, nil}
}

// Join creates a new relation by performing a natural join on the inputs
//
func (r1 *chanLiteral) Join(r2 Relation, zero T) Relation {
	if r1.Err() != nil {
		return r1
	}
	if r2.Err() != nil {
		return r2
	}
	return &JoinExpr{r1, r2, zero, nil}
}

// GroupBy creates a new relation by grouping and applying a user defined func
//
func (r1 *chanLiteral) GroupBy(t2, vt T, gfcn func(<-chan T) T) Relation {
	if r1.Err() != nil {
		return r1
	}
	return &GroupByExpr{r1, t2, vt, gfcn, nil}
}

// Error returns an error encountered during construction or computation
func (r1 *chanLiteral) Err() error {
	return r1.err
}