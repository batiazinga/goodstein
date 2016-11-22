package main

import "strconv"

// A simple AST for the hereditary base-n decomposition problem

// expr is the interface of the AST
//
// Note that AST must be a tree,
type expr interface {
	String() string

	// LaTeX returns the expression in a string
	// which can be copied to a LaTeX file.
	LaTeX() string

	// isZero return true if the expression is the 0 integer literal.
	isZero() bool

	// isOne returns true if the expression is the 1 integer literal.
	isOne() bool

	// clean simplifies an expr by removing spurious '0+...', '0*...', '1*...' and so one.
	// A cleaned copy is returned and the original expr is not modified.
	clean() expr

	// iSubstitute replace all n-literals by m-literals.
	// The original expr is affected (in-place substitution).
	iSubstitute(n, m int)
}

type operation int

const (
	sum operation = iota
	prod
	power
)

type literal struct {
	value int
}

func (l *literal) String() string { return strconv.FormatInt(int64(l.value), 10) }
func (l *literal) LaTeX() string  { return l.String() }

func (l *literal) isZero() bool { return l.value == 0 }
func (l *literal) isOne() bool  { return l.value == 1 }

func (l *literal) clean() expr { return &literal{l.value} }

func (l *literal) iSubstitute(n, m int) {
	if l.value == n {
		l.value = m
	}
}

type binary struct {
	left, right expr
	op          operation
}

func (b *binary) String() string {
	switch b.op {
	case sum:
		return b.left.String() + " + " + b.right.String()
	case prod:
		return b.left.String() + " * " + b.right.String()
	case power:
		// do we need parenthesis for the exponent expression?
		if _, cast := b.right.(*literal); cast {
			// no parenthesis
			return b.left.String() + "^" + b.right.String()
		}
		// parenthesis
		return b.left.String() + "^(" + b.right.String() + ")"
	default:
		return "invalid operation"
	}
}

func (b *binary) LaTeX() string {
	switch b.op {
	case sum:
		return b.left.LaTeX() + " + " + b.right.LaTeX()
	case prod:
		return b.left.LaTeX() + " \\times " + b.right.LaTeX()
	case power:
		return b.left.LaTeX() + "^{" + b.right.LaTeX() + "}"
	default:
		return "invalid operation"
	}
}

func (b *binary) isZero() bool { return false }
func (b *binary) isOne() bool  { return false }

func (b *binary) clean() expr {
	// clean left and right terms:
	left := b.left.clean()
	right := b.right.clean()

	// clean the binary itself
	switch b.op {

	case sum:
		if left.isZero() {
			// 0+n=n
			return right
		}
		if right.isZero() {
			// n+0=n
			return left
		}
		// n+m=n+m... no further simplification
		return &binary{
			op:    sum,
			left:  left,
			right: right,
		}

	case prod:
		if left.isZero() || right.isZero() {
			// 0*n=0
			return &literal{0}
		}
		if left.isOne() {
			// 1*n=n
			return right
		}
		if right.isOne() {
			// n*1=n
			return left
		}
		// n*m=n*m... no further simplification
		return &binary{
			op:    prod,
			left:  left,
			right: right,
		}

	case power:
		if right.isZero() {
			// n^0=1
			return &literal{1}
		}
		if right.isOne() {
			// n^1=n
			return left
		}
		// n^m=n^m... no further simplification
		return &binary{
			op:    power,
			left:  left,
			right: right,
		}
	}

	return nil // this should never happen

}

func (b *binary) iSubstitute(n, m int) {
	b.left.iSubstitute(n, m)
	b.right.iSubstitute(n, m)
}
