// chanLiteral is a relation with underlying data stored in a channel.
// It is intended to be used for general purpose source data that can only be
// queried as a whole, or that has been preprocessed.  It can also be used as
// an adapter to interface with other sources of data.  Basically anything
// that can produce a chan of structs.

package rel

import (
	"reflect"
)

// chanLiteral is an implementation of Relation using a channel.
// One issue is that it does not pass along any cancellation.  Maybe it would
// be better if there were two chanLiterals: one that allowed cancellation and
// one that did not.
type chanLiteral struct {
	// the channel of tuples in the relation
	rbody reflect.Value

	// set of candidate keys
	cKeys CandKeys

	// the type of the tuples contained within the relation
	zero interface{}

	// sourceDistinct indicates if the source chan was already distinct or if a
	// distinct has to be performed when sending tuples
	sourceDistinct bool

	err error
}

// TupleChan sends each tuple in the relation to a channel
// note: this consumes the values of the relation, and when it is finished it
// closes the input channel.
func (r1 *chanLiteral) TupleChan(t interface{}) chan<- struct{} {
	cancel := make(chan struct{})
	// reflect on the channel
	chv := reflect.ValueOf(t)
	err := EnsureChan(chv.Type(), r1.zero)
	if err != nil {
		r1.err = err
		return cancel
	}
	if r1.err != nil {
		chv.Close()
		return cancel
	}
	if r1.sourceDistinct {
		go func(rbody, res reflect.Value) {
			// input channel
			sourceSel := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: rbody}
			canSel := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(cancel)}
			inCases := []reflect.SelectCase{canSel, sourceSel}

			// output channel
			resSel := reflect.SelectCase{Dir: reflect.SelectSend, Chan: res}

			for {
				chosen, tup, ok := reflect.Select(inCases)

				// cancel has been closed, so close the results
				if chosen == 0 {
					return
				}
				if !ok {
					// source channel was closed
					break
				}
				resSel.Send = tup
				outCases := []reflect.SelectCase{canSel, resSel}
				chosen, _, _ = reflect.Select(outCases)

				if chosen == 0 {
					// cancel has been closed, so close the results
					return
				}
			}
			res.Close()
		}(r1.rbody, chv)
		return cancel
	}
	// Build up a map where each key is one of the tuples.  This consumes
	// memory.  It may be better to have a seperate distinctExpr Relation,
	// which can transform an input relation to be distinct.  On the other
	// hand, that would produce an additional level of reflection and chan
	// communication.
	mem := map[interface{}]struct{}{}
	go func(rbody, res reflect.Value) {
		sourceSel := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: rbody}
		canSel := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(cancel)}
		inCases := []reflect.SelectCase{canSel, sourceSel}

		resSel := reflect.SelectCase{Dir: reflect.SelectSend, Chan: res}

		for {
			chosen, rtup, ok := reflect.Select(inCases)
			// cancel has been closed, so close the results
			if chosen == 0 {
				return
			}
			if !ok {
				// source channel was closed
				break
			}

			tup := interface{}(rtup.Interface())
			if _, dup := mem[tup]; !dup {
				resSel.Send = rtup
				chosen, _, ok = reflect.Select([]reflect.SelectCase{canSel, resSel})
				if chosen == 0 {
					// cancel has been closed, so close the results
					return
				}
				mem[tup] = struct{}{}
			}
		}
		res.Close()
	}(r1.rbody, chv)
	return cancel
}

// Zero returns the zero value of the relation (a blank tuple)
func (r1 *chanLiteral) Zero() interface{} {
	return r1.zero
}

// CKeys is the set of candidate keys in the relation
func (r1 *chanLiteral) CKeys() CandKeys {
	return r1.cKeys
}

// GoString returns a text representation of the Relation
func (r1 *chanLiteral) GoString() string {
	return goStringTabTable(r1)
}

// String returns a text representation of the Relation
func (r1 *chanLiteral) String() string {
	return "Relation(" + HeadingString(r1) + ")"
}

// Project creates a new relation with less than or equal degree
// t2 has to be a new type which is a subdomain of r.
func (r1 *chanLiteral) Project(z2 interface{}) Relation {
	return NewProject(r1, z2)
}

// Restrict creates a new relation with less than or equal cardinality
// p has to be a func(tup T) bool where tup is a subdomain of the input r.
func (r1 *chanLiteral) Restrict(p Predicate) Relation {
	return NewRestrict(r1, p)
}

// Rename creates a new relation with new column names
// z2 has to be a struct with the same number of fields as the input relation
func (r1 *chanLiteral) Rename(z2 interface{}) Relation {
	return NewRename(r1, z2)
}

// Union creates a new relation by unioning the bodies of both inputs
func (r1 *chanLiteral) Union(r2 Relation) Relation {
	return NewUnion(r1, r2)
}

// Diff creates a new relation by set minusing the two inputs
func (r1 *chanLiteral) Diff(r2 Relation) Relation {
	return NewDiff(r1, r2)
}

// Join creates a new relation by performing a natural join on the inputs
func (r1 *chanLiteral) Join(r2 Relation, zero interface{}) Relation {
	return NewJoin(r1, r2, zero)
}

// GroupBy creates a new relation by grouping and applying a user defined func
func (r1 *chanLiteral) GroupBy(t2, gfcn interface{}) Relation {
	return NewGroupBy(r1, t2, gfcn)
}

// Map creates a new relation by applying a function to tuples in the source
func (r1 *chanLiteral) Map(mfcn interface{}, ckeystr [][]string) Relation {
	return NewMap(r1, mfcn, ckeystr)
}

// Err returns an error encountered during construction or computation
func (r1 *chanLiteral) Err() error {
	return r1.err
}
