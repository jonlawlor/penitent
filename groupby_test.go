package rel

import (
	"testing"
)

// tests for projection
func TestGroupBy(t *testing.T) {

	// TODO(jonlawlor): replace with table driven test?
	type r1tup struct {
		PNO int
		Qty int
	}
	type valtup struct {
		Qty int
	}

	// a simple summation
	groupFcn := func(val chan interface{}) interface{} {
		res := valtup{}
		for vi := range val {
			v := vi.(valtup)
			res.Qty += v.Qty
		}
		return res
	}
	wantString := `rel.New([]struct {
 PNO int 
 Qty int 
}{
 {4, 900,  },
 {1, 1300, },
 {2, 700,  },
 {3, 200,  },
})`
	r1 := orders.GroupBy(r1tup{}, valtup{}, groupFcn)
	if r1.GoString() != wantString {
		t.Errorf("orders.Groupby = \"%s\", want \"%s\"", r1.GoString(), wantString)

	}
	return
}