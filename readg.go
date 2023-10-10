// gtools - a collection of grammar manipulation tools
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package gtools

func ReadGrammar() {
	readg()
}

// global variables private to this package
var (
	ch   byte // static // most recent char read from stdin (could be EOF)
	line int  // static // current line number on stdin, used in error reports
)

// parsing utility
var (
	endlist bool // set by nonblank at end of list, reset when understood
	endrule bool // set by nonblank at end of rule, reset when understood
)

// PSYMBOL definesym( char * str )
// define str in the main symbol list, it must not be already there
func definesym(str []byte) PSYMBOL {
	var s PSYMBOL
	var i int   // character index in str
	var len int // length of string str

	// create the symbol itself and link it into place
	s = NEWSYMBOL()
	*symlistend = s
	s.line = line
	s.name = stringlim
	s.data = nil
	s.state = UNTOUCHED
	s.starter = nil
	s.follows = nil
	s.next = nil
	symlistend = &(s.next)

	// copy the characters of the symbol name into place
	len = int(str[0]) // works because length is encoded as first byte
	for i = 0; i <= len; i++ {
		if stringlim >= STRINGLIMIT {
			errormsg("STRING POOL OVERVFLOW", line)
		} else {
			stringtab[stringlim] = str[i]
			stringlim = stringlim + 1
		}
	}
	return s
}

// void errormsg( char * msg, int line ) {
// output an error message msg attributed to the given source line
// use -1 as a line number if line attribution does not work
func errormsg(msg string, line int) {
	fputs(" >>", stderr)
	fputs(msg, stderr)
	if line > 0 {
		fprintf(stderr, " on line %d", line)
	}
	fputs("<<\n", stderr)
}

// static void extendsym( int * len, char * str, char ch )
func extendsym(len *int, str []byte, ch byte) {
	if *len >= SYMLEN {
		errormsg("SYMBOL TOO LONG", line)
	} else {
		*len = *len + 1
		str[*len] = ch
	}
}

// static PPRODUCTION getprod()
// get a list of production rules
func getprod() PPRODUCTION {
	var ph PPRODUCTION // the head of the production list
	var p PPRODUCTION  // the current production
	var np PPRODUCTION // the new production

	ph = nil
	p = nil

	// do {...} while (!endrule);
	for firstTime := true; firstTime || !endrule; firstTime = false {
		np = NEWPRODUCTION()
		np.line = line
		nonblank()
		if !endlist { // the normal case
			np.data = getsymlist()
		} else { // nothing there
			errormsg("EMPTY PRODUCTION RULE", np.line)
			if emptypt != nil {
				np.data = NEWELEMENT()
				np.data.line = line
				np.data.next = nil
				np.data.data = emptypt
			} else {
				np.data = nil
			}
			endlist = false
		}
		// link it in place
		if ph == nil {
			ph = np
		}
		if p != nil {
			p.next = np
		}
		p = np
	}
	endrule = false
	endlist = false
	p.next = nil
	return ph
}

// static PSYMBOL getsymbol()
// get symbol from input to str
func getsymbol() PSYMBOL {
	var str [SYMLEN + 1]byte // most recent symbol from stdin
	// string length is encoded in str[0]

	var len int // index of last used space in str

	// Must be called with ch nonblank, first char of symbol
	len = 1
	str[len] = ch

	if ch == '<' { // may be a < quoted symbol
		ch = getchar()
		if ((ch <= 'z') && (ch >= 'a')) || ((ch <= 'Z') && (ch >= 'A')) || ((ch <= '9') && (ch >= '0')) { // definitely < quoted
			for { // consume bracketed symbol
				extendsym(&len, str[:], ch)
				if ch == '>' {
					break
				}
				ch = getchar()
				if ch == '\n' {
					break
				}
				if ch == EOF {
					break
				}
			}
			if ch == '>' { // normal end of symbol
				ch = getchar() // skip trailing >
			} else { // abnormal end of symbol
				errormsg("MISSING CLOSING > MARK", line)

				// fake it
				extendsym(&len, str[:], '>')
			}
		} else { // symbol ends at next blank (broadly speaking)
			for { // symbol
				extendsym(&len, str[:], ch)
				ch = getchar()
				if ch == ' ' {
					break
				}
				if ch == '\t' {
					break
				}
				if ch == '\n' {
					break
				}
				if ch == EOF {
					break
				}
			}
		}
	} else if (ch == '"') || (ch == '\'') { // quoted
		ch = getchar()
		for (ch != str[1]) && (ch != '\n') && (ch != EOF) {
			extendsym(&len, str[:], ch)
			ch = getchar()
		}
		if ch == str[1] {
			extendsym(&len, str[:], ch)
			ch = getchar()
		} else {
			errormsg("MISSING CLOSING QUOTE", line)

			// fake it
			extendsym(&len, str[:], str[1])
		}
	} else { // symbol did not begin with < or quote, ends with space
		ch = getchar()
		for (ch != ' ') && (ch != '\t') && (ch != '\n') && (ch != EOF) {
			extendsym(&len, str[:], ch)
			ch = getchar()
		}
	}

	// we now have a symbol in str[1...len] !
	if !(0 <= len && len <= 255) {
		panic("assert(0 <= len <= 255")
	}
	str[0] = byte(len) // record symbol length
	return lookupordefine(str[:])
}

// static PELEMENT getsymlist()
// get the list of symbols on RHS of rule
func getsymlist() PELEMENT {
	var s PELEMENT

	nonblank()
	if endlist {
		endlist = false
		return nil
	} else {
		s = NEWELEMENT()
		s.line = line
		s.data = getsymbol()
		s.next = getsymlist()
		return s
	}
}

// static PSYMBOL lookupordefine( char * str )
// lookup str in the main symbol list, and add it if required
func lookupordefine(str []byte) PSYMBOL {
	// var ps *PSYMBOL = &symlist // reference to current symbol
	var s PSYMBOL // current symbol
	// var j, k int               // character indices
	// var m PSYMBOL
	s = lookupsym(str)
	if s != nil {
		return s
	}
	return definesym(str)
}

// PSYMBOL lookupsym( char * str )
// lookup str in the main symbol list, return NULL if not found
// str[0] is length of symbol, in characters
func lookupsym(str []byte) PSYMBOL {
	var s PSYMBOL
	var i int      // character index in str
	var j STRINGPT // character index in stringtab
	var len int    // length of str

	len = int(str[0]) // works because length is encoded as first byte
	for s = symlist; s != nil; s = s.next {
		i = 0
		j = s.name
		// bug: panic: runtime error: index out of range [51] with length 51
		// mdh: limit i to 0...len
		for i <= len && str[i] == stringtab[j] {
			i = i + 1
			j = j + 1
		}
		if i > len { // we found it
			return s
		}
		// the above works because length is encoded as first byte
	}
	// we only get here if we don't find a symbol
	return nil
}

// static void newline()
// advance to next line, called when ch == '\n'
func newline() {
	line = line + 1
	ch = getchar()
}

// static void nonblank()
// fancy scan for a nonblank character in ch
func nonblank() {
	for ch == '|' || ch == ' ' || ch == '\t' || ch == '\n' || ch == EOF {
		if ch == '|' {
			endlist = true
			ch = getchar()
			return
		} else if ch == EOF {
			endrule = true
			endlist = true
			return
		} else if ch == '\n' {
			newline()
			if ch != ' ' && ch != '\t' { // line starts with nonblank
				endrule = true
				endlist = true
				return
			}
		} else { // must have been blank or tab
			ch = getchar()
		}
	}
	return
}

// static void skipline()
// skip the rest of this line
func skipline() {
	for ch != '\n' && ch != EOF {
		ch = getchar()
	}
	if ch == '\n' {
		newline()
	}
}

// static void skipwhite()
// simple scan for a nonblank character in ch
func skipwhite() {
	for ch == '\t' || ch == ' ' {
		ch = getchar()
	}
}

// readg: read grammar into global grammar structure in grammar.h
func readg() {
	var s PSYMBOL
	var p PPRODUCTION
	var ok bool

	// global initialization
	stringlim = 0 // no characters have been put in stringtab
	symlist = nil // no symbols have been encountered
	symlistend = &symlist
	head = nil    // we have no distinguished symbol
	emptypt = nil // we have no empty symbol

	// prime the input stream
	line = 1
	ch = getchar()
	endlist = false
	endrule = false

	// while (ch != EOF) {
	for ch != EOF {
		// while (ch == '\n') newline();
		for ch == '\n' {
			newline()
		}
		if ch == '>' { // Identify distinguished symbol
			if head != nil {
				errormsg("EXTRA DISTINGUISHED SYMBOL", line)
			} else {
				ch = getchar() /* skip > */
				skipwhite()
				if (ch == '\n') || (ch == EOF) {
					errormsg("NO DISTINGUISHED SYMBOL", line)
				} else {
					head = getsymbol()
				}
			}
			skipline()
		} else if ch == '/' { // Identify the empty (/)symbol
			if emptypt != nil {
				errormsg("EXTRA EMPTY SYMBOL", line)
				skipline()
			} else {
				ch = getchar() /* skip */
				skipwhite()
				if ch == '\n' || ch == EOF {
					errormsg("NO EMPTY SYMBOL", line)
				} else {
					emptypt = getsymbol()
				}
			}
			skipline()
		} else if ch == COMMENT { // COMMENT
			skipline()
		} else if ch != EOF { // WE MIGHT HAVE A RULE
			s = getsymbol()
			skipwhite()

			ok = false
			if ch == ':' { // consume ::= or := or : or =
				ok = true
				ch = getchar()
				if ch == ':' {
					ch = getchar()
					if ch == '=' {
						ch = getchar()
					}
				} else if ch == '=' {
					ch = getchar()
				}
			} else if ch == '=' {
				ok = true
				ch = getchar()
			}

			if ok { // WE HAVE A RULE s ::= rule
				p = s.data
				if p == nil {
					s.data = getprod()
				} else {
					for p.next != nil {
						p = p.next
					}
					p.next = getprod()
				}
			} else { // NOT A RULE, JUST s ...comment
				errormsg("MISSING ::= OR EQUIVALENT", line)
				skipline()
			}
		}
	}

	if head == nil {
		errormsg("DISTINGUISHED SYMBOL NOT GIVEN", -1)
	} else if TERMINAL(head) {
		errormsg("DISTINGUISHED SYMBOL IS TERMINAL", head.line)
	}
	if (emptypt != nil) && (NONTERMINAL(emptypt)) {
		errormsg("EMPTY SYMBOL IS NONTERMINAL", emptypt.data.line)
	}
	line = -1 // mark any new line numbers as fictional
}
