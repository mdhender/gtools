// gtools - a collection of grammar manipulation tools
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package gtools

// /* grammar.h */
//
// /* written by Douglas Jones, June 2010,
//    based on pieces of cruncher,
//      written in Pascal by Douglas Jones, March 1990
//      rewritten in C, Jan 2007
// */

// all data structures used to represent the grammar are here
// there is no use of serious data abstraction in this code

// Configuration constants

// Limit on the number of nonterminal characters
const STRINGLIMIT = 10_000

// Limit on the number of chars in any symbol (T or NT)
const SYMLEN = 50

// Configurable syntactic details
const COMMENT = '#'

// #define RULESYM "::="
var RULESYM = "::="

// note:  RULESYM is for writeg, while readg accepts any of :, =, := or ::=

// Types

// Handle on strings in the string table (0 .. STRINGLIMIT)
type STRINGPT int

type STYPE int

const (
	ISEMPTY    STYPE = iota // known to be equivalent to the empty symbol
	CANBEEMPTY              // may be empty, marks optional syntactic elements
	TOUCHED                 // is reachable from the root of the grammar
	UNTOUCHED               // has not yet been reached from the root
	STYPES                  // NOT A TYPE, rather, the number of types permitted
)

// Pointer types

type PSYMBOL *symbol

type PPRODUCTION *production

type PELEMENT *element

type symbol struct {
	name    STRINGPT    // symbols have names
	next    PSYMBOL     // symbols may occur in lists of symbols
	data    PPRODUCTION // the head of the list of productions
	state   STYPE       // the state of this symbol
	starter PELEMENT    // the head of the terminal list in the start set
	follows PELEMENT    // the head of the terminal list in the follow set
	line    int         // source line number on which symbol first seen
}

type production struct {
	next    PPRODUCTION // each symbol has a linked list of productions
	data    PELEMENT    // the body of a production is a list of elements
	state   STYPE       // the state of the production rule
	starter PELEMENT    // the head of the terminal list in the start set
	ender   PELEMENT    // the head of the terminal list in the follow set
	line    int         // source line number on which production starts
}

type element struct {
	next PELEMENT // each production is a list of elements
	data PSYMBOL  // an element is a handle on a symbol, never NULL
	line int      // source line number on which element occurs
}

// storage allocators

func NEWSYMBOL() PSYMBOL {
	return &symbol{}
}

func NEWPRODUCTION() PPRODUCTION {
	return &production{}
}

func NEWELEMENT() PELEMENT {
	return &element{}
}

// predicates

func TERMINAL(ps PSYMBOL) bool {
	// #define TERMINAL(s)    ((s)->data == NULL)
	return ps.data == nil
}

func NONTERMINAL(ps PSYMBOL) bool {
	// #define NONTERMINAL(s) ((s)->data != NULL)
	return ps.data != nil
}

// global variables

// the string table
var stringtab [STRINGLIMIT]byte

// the lowest free element in the string table
var stringlim STRINGPT

// head and address of null pointer for the main list of all symbols
var symlist PSYMBOL
var symlistend *PSYMBOL

// Identity of the empty symbol
var emptypt PSYMBOL

// Identity of the distinguished symbol in the grammar
var head PSYMBOL
