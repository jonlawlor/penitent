rel
========

Relational Algebra in Go.

This implements most (?) of the elements of relational algebra, including project, restrict, join, intersect, setdiff, and union.  It also implements some of the common non-relational operations, including groupby, order, insert, and update.  To learn more about relational algebra, C. J. Date's Database in Depth is a great place to start, and is used as the source of terminology in the rel package.

The semantics of this package are very similar to Microsoft's LINQ, although the syntax is different.

Interfaces
==========
The uses of the interfaces defined in the rel package are outlined here.

Relation Interface
------------------

Relations are channels of tuples, and operations on those channels.  The relational algebra operations of project, restrict, join, intersect, setdiff, and union all take at least one relation input and result in a relation output.

Predicate Interface
-------------------
Predicates are used in the restrict operation.
