package main

import (
	"fmt"
	"os"

	"github.com/batiazinga/goodstein/decomposition"
)

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
		d, err := decomposition.New(input.base, input.n)
		if err != nil {
			os.Exit(1)
		}
		fmt.Printf("%v = %s\n", d.Eval(), d.LaTeX())
	}
}
