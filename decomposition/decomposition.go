package decomposition

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// monome is an expression of the form 'coeff * base ^ exponent'
// where coeff and base are integers and exponent is a
// hereditary base-(base) decomposition.
type monome struct {
	coeff, base int
	exponent    Decomposition
}

// copyMonome returns a deep copy (i.e. the exponent is also a copy)
// of the original monome.
func copyMonome(m *monome) *monome {
	return &monome{
		coeff:    m.coeff,
		base:     m.base,
		exponent: copyDecomposition(m.exponent),
	}
}

// isZero returns true if the monome is equal to zero.
func (m *monome) isZero() bool { return m.coeff == 0 }

// isOne returns true if the monome is equal to one.
func (m *monome) isOne() bool {
	return m.coeff == 1 && m.exponent.IsZero()
}

// string is helper for the String and LaTeX methods.
// It returns a human readable version of the monome
// with given symbols for the multiplication and
// left and right 'groupers' around the exponent.
func (m *monome) string(times, leftGroup, rightGroup string) string {
	// if monome is zero, just return 0
	if m.isZero() {
		return "0"
	}

	// elementary blocks
	strCoeff := strconv.FormatInt(int64(m.coeff), 10)
	strBase := strconv.FormatInt(int64(m.base), 10)
	times = " " + times + " "

	switch {
	case m.exponent.IsZero():
		// base ^ exponent is one, so monome is equal to its coeff
		return strCoeff

	case m.exponent.isOne():
		// base ^ exponent is base
		if m.coeff == 1 {
			// 1 times base is useless, just return the base
			return strBase
		}
		// result is coeff times base
		return strCoeff + times + strBase

	default:
		// general case for the base ^ exponent part
		result := strBase + " ^ " + leftGroup + m.exponent.string(times, leftGroup, rightGroup) + rightGroup
		if m.coeff == 1 {
			// 1 times ... is useless
			return result
		}
		// most general case
		return strCoeff + times + result
	}
}

// eval returns the numeric value of a monome as a big int.
func (m *monome) eval() *big.Int {
	c := big.NewInt(int64(m.coeff))
	b := big.NewInt(int64(m.base))

	result := big.NewInt(0)
	result.Exp(b, m.exponent.Eval(), nil)
	result.Mul(c, result)
	return result
}

// Decomposition is slice of (pointer to) monomes.
// Order of the monomes matters.
type Decomposition []*monome

// New returns the hereditary base-b decomposition of n.
// n must be non negative and base must be at least 2.
func New(base, n int) (Decomposition, error) {
	// n must be non negative
	if n < 0 {
		return nil, fmt.Errorf("n must be non negative")
	}

	// base must at least 2
	if base < 2 {
		return nil, fmt.Errorf("base must be at least 2")
	}

	return recDecompose(base, n, 0).clean(), nil
}

// recDecompose recursively builds the hereditary base-b decomposition of n.
// Monomes are sorted from the least significant to the most significant one.
func recDecompose(b, n, k int) Decomposition {
	// stop condition: nothing to decompose, return the nil Decomposition
	if n == 0 {
		return nil
	}

	singleton := Decomposition([]*monome{
		&monome{
			coeff:    n % b,
			base:     b,
			exponent: recDecompose(b, k, 0),
		},
	})
	return append(singleton, recDecompose(b, n/b, k+1)...)
}

// copyDecomposition returns a deep copy of the decomposition.
func copyDecomposition(d Decomposition) Decomposition {
	copied := make(Decomposition, 0, len(d))
	for _, m := range d {
		copied = append(copied, copyMonome(m))
	}
	return copied
}

// IsZero returns true if the decomposition is the decomposition of 0 (in any base).
// This applies only to a cleaned decomposition.
func (d Decomposition) IsZero() bool {
	return len(d) == 0
}

// isOne returns true if the decomposition is the decomposition of 1 (in any base).
// This applies only to a cleaned decomposition.
func (d Decomposition) isOne() bool {
	return len(d) == 1 && d[0].isOne()
}

// clean removes all zero-monomes from the Decomposition.
func (d Decomposition) clean() Decomposition {
	var cleaned Decomposition
	for _, m := range d {
		// remove zero monome
		if m.isZero() {
			continue
		}

		// non zero monome:
		// clean its exponent and add it to the cleaned Decomposition
		cleaned = append(cleaned, &monome{
			coeff:    m.coeff,
			base:     m.base,
			exponent: m.exponent.clean(),
		})
	}

	return cleaned
}

// string is a helper for the String and LaTeX methods.
// It returns a human-readable decomposition with
// the given symbols for multiplication and left and right
// groupers around the exponents.

func (d Decomposition) string(times, leftGroup, rightGroup string) string {
	// length of the decomposition
	l := len(d)
	// if there is no monome, decompostion is zero
	if l == 0 {
		return "0"
	}

	// write all monomes in reverse order
	strMonomes := make([]string, l)
	for i, m := range d {
		strMonomes[l-1-i] = m.string(times, leftGroup, rightGroup)
	}
	return strings.Join(strMonomes, " + ")
}

// String returns a human readable decomposition
// where most significant monomes lie on the left
// and least significant ones on the right.
// The decomposition does not contain any spurrious '0+...', '0*...' expressions.
func (d Decomposition) String() string {
	return d.string("*", "(", ")")
}

// LaTeX is similar to String but it returns a valid LaTeX command.
func (d Decomposition) LaTeX() string {
	return d.string("\\times", "{", "}")
}

// Eval computes and returns the value of the decomposition.
// It returns a big int since huge numbers are expected.
// Note that even if the value of the expression may be huge,
// integer literals in it should remain small enough for type int.
func (d Decomposition) Eval() *big.Int {
	result := big.NewInt(0)
	for _, m := range d {
		result.Add(result, m.eval())
	}
	return result
}

// IncrementBase returns a new Decomposition with base incremented by one.
func IncrementBase(d Decomposition) Decomposition {
	incremented := make(Decomposition, 0, len(d))
	for _, m := range d {
		incremented = append(incremented, &monome{
			coeff:    m.coeff,
			base:     m.base + 1,
			exponent: IncrementBase(copyDecomposition(m.exponent)),
		})
	}
	return incremented
}

// Decrement returns a new Decomposition
// which has been symbolically decremented.
// If the decomposition is equal to zero, it returns the zero Decomposition
func Decrement(d Decomposition) Decomposition {
	// if decomposition is zero, return zero
	if d.IsZero() {
		return nil
	}

	// to be decremented
	decremented := copyDecomposition(d)

	// find the least significant monome
	lsm := decremented[0]
	// and decrease its coefficient by one
	lsm.coeff -= 1

	// append all monomes from this one to zero
	// with coefficient (base-1).
	exp := lsm.exponent
	for !exp.IsZero() {
		// decrease exponent
		exp = Decrement(exp)

		// new monome is the least significant one.
		// prepend it
		decremented = append([]*monome{&monome{
			coeff:    lsm.base - 1,
			base:     lsm.base,
			exponent: copyDecomposition(exp),
		}}, decremented...)
	}

	// clean the decomposition
	return decremented.clean()
}
