// gtools - a collection of grammar manipulation tools
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package gtools

// convert Wirth style EBNF to pure BNF
//   <a> ::= <b> [ <c> ] <d> { <e> } <f>
// is replaced with
//   <a> ::= <b> <q01q> <d> <q02q> <f>
//   <q01q> ::= <c>
//           |  <empty>
//   <q02q> ::= <e>
//           |  <e> <q02q>
//           |  <empty>
// This is a kluge!  The metasymbols ( ) [ ] { } are parsed as terminals
// in the grammar by the readgram routine.

func GDeEBNF() {
	gdeebnf()
}

// Global
var (
	lparen  PSYMBOL // static // '('
	rparen  PSYMBOL // static // ')'
	lsquare PSYMBOL // static // '['
	rsquare PSYMBOL // static // ']'
	lcurly  PSYMBOL // static // '{'
	rcurly  PSYMBOL // static // '}'
)

// Worker routines

// static void addemptyrule( PSYMBOL s ) {
// add an empty production to nonterminal s
func addemptyrule(s PSYMBOL) {
	var p PPRODUCTION // a new production | <empty>
	var e PELEMENT    // a new element <empty>

	if emptypt == nil {
		errormsg("EMPTY SYMBOL MUST BE DEFINED", s.data.line)
		return // quit if can't add empty rule
	}

	e = NEWELEMENT()
	e.line = s.data.line // take the line number from its first rule
	e.next = nil
	e.data = emptypt

	p = NEWPRODUCTION()
	p.line = s.data.line // take the line number as above
	p.next = s.data
	p.data = e
	p.state = UNTOUCHED
	s.data = p
}

// static PSYMBOL extractuntil( PSYMBOL s, PPRODUCTION p, PELEMENT e, PSYMBOL rsym, int * nsc )
// delete from e.next of production rule p and successors
// until an element that references rsym; element e is edited
// to refer to a new symbol with a name derived from the name of s
// using *nsc and a rule holding the deleted parts, minus the
// parentheses.  Returns a pointer to the new symbol, never nil. */
func extractuntil(s PSYMBOL, p PPRODUCTION, e PELEMENT, rsym PSYMBOL, nsc *int) PSYMBOL {
	var ee PELEMENT    // the element holding the end paren, if we find it
	var eep *PELEMENT  // points to ee, once found, so we can delete it
	var ns PSYMBOL     // a new symbol
	var nr PPRODUCTION // the tail of the rule list hanging from ns
	var ne PELEMENT    // first element within parens starts body of rule
	var lsym PSYMBOL   // the left brace that balances rsym
	var nest int       // bracket nesting level

	lsym = e.data // we were called with left brace current element

	// make new symbol to refer to a new set of rules
	ns = inventsymbol(s, nsc)
	ns.line = e.line

	// make the first new rule that will hang under ns
	nr = NEWPRODUCTION()
	nr.line = e.line
	nr.state = UNTOUCHED
	nr.next = nil
	nr.data = ne

	// set up preliminary linkage
	ne = e.next  // the new rule begins after the opening bracket
	nr.data = ne // the new rule begins after the opening bracket
	ns.data = nr // hang the new rule in place
	e.data = ns  // replace open bracket with new symbol

	// look for end paren of set of rules -- eep will point at it
	nest = 0
	eep = &(nr.data)
	ee = *eep
	for {
		if ee == nil {
			// we hit the end of a rule, either because of missing end bracket or bracketed alternatives
			if nr.data == nil { // previous rule was empty!
				errormsg("EMPTY BRACKETED RULE", nr.line)
				// add empty element to rule if possible
				if emptypt != nil {
					ne = NEWELEMENT()
					nr.data = ne
					ne.line = nr.line
					ne.next = nil
					ne.data = emptypt
					ne = ee
				}
			}

			if p.next == nil { // quit loop
				break
			}

			// we have hit ( | and there is another production,
			// so swipe a rule from p's rule list, hang it under ns
			nr.next = p.next
			nr = nr.next
			p.next = nr.next
			nr.next = nil

			// move to next rule looking for end paren
			eep = &(nr.data)
			ee = *eep
			ne = ee

		} else { // ee != nil
			// scan for end paren, accounting for nesting
			if ee.data == rsym {
				if nest == 0 { // quit loop
					break
				}
				nest = nest - 1
			} else if ee.data == lsym {
				nest = nest + 1
			}

			// march down this rule looking for end paren
			eep = &(ee.next)
			ee = *eep
		}
	}

	if ee != nil { // we found the balancing parenthesis
		*eep = nil       // body of list does not include rsym
		e.next = ee.next // snip bracketed body out of rule

		if nr.data == nil { // final rule of set was empty!
			errormsg("EMPTY BRACKETED RULE", ee.line)
			// add empty element to rule if possible
			if emptypt != nil {
				nr.data = NEWELEMENT()
				nr.data.line = nr.line
				nr.data.next = nil
				nr.data.data = emptypt
			}
		}
	} else { // failure, unbalanced parens, treat as end paren
		// assert ee == *eep == nil
		e.next = nil /* snip body out of rule */

		if rsym == rparen {
			errormsg("MISSING )", nr.line)
		} else if rsym == rsquare {
			errormsg("MISSING ]", nr.line)
		} else /* rsym == rcurly */ {
			errormsg("MISSING }", nr.line)
		}
		// assert complaint about ne == nil was already done
	}
	return ns
}

// static PSYMBOL getdelsym( char * str )
// get symbol with name str and delete it from the master symbol list;
// str must begin with the length of the string, in characters;
// str must be nil terminated so error messages work;
// the symbol str must terminal; from now on, it is a metasymbol
func getdelsym(str []byte) PSYMBOL {
	var s PSYMBOL    // the symbol we are looking up
	var pss *PSYMBOL // tools for deletion
	var ss PSYMBOL   // tools for deletion

	s = lookupsym(str)

	if s == nil { // nothing to do here
		return nil
	}

	if s.data != nil { // it's nonterminal?
		errormsg("BRACE SHOULD BE NONTERMINAL", s.data.line)
	}

	// now find pss, pointer to s in symlist -- we know it's there
	pss = &symlist
	ss = *pss
	for ss != s {
		// walk onward in symbol list
		pss = &(ss.next)
		ss = *pss
	}

	// unlink s from symlist
	if s.next == nil { // s was final element of list
		symlistend = pss
	}
	*pss = s.next

	// done
	return s
}

// static PSYMBOL inventsymbol( PSYMBOL s, int * nsc ) {
// invent and return a new symbol; as initialized, it is a terminal
// symbol a unique name derived from the name of s and *nsc.
// Add rules to it and it will become nonterminal.
func inventsymbol(s PSYMBOL, nsc *int) PSYMBOL {
	var newname [SYMLEN]byte
	var i int      // index into newname
	var j STRINGPT // stringpool index
	var len int    // length of symbol s
	var dash int   // location of dash in newname
	var ext int    // name extension
	var quote byte // quotation mark at end of name

	// first, copy base name of new symbol into place
	j = s.name
	len = int(stringtab[j]) // works because length is encoded as first byte
	for i = 0; i <= len; i++ {
		newname[i] = stringtab[j]
		j++
	}

	// figure out what kind of quotes are in use, if any
	if (newname[1] == '<' && newname[len] == '>') || (newname[1] == '"' && newname[len] == '"') || (newname[1] == '\'' && newname[len] == '\'') {
		quote = newname[len]
		dash = len
		/* extension will wipe out trailing quote */
	} else {
		quote = ' '
		dash = len + 1
		// extension will be appended
	}

	newname[dash] = '-'

	// now, try to find a name extension that is not already in use
	// do {...} while (lookupsym( newname ) != nil)
	for firstTime := true; firstTime || lookupsym(newname[:]) != nil; firstTime = false {
		ext = *nsc
		*nsc = ext + 1

		// create name extension of the form -a or -b
		i = dash
		// do {...} while (ext > 0)
		for firstTime := true; firstTime || ext > 0; firstTime = false {
			i = i + 1
			newname[i] = byte(ext%26) + 'a'
			ext = ext / 26
		}

		// put back the quote that extension overwrote
		if quote != ' ' {
			i = i + 1
			newname[i] = quote
		}

		// set the length to reflect extension etc
		if !(0 <= i && i <= SYMLEN) {
			panic("assert(0 <= i <= SYMLEN)")
		}
		newname[0] = byte(i)

		// keep trying until new name is genuinely new
	}

	// newname is genuinely new
	s = definesym(newname[:])
	return s
}

// static void makeiterative( PSYMBOL s )
// make nonterminal s iterate
func makeiterative(s PSYMBOL) {
	var p PPRODUCTION // a production under s
	var e PELEMENT    // an element of p
	var pe *PELEMENT  // the pointer to e
	var line int      // best guess at source line number for empty element

	// add a self reference to end of each rule of s
	for p = s.data; p != nil; p = p.next { // for each production
		// find the end of the production
		pe = &(p.data)
		line = p.line
		e = *pe
		for e != nil {
			pe = &(e.next)
			line = e.line
			e = *pe
		}

		// tack on a new element at the end that references s
		e = NEWELEMENT()
		e.line = line
		e.next = nil
		e.data = s
		*pe = e
	}

	// add an empty production to s to terminate iteration
	addemptyrule(s)
}

// static void process( PSYMBOL s )
// process one symbol
func process(s PSYMBOL) {
	var p PPRODUCTION   // this production rule
	var pp *PPRODUCTION // the pointer to p so we can delete rules
	var e PELEMENT      // this element of p
	var ep *PELEMENT    // the pointer to e so we can delete elements
	var ns PSYMBOL      // a new symbol
	var nsc int         // count used to uniquely name new symbols

	nsc = 0 // any added symbols start with -a if possible

	pp = &(s.data)
	p = *pp
	for p != nil {
		// for each production rule of symbol, allowing for deletion
		ep = &(p.data)
		e = *ep
		for e != nil {
			// for each element of rule e
			if e.data == rparen { // syntax error
				errormsg("UNEXPECTED )", e.line)
				e = e.next
				*ep = e // clip it from rule

			} else if e.data == rsquare { // syntax error
				errormsg("UNEXPECTED ]", e.line)
				e = e.next
				*ep = e /* clip it from rule */

			} else if e.data == rcurly { // syntax error
				errormsg("UNEXPECTED }", e.line)
				e = e.next
				*ep = e /* clip it from rule */

			} else if e.data == lparen {
				ns = extractuntil(s, p, e, rparen, &nsc)
				// above may delete rules following p
				// no extra work to be done

				process(ns)

				// move onward down the list
				ep = &(e.next)
				e = *ep

			} else if e.data == lsquare {
				ns = extractuntil(s, p, e, rsquare, &nsc)
				// above may delete rules following p

				process(ns)

				// make new symbol possibly empty
				addemptyrule(ns)

				// move onward down the list
				ep = &(e.next)
				e = *ep

			} else if e.data == lcurly {
				ns = extractuntil(s, p, e, rcurly, &nsc)
				// above may delete rules following p

				process(ns)

				// make new symbol possibly iterate
				makeiterative(ns)

				// move onward down the list
				ep = &(e.next)
				e = *ep

			} else {
				// move onward down the list
				ep = &(e.next)
				e = *ep
			}
		}
		pp = &(p.next)
		p = *pp
	}
	s.state = TOUCHED
}

// The interface

// void gdeebnf()
// remove Wirth-style EBNF features
func gdeebnf() {
	var s PSYMBOL

	// initializations (remember that the first byte holds the length!)
	lparen = getdelsym([]byte{1, '('})
	rparen = getdelsym([]byte{1, ')'})
	lsquare = getdelsym([]byte{1, '['})
	rsquare = getdelsym([]byte{1, ']'})
	lcurly = getdelsym([]byte{1, '{'})
	rcurly = getdelsym([]byte{1, '}'})

	// prevent duplicate processing
	for s = symlist; s != nil; s = s.next {
		s.state = UNTOUCHED
	}

	// now process all the symbols
	s = symlist
	for s != nil {
		// for all symbols s
		if s.state == UNTOUCHED {
			process(s)
		}
		// process may have added to tail of symlist
		s = s.next
	}
}
