// gtools - a collection of grammar manipulation tools
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package gtools

// written by Douglas Jones, July 2010,
// based on pieces of cruncher, written in Pascal by Douglas Jones, March 1990
// rewritten in C, Jan 2007

// OUTPUT FORMAT:
// The distinguished symbol from which the grammar hangs is always output
// first, followed the designated empty symbol, if one was specified.
//    |
//    |> <distinguished>
//    |/ <empty>
//    |
// Terminal and nonterminal symbols are preserved unchanged, in whatever
// form they were input.  Rules are output as follows:
//    |
//    |<nonterminal> ::= <a> 'b' <c> 'd'
//    |                  <continuation>
//    |
// The delimiter ::= was used in the original Algol 60 report.  Indenting
// indicated continuation.  Multiple rules that reduce to the same
// nonterminal are output in succession, using the vertical bar to separate
// successive rules instead of re-stating the lefthand side of the rule.
//    |
//    |<nonterminal> ::= <a> 'b'
//    |               |  <c> 'd'
//    |
// If the start set and follow set of a nonterminal have been computed, these
// are included in the output grammar as comments after all the rules for
// that nonterminal.
//    |
//    |<nonterminal> ::= <a> 'b'
//    |               |  <c> 'd'
//    |\ start:  'a' 'c'
//    |\ follow: 'g' 'h'
//    |
// At the end, unreachable symbols and unused production rules are listed.

// global variables private to this package
var (
	barcol  int // static // column number for vertical bar
	contcol int // static // column number for continuation
)

// printing utility

func WriteGrammar() {
	writeg()
}

// static void outprod( PPRODUCTION p )
// put out RHS of rule p
func outprod(p PPRODUCTION) {
	var e PELEMENT
	var s PSYMBOL

	if p != nil { // empty rules should never happen, but be safe
		for e = p.data; e != nil; e = e.next {
			s = e.data
			outspacesym(s, contcol, ' ')
			outsymbol(s)
		}
	}
}

// static void outprodgroup( PSYMBOL s )
// put out all rules with s on LHS
func outprodgroup(s PSYMBOL) {
	var p PPRODUCTION
	var e PELEMENT
	var ss PSYMBOL

	outline()
	outsymbol(s)
	outchar(' ')
	barcol = getoutcol() + 1 // remember indent for next rule
	outstr(RULESYM)
	contcol = getoutcol() + 1 // remember indent for continuation

	// output first production on same line
	p = s.data
	outprod(p)

	for p.next != nil { // output successive productions
		p = p.next

		outline()
		outspaces(barcol)
		outstr("| ")

		outprod(p)
	}

	if s.starter != nil { // output start set
		outline()
		outchar(COMMENT)
		outstr(" start set:  ")
		contcol = getoutcol() + 1 // remember indent

		for e = s.starter; e != nil; e = e.next {
			ss = e.data
			outspacesym(ss, contcol, COMMENT)
			outsymbol(ss)
		}
	}
	if s.follows != nil { // output follow set
		outline()
		outchar(COMMENT)
		outstr(" follow set: ")
		contcol = getoutcol() + 1 // remember indent

		for e = s.follows; e != nil; e = e.next {
			ss = e.data
			outspacesym(ss, contcol, COMMENT)
			outsymbol(ss)
		}
	}
	if s.follows != nil || s.starter != nil {
		// output blank line to separate from the next rule
		outline()
	}
}

// static void outreachable( PSYMBOL s )
// recursively output reachable production
func outreachable(s PSYMBOL) {
	// handles used in list traversals
	var p PPRODUCTION
	var e PELEMENT
	var ss PSYMBOL

	// only output nonterminals with their rules
	if s.data != nil {
		outprodgroup(s)
	}

	// mark that it is printed
	s.state = TOUCHED

	// recursive walk through all of its productions
	p = s.data
	for p != nil {
		// walk through all elements of each production
		e = p.data
		for e != nil {
			// touch the associated symbol
			ss = e.data
			if ss.state == UNTOUCHED {
				outreachable(ss)
			}
			// move to next element
			e = e.next
		}
		// move to next rule
		p = p.next
	}
}

// void writeg()
// write grammar structure documented in grammar.h
func writeg() {
	var s PSYMBOL
	var header bool

	outsetup()

	if head != nil { // there is a distinguished symbol
		outstr("> ")
		outsymbol(head)
	} else {
		outchar(COMMENT)
		outstring(tocstring(" no distinguished symbol!"))
	}

	if emptypt != nil { // there is an empty symbol
		outline()
		outstr("/ ")
		outsymbol(emptypt)
	}

	for s = symlist; s != nil; s = s.next {
		s.state = UNTOUCHED
	}
	if head != nil {
		outline()
		outreachable(head)
	}

	header = false
	for s = symlist; s != nil; s = s.next {
		if (s.data == nil) && (s.state == TOUCHED) {
			if !header {
				outline()
				outline()
				outchar(COMMENT)
				outstr(" terminals:  ")
				contcol = getoutcol() + 1 // remember indent
				header = true
			}
			outspacesym(s, contcol, COMMENT)
			outsymbol(s)
		}
	}

	header = false
	for s = symlist; s != nil; s = s.next {
		if (s.data != nil) && (s.state == UNTOUCHED) {
			if !header {
				outline()
				outline()
				outchar(COMMENT)
				outstr(" unused productions")
				header = true
			}
			outprodgroup(s)
		}
	}

	header = false
	for s = symlist; s != nil; s = s.next {
		if (s.data == nil) && (s.state == UNTOUCHED) {
			if !header {
				outline()
				outline()
				outchar(COMMENT)
				outstr(" unused terminals: ")
				contcol = getoutcol() + 1 // remember indent
				header = true
			}
			outspacesym(s, contcol, COMMENT)
			outsymbol(s)
		}
	}
	outline()
}
