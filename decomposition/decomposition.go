package decomposition

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

// Decomposition is a hereditary base-b decomposition.
type Decomposition struct {
	// order of the monomes matter:
	// they are sorted from least to most significant
	monomes []monome
}

// New returns the hereditary base-b decomposition of n.
// n must be non negative and b must be at least 2.
func New(b, n int) (Decomposition, error) {
	// n must be non negative
	if n < 0 {
		return Decomposition{}, fmt.Errorf("n must be non negative")
	}

	// base must at least 2
	if b < 2 {
		return Decomposition{}, fmt.Errorf("base must be at least 2")
	}

	return Decomposition{recDecompose(b, n, 0)}.clean(), nil
}

// recDecompose recursively builds the hereditary base-b decomposition of n.
// Monomes are sorted from the least significant to the most significant one.
func recDecompose(b, n, k int) []monome {
	// stop condition: nothing to decompose, return the empty Decomposition
	if n == 0 {
		return nil
	}

	// init decomposition with its least significant monome
	singleton := []monome{
		monome{
			coeff:    n % b,
			base:     b,
			exponent: Decomposition{recDecompose(b, k, 0)},
		},
	}
	return append(singleton, recDecompose(b, n/b, k+1)...)
}

// copyDecomposition returns a deep copy of the decomposition.
func copyDecomposition(d Decomposition) Decomposition {
	copied := make([]monome, len(d.monomes))
	for i, m := range d.monomes {
		copied[i] = copyMonome(m)
	}
	return Decomposition{copied}
}

// IsZero returns true if the decomposition is the decomposition of 0 (in any base).
// The default value of Decomposition is a zero decomposition.
func (d Decomposition) IsZero() bool {
	return len(d.monomes) == 0
}

// isOne returns true if the decomposition is the decomposition of 1 (in any base).
// This applies only to a cleaned decomposition.
func (d Decomposition) isOne() bool {
	return len(d.monomes) == 1 && d.monomes[0].isOne()
}

// clean removes all zero-monomes from the Decomposition.
func (d Decomposition) clean() Decomposition {
	var cleaned []monome
	for _, m := range d.monomes {
		// remove zero monome
		if m.isZero() {
			continue
		}

		// non zero monome:
		// clean its exponent and add it to the cleaned Decomposition
		cleaned = append(cleaned, monome{
			coeff:    m.coeff,
			base:     m.base,
			exponent: m.exponent.clean(),
		})
	}

	return Decomposition{cleaned}
}

// string is a helper for the String and LaTeX methods.
// It returns a human-readable decomposition with
// the given symbols for multiplication and left and right
// groupers around the exponents.
func (d Decomposition) string(times, leftGroup, rightGroup string) string {
	// length of the decomposition
	l := len(d.monomes)
	// if there is no monome, decompostion is zero
	if l == 0 {
		return "0"
	}

	// write all monomes in reverse order
	strMonomes := make([]string, l)
	for i, m := range d.monomes {
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
// Special characters are not escaped so it must not be formatted with the %s verb.
// Instead, the %q one must be used.
func (d Decomposition) LaTeX() string {
	return d.string("\times", "{", "}")
}

// Eval computes and returns the value of the decomposition.
// It returns a *big.Int since huge numbers are expected.
// Note that even if the value of the expression may be huge,
// integer literals in it should remain small enough for type int.
func (d Decomposition) Eval() *big.Int {
	result := big.NewInt(0)
	for _, m := range d.monomes {
		result.Add(result, m.eval())
	}
	return result
}

// IncrementBase returns a new Decomposition with base incremented by one.
// Original decomposition is left unchanged.
func (d Decomposition) IncrementBase() Decomposition {
	incremented := make([]monome, len(d.monomes))
	for i, m := range d.monomes {
		incremented[i] = monome{
			coeff:    m.coeff,
			base:     m.base + 1,
			exponent: m.exponent.IncrementBase(),
		}
	}
	return Decomposition{incremented}
}

// Decrement returns a new Decomposition
// which has been symbolically decremented.
// If the decomposition is already equal to zero it returns the zero Decomposition.
// The original decomposition is left unchanged.
func (d Decomposition) Decrement() Decomposition {
	// if decomposition is zero, return zero
	if d.IsZero() {
		return Decomposition{}
	}

	// to be decremented
	decremented := copyDecomposition(d).monomes

	// find the least significant monome
	// and decrease its coefficient by one
	decremented[0].coeff -= 1

	// prepend all monomes from this one to zero
	// with coefficient (base-1).
	exp := decremented[0].exponent
	var lsms []monome
	for !exp.IsZero() {
		// decrease exponent
		exp = exp.Decrement()

		// new monome is the least significant one.
		// prepend it
		lsms = append(lsms, monome{
			coeff:    decremented[0].base - 1,
			base:     decremented[0].base,
			exponent: copyDecomposition(exp),
		})
	}

	// lsms are in the wrong order
	// reorder them and prepend them to the decremented list of monomes
	for i := 0; i < len(lsms)/2; i++ {
		lsms[i], lsms[len(lsms)-1-i] = lsms[len(lsms)-1-i], lsms[i]
	}
	decremented = append(lsms, decremented...)

	// clean the decomposition
	return Decomposition{decremented}.clean()
}

// monome is an expression of the form 'coeff * base ^ exponent'
// where coeff and base are integers and exponent is a
// hereditary base-b decomposition with base 'base'.
type monome struct {
	coeff, base int
	exponent    Decomposition
}

// copyMonome returns a deep copy (i.e. the exponent is also a copy)
// of the original monome.
func copyMonome(m monome) monome {
	return monome{
		coeff:    m.coeff,
		base:     m.base,
		exponent: copyDecomposition(m.exponent),
	}
}

// isZero returns true if the monome is equal to zero.
func (m monome) isZero() bool { return m.coeff == 0 }

// isOne returns true if the monome is equal to one.
func (m monome) isOne() bool {
	return m.coeff == 1 && m.exponent.IsZero()
}

// string is a helper for the String and LaTeX methods.
// It returns a human readable version of the monome
// with given symbols for the multiplication and
// left and right 'groupers' around the exponent.
func (m monome) string(times, leftGroup, rightGroup string) string {
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

func (m monome) String() string {
	return m.string("*", "(", ")")
}

// eval returns the numeric value of a monome as a *big.Int.
func (m monome) eval() *big.Int {
	c := big.NewInt(int64(m.coeff))
	b := big.NewInt(int64(m.base))

	result := big.NewInt(0)
	result.Exp(b, m.exponent.Eval(), nil)
	result.Mul(c, result)
	return result
}
