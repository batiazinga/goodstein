package decomposition

import (
	"fmt"
	"math/big"
	"testing"
)

func Example() {
	// decomposition of 0 (in any base)
	d0, _ := New(2, 0)
	fmt.Println(d0.String())

	// decomposition of 1 (in any base)
	d1, _ := New(3, 1)
	fmt.Println(d1.String())

	// base-2 decomposition of 10
	d2_10, _ := New(2, 10)
	fmt.Println(d2_10.String())

	// base-3 decomposition of 10
	d3_10, _ := New(3, 10)
	fmt.Println(d3_10.String())

	// Output:
	// 0
	// 1
	// 2 ^ (2 + 1) + 2
	// 3 ^ (2) + 1
}

func ExampleDecomposition_LaTeX() {
	// decomposition of 0 (in any base)
	d0, _ := New(2, 0)
	fmt.Println(d0.LaTeX())

	// decomposition of 1 (in any base)
	d1, _ := New(3, 1)
	fmt.Println(d1.LaTeX())

	// base-2 decomposition of 10
	d2_10, _ := New(2, 10)
	fmt.Println(d2_10.LaTeX())

	// base-3 decomposition of 10
	d3_10, _ := New(3, 10)
	fmt.Println(d3_10.LaTeX())

	// Output:
	// 0
	// 1
	// 2 ^ {2 + 1} + 2
	// 3 ^ {2} + 1
}

// unit tests for monomes

type goldenMonome struct {
	m             monome
	isZero, isOne bool
	value         *big.Int
}

var goldenMonomes = []goldenMonome{
	{
		m: monome{
			coeff:    0,
			base:     2,
			exponent: Decomposition{},
		},
		isZero: true,
		isOne:  false,
		value:  big.NewInt(0),
	},
	{
		m: monome{
			coeff:    1,
			base:     2,
			exponent: Decomposition{},
		},
		isZero: false,
		isOne:  true,
		value:  big.NewInt(1),
	},
	{
		m: monome{
			coeff:    2,
			base:     3,
			exponent: Decomposition{},
		},
		isZero: false,
		isOne:  false,
		value:  big.NewInt(2),
	},
}

func TestMonomeIsZero(t *testing.T) {
	for _, g := range goldenMonomes {
		if g.m.isZero() != g.isZero {
			t.Errorf("wrong 'isZero' for %q", g.m)
		}
	}
}
func TestMonomeIsOne(t *testing.T) {
	for _, g := range goldenMonomes {
		if g.m.isOne() != g.isOne {
			t.Errorf("wrong 'isOne' for %q", g.m)
		}
	}
}
func TestMonomeEval(t *testing.T) {
	for _, g := range goldenMonomes {
		if g.m.eval().Cmp(g.value) != 0 {
			t.Errorf("wrong value %v for %q", g.m.eval(), g.m)
		}
	}
}
