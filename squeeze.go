// gtools - a collection of grammar manipulation tools
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package gtools

func Squeeze() {
	squeeze()
}

// Global
var (
	// static variables for squeeze file
	_squeeze struct {
		change bool // record that a change was made to the grammar
	}
)

// Worker routines

// static void squeezesymbols()
// squeeze out redundant symbols
func squeezesymbols() {
	/* handles used in list traversals */
	var s PSYMBOL
	var p PPRODUCTION
	var e PELEMENT
	var s1 PSYMBOL
	var p1 PPRODUCTION
	var e1, e2 PELEMENT

	/* for all symbols */
	for s = symlist; s != nil; s = s.next {

		/* for all production rules of that symbol */
		for p = s.data; p != nil; p = p.next {

			/* for all elements of that production rule */
			for e = p.data; e != nil; e = e.next {

				s1 = e.data
				p1 = s1.data
				if (p1 != nil) && (p1.next == nil) {
					/* symbol s1 has just 1 rule p1
					   substitute that rule for s1 in p */

					e1 = p1.data
					if e1 != nil { /* if well formed */
						/* first element of rule
						   overwrites first element */
						e.data = e1.data
						_squeeze.change = true

						/* now copy rest of rule */
						e1 = e1.next
						for e1 != nil {
							e2 = NEWELEMENT()
							e2.data = e1.data
							e2.next = e.next
							e.next = e2
							e = e2
							e1 = e1.next
						}
					}
				}
			}
		}
	}
}

// static void squeezerules()
// squeeze out redundant production rules
//
//	/* handles used in list traversals */
func squeezerules() {
	var s PSYMBOL
	var p PPRODUCTION
	var qp *PPRODUCTION /* pointer to q */
	var q PPRODUCTION

	/* for all symbols */
	for s = symlist; s != nil; s = s.next {

		/* for all productions */
		for p = s.data; p != nil; p = p.next {

			/* for all additional productions of symbol s */
			qp = &(p.next)
			q = *qp /* that is, q = p.next */
			for q != nil {
				if samerule(p, q) {
					/* rule q is redundant, eliminate it */
					*qp = q.next
					_squeeze.change = true
				} else {
					/* move to next production */
					qp = &(q.next)
				}
				q = *qp
			}
		}
	}
}

// The interface

// void squeeze()
// squeeze out redundant rules and symbols */
func squeeze() {
	/* count the symbols and setup for reachability analysis */
	// do {...} while _squeeze.change
	for firstTime := true; firstTime || _squeeze.change; firstTime = false {
		_squeeze.change = false
		squeezerules()
		squeezesymbols()
	}
}
