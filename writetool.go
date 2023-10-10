// gtools - a collection of grammar manipulation tools
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package gtools

// global variables private to this package
var (
	column int // static // last column number filled on line
)

// int getoutcol()
// note what column we are on
func getoutcol() int {
	return column
}

// void outchar( char ch )
// put out one character
func outchar(ch byte) {
	putchar(ch)
	column++
}

// void outline()
// put out newline
func outline() {
	putchar('\n')
	column = 0
}

// void outsetup()
// setup for package use
func outsetup() {
	column = 0
}

// void outspaces( int c )
// put out spaces until column = c
func outspaces(c int) {
	for column < c {
		outchar(' ')
	}
}

// void outspacesym( PSYMBOL s, int c, char ch )
// put out a space, or if s won't fit, return to column c starting the line with ch
func outspacesym(s PSYMBOL, c int, ch byte) {
	var len int      // length of s, in chars
	var pos STRINGPT // position of s in stringtab

	pos = s.name
	len = int(stringtab[pos]) // works because length is encoded as first byte

	// does s fit on the line?
	if (column + 1 + len) > 80 { // no, move to next line
		outline()
		if c > 1 {
			outchar(ch)
		}
		outspaces(c)
	} else { // yes, output just one space
		outchar(' ')
	}
}

// copy of outstring for Go strings
func outstr(s string) {
	for _, r := range s {
		if !(0 <= r && r <= 255) {
			panic("assert(0 <= r <= 255")
		}
		outchar(byte(r))
	}
}

// void outstring( char * p )
// put out one null-terminated string
func outstring(p []byte) {
	for _, ch := range p {
		if ch == 0 {
			break
		}
		outchar(ch)
	}
}

// void outsymbol( PSYMBOL s )
// put symbol to output
func outsymbol(s PSYMBOL) {
	var len int      // length of s, in chars
	var pos STRINGPT // positiion of s in stringtab

	pos = s.name
	len = int(stringtab[pos]) // works because length is encoded as first byte

	for len != 0 {
		len = len - 1
		pos = pos + 1
		outchar(stringtab[pos])
	}
}
