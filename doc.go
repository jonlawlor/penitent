// Package rel implements relational algebra, a set of operations on
// sets of tuples which result in relations, as defined by E. F. Codd.
//
// Basics
//
// What folows is a brief introduction to relational algebra.  For a more
// complete introduction, please read C. J. Date's book "Database in Depth".
// This package uses the same terminology.
//
// Relations are sets of named tuples with identical attributes.  The primative
// operations which define the relational algebra are:
//
// Union, which adds two sets together.
//
// Diff, which removes all elements from one set which exist in another.
//
// Restrict, which removes values from a relation that do not satisfy a
// predicate.
//
// Project, which removes zero or more attributes from the tuples the relation
// is defined on.
//
// Rename, which changes the names of the attributes in a relation.
//
// Join, which can multiply two relations together (which may have different
// types of tuples) by returning all combinations of tuples in the two
// relations where all attributes in one relation are equal to the attributes
// in the other where the names are the same.  This is sometimes called a
// natural join.
//
// This package represents tuples as structs with no unexported or anonymous
// fields.  The fields of the struct are the attributes of the tuple it
// represents.
//
// Attributes are strings with some additional methods that are useful for
// constructing predicates and candidate keys.  They have to be valid field
// names in go.
//
// Predicates are functions which take a tuple and return a boolean, and are
// used as an input for Restrict expressions.
//
// Candidate keys are the sets of attributes which define unique tuples in a
// relation.  Every relation has at least one candidate key, because every
// relation only contains unique tuples.  Some relations may contains several
// candidate keys.
//
// Relations in this package can be either literal, such as a relation from a
// map of tuples, or an expression of other relations, such as a join between
// two source relations.
//
// Literal Relations can be defined using the rel.New function.  Given a slice,
// map, or channel of tuples, the New function constructs a new "essential"
// relation, with those values as tuples.  Other packages can create essential
// relations from other sources of data, such as the github.com/jonlawlor/relcsv
// package, or the github.com/jonlawlor/relsql package.
//
// Relational Expressions are generated when one of the methods Project,
// Restrict, Union, Diff, Join, Rename, Map, or GroupBy.  During their
// construction, the rel package checks to see if they can be distributed over
// the source relations that they are being called on, and if so, it attempts
// to push the expressions down the tree of relations as far as they can go,
// with the end goal of getting pushed all the way to the "essential" source
// relations.  In this way, relational expressions can (hopefully) reduce the
// amount of computation done in total and / or done in the go runtime.
//
package rel

// variable naming conventions
//
// r1, r2, r3, ... all represent relations.  If there is an operation which
// has an output relation, the output relation will have the highest number
// after the r.
//
// body, body1, body2, b, b1, b2, ... all represent channels of tuples.
//
// zero, z, z1, z2, ... all represent a tuple's zero value, with defaults in
// all of the fields.
//
// elem, e, e1, e2, ... all represent the reflect.ValueOf(z) with the
// appropriate identification.
//
// tup, tup1, tup2, ... all represent actual tuples going through some
// relational transformation.
//
// rtup, rtup1, rtup2, ... all represent the reflect.ValueOf(tup) with the
// appropriate identification.
