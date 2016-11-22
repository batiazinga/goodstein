package main

// decompose returns the hereditary base-(base) decomposition of n:
//
//     https://en.wikipedia.org/wiki/Goodstein's_theorem#Hereditary_base-n_notation.
//
// The result is an AST like
//
//     ( ( (...) + n_2 * b^{2} ) + n_1 * b ) + n_0
//
// where all unecessary operations have been removed (0+..., 1*..., ...).
func decompose(base, n int) expr {
	// dumb decomposition
	e := recDecompose(base, n, 0)

	// clean result and return
	return e.clean()
}

// recDecompose recursively builds the AST of the hereditary base-(base)
// decomposition of n.
func recDecompose(base, n, k int) expr {
	// stop condition:
	if n == 0 {
		return &literal{0}
	}

	return &binary{
		op:   sum,
		left: recDecompose(base, n/base, k+1),
		right: &binary{
			op:   prod,
			left: &literal{n % base},
			right: &binary{
				op:    power,
				left:  &literal{base},
				right: recDecompose(base, k, 0),
			},
		},
	}

}
