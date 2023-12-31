# History
Gtools is a set of grammar tools,
a set of tools for manipulating BNF grammars,
created by Douglas W. Jones.

The remainder of this is copied from his README file,
lightly edited to be markdown-ish.

## The tools

*  gcopy 	-- copy (and pretty print) a grammar, stdin to stdout
*  gdeebnf	-- remove Wirth-style EBNF extensions
*  gdeempty	-- rewrite grammar to remove empty symbol, if possible
*  gsqueeze	-- like gramcopy but also eliminate redundant rules and symbols
*  gsample	-- output one random string from the grammar, stdin to stdout
*  gstartfollow 	-- output the start set and follow set of all symbols
*  gstats	-- gather basic statistics on grammar, stdin to stdout
*  Makefile	-- make one or more of the above, or make with no args for help

## Source files for main programs that weave together components

*  gcopy.c	-- gcopy
*  gdeebnf.c	-- gdeebnf
*  gdeempty.c	-- gdeempty
*  gsqueeze.c	-- gsqueeze
*  gsample.c	-- gsample
*  gstartfollow 	-- gstartfollow
*  gstats.c	-- gstats

## Source files for components

*  grammar.h	-- framework
*  readgram.c	-- read grammar from stdin
*  readgram.h
*  reachable.c	-- mark reachable symbols in the grammar
*  reachable.h
*  writetool.c	-- tools for formatting output to stdout
*  writetool.h
*  writegram.c	-- pretty print grammar to stdout
*  writegram.h
*  gramstats.c	-- statistics gathering mechanism
*  gramstats.h
*  squeezegram.c	-- eliminate redundant production rules and nonterminals
*  squeezegram.h
*  samplegram.c	-- generate a random string from the grammar
*  samplegram.h
*  deempty.c	-- rewrite grammar to remove empty symbol, if possible
*  deempty.h
*  deebnf.c	-- remove Wirth-style EBNF extensions
*  deebnf.h

## Some example grammars

*  gram.gr	-- the EBNF grammar of the BNF and EBNF notation used here
*  errors.gr	-- generate as many error messages as possible
*  bnf.gr	-- simple BNF for expressions
*  ebnf.gr	-- extended BNF for expressions

## Warranty

The gtools code is basically junk, so it is offered here as freeware.
You get what you pay for,
and if it does anything you consider useful, that is good,
but if it breaks,
you bear the entire burden of responsibility for using such junk.

## License
The gtools code is basically junk, so there's no point to dolling it up with fancy license agreements.
Basically, it's free;
if you can make something useful from it,
you are welcome to assert any rights you wish to your improvements,
so long as you credit me (Douglas W. Jones at the University of Iowa) for the foundation on which you built.
Feel free to redistribute this code to anyone,
as is or with your additions and improvements.
