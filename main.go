package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/batiazinga/goodstein/decomposition"
)

var (
	it     = flag.Int("it", 10, "maximum number of iterations")
	latex  = flag.Bool("latex", false, "if true, results are valid LaTeX commands")
	header = flag.Bool("header", true, "if true, a header is displayed")
)

func main() {
	flag.Parse()

	// check command validity

	// check number of iterations
	if *it < 0 {
		log.Print("it must be positive")
		os.Exit(1)
	}

	// check number of arguments
	if len(flag.Args()) != 1 {
		log.Print("expecting one and only one argument")
		os.Exit(1)
	}

	// validate argument
	n, err := strconv.ParseInt(flag.Arg(0), 10, 64)
	if err != nil {
		log.Printf("invalid argument, expecting integer: %v", err)
		os.Exit(1)
	}
	// it must be positive too
	if n < 0 {
		log.Print("invalid argument, expecting positive integer")
		os.Exit(1)
	}

	// compute first decomposition
	b := 2 // initial base
	// compute hereditary base-2 decomposition of n
	d, err := decomposition.New(b, int(n))
	if err != nil {
		log.Printf("error while computing hereditary base-%b decomposition of %v: %v", b, n, err)
		os.Exit(2)
	}

	// print header (or not)
	if *header {
		fmt.Fprintln(os.Stdout, "iteration base value decomposition")
	}

	// start iterations
	for i := 0; i < *it; i++ {
		// print result to stdout
		var strDecomposition string
		if *latex {
			strDecomposition = d.LaTeX()
		} else {
			strDecomposition = d.String()
		}
		fmt.Fprintf(os.Stdout, "%v %v %v %q\n", i, b, d.Eval(), strDecomposition)

		// if decomposition is zero, stop
		if d.IsZero() {
			os.Exit(0)
		}

		// increment base and remove one
		b++ // for reporting only
		d = decomposition.Decrement(decomposition.IncrementBase(d))
	}
}
