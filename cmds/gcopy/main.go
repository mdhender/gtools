// gtools - a collection of grammar manipulation tools
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

// Package main implements the command line tool gcopy
// to copy (and pretty print) a BNF grammar.
package main

import (
	"flag"
	"github.com/mdhender/gtools"
	"log"
)

// written by Douglas Jones, July 2013,
// based on pieces of cruncher, written in Pascal by Douglas Jones, March 1990
// rewritten in C, Jan 2007

// main program to copy a grammar
func main() {
	var input string
	flag.StringVar(&input, "input", input, "grammar to process")
	flag.Parse()

	if input != "" {
		if err := gtools.SetStdin(input); err != nil {
			log.Fatal(err)
		}
	}

	readg()
	writeg()
}

func readg() {
	gtools.ReadGrammar()
}

func writeg() {
	gtools.WriteGrammar()
}
