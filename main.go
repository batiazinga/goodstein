package main

import "fmt"

var inputs = []struct{ base, n int }{
	{2, 5},
	{2, 10},
	{2, 1023},
	{2, 1024},
	{3, 5},
	{3, 10},
	{3, 1023},
	{3, 1024},
}

func main() {
	for _, input := range inputs {
		fmt.Printf("%v = %s\n", input.n, decompose(input.base, input.n).LaTeX())
	}
}
