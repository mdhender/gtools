// gtools - a collection of grammar manipulation tools
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

// Package main implements the command line tool gsqueeze
// to eliminate redundancy from a BNF grammar.
package main

import (
	"flag"
	"github.com/mdhender/gtools"
	"log"
)

func main() {
	var input string
	flag.StringVar(&input, "input", input, "grammar to process")
	flag.Parse()

	if input != "" {
		if err := gtools.SetStdin(input); err != nil {
			log.Fatal(err)
		}
	}

	gtools.ReadGrammar()
	gtools.Squeeze()
	gtools.WriteGrammar()
}
