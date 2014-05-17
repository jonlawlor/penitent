package rel

import (
	"testing"
)

// test creation of relations, including tests to determine the cost of
// representing slices of structs as relations instead of native.
// TODO(jonlawlor): tests involving []struct, map[struct], and chan struct
// representations

// type of the example relations
type exTup2 struct {
	Foo int
	Bar string
}

// exampleRel2 creates an example relation with given cardinality
// and degree 2.
func exampleRel2(c int) (recs []exTup2) {
	for i := 0; i < c; i++ {
		recs = append(recs, exTup2{i, "test"})
	}
	return
}

func BenchmarkSimpleNewTiny(b *testing.B) {
	// test the time it takes to make a new relation with a given size
	exRel := exampleRel2(10)
	ck := [][]string{[]string{"foo"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(exRel, ck)
	}
}
func BenchmarkNonDistinctNewTiny(b *testing.B) {
	// test the time it takes to make a new relation with a given size,
	// but without any candidate keys.  The New function will run
	// a distinct on the input data.
	exRel := exampleRel2(10)
	ck := [][]string{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(exRel, ck)
	}
}

func BenchmarkSimpleNewSmall(b *testing.B) {
	// test the time it takes to make a new relation with a given size
	exRel := exampleRel2(1000)
	ck := [][]string{[]string{"foo"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(exRel, ck)
	}
}
func BenchmarkNonDistinctNewSmall(b *testing.B) {
	// test the time it takes to make a new relation with a given size,
	// but without any candidate keys.  The New function will run
	// a distinct on the input data.
	exRel := exampleRel2(1000)
	ck := [][]string{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(exRel, ck)
	}
}

func BenchmarkSimpleNewMedium(b *testing.B) {
	// test the time it takes to make a new relation with a given size
	exRel := exampleRel2(100000)
	ck := [][]string{[]string{"foo"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(exRel, ck)
	}
}
func BenchmarkNonDistinctNewMedium(b *testing.B) {
	// test the time it takes to make a new relation with a given size,
	// but without any candidate keys.  The New function will run
	// a distinct on the input data.
	exRel := exampleRel2(100000)
	ck := [][]string{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(exRel, ck)
	}
}

func BenchmarkSimpleNewLarge(b *testing.B) {
	// test the time it takes to make a new relation with a given size
	exRel := exampleRel2(10000000)
	ck := [][]string{[]string{"foo"}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(exRel, ck)
	}
}
func BenchmarkNonDistinctNewLarge(b *testing.B) {
	// test the time it takes to make a new relation with a given size,
	// but without any candidate keys.  The New function will run
	// a distinct on the input data.
	exRel := exampleRel2(10000000)
	ck := [][]string{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(exRel, ck)
	}
}

// test the degrees
func TestDeg(t *testing.T) {
	fix := []struct {
		name string
		in   int
		out  int
	}{
		{"suppliers", Deg(suppliers), 4},
		{"parts", Deg(parts), 5},
		{"orders", Deg(orders), 3},
	}
	for i, dt := range fix {
		if dt.in != dt.out {
			t.Errorf("%d. %s.Deg() => %d, want %d", i, dt.name, dt.in, dt.out)
		}
	}
}

func TestCard(t *testing.T) {
	fix := []struct {
		name string
		in   int
		out  int
	}{
		{"suppliers", Card(suppliers), 5},
		{"parts", Card(parts), 6},
		{"orders", Card(orders), 12},
	}
	for i, dt := range fix {
		if dt.in != dt.out {
			t.Errorf("%d. %s.Card() => %d, want %d", i, dt.name, dt.in, dt.out)
		}
	}
}
