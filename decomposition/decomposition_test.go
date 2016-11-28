package decomposition

import "fmt"

func ExampleString() {
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

func ExampleLaTeX() {
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
