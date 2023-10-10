# gtools

GTools is a port of grammar tools created by Douglas W. Jones,
a professor at the University of Iowa.

# Copying
This package is a port of
[Douglas W. Jones](https://homepage.cs.uiowa.edu/~dwjones/)
[gtools](https://homepage.cs.uiowa.edu/~dwjones/compiler/gtools/).

Copyright in the original source code and documentation is owned by Professor Jones.
The file `gtools.txt` in the root of this package contains the unmodified source code and documentation.

Most of the content of this document has been sourced from the `gtools` web page.

# Building

# Makefile

# Documentation
Documentation from the web site.

## Introduction
## Installing gtools
## Notation
## The Tools
### gcopy — copy (and pretty print) a *BNF* grammar
### deebnf — convert an *EBNF* grammar to *BNF*
### gdeempty — eliminate empty symbols from a *BNF* grammar
### gsqueeze — eliminate redundancy from a *BNF* grammar
### gsample — generate an example string from a *BNF* grammar
### gstartfollow — find the start and follow sets of all non-terminals
### gstats — count the rules and symbols of a *BNF* grammar
## Notes

## Introduction
While there may be an officially designated reference grammar for any particular language,
it is important to remember that there are an infinite number of grammars for any particular language.
One grammar may be better suited to pedagogy,
arranged so that each memorable component of the language is named,
while another grammar may be closely matched to some particular approach to parsing.

A grammar intended for pedagogic use may well contain numerous non-terminal symbols that are primarily intended as anchors for discussions of features of the language but that play no useful role in parsing the language.
Another grammar might be constructed under the constraints of a particular parsing method,
with rules that fit the constraints of the parsing algorithm but serve no useful pedagogical purpose.
Manual transformation between such grammars is always possible,
but it is an error-prone clerical task.

The `gtools` package contains a number of utilities that transform and compute properties of grammars written in *BNF* and *EBNF*.
These tools all share both a common internal data structure and a common syntax for textual expression of grammars.
The syntactic notation is basically *BNF* (Backus–Naur Form) and support is included for the most common extensions found in *EBNF* (Extended Backus-Naur Form),
coming close to supporting Wirth syntax notation although not compatible with *ISO EBNF*.

## Installing gtools
TBD.

## Notation
The file `gram.gr` included in the repository gives a grammar for the *BNF* and *EBNF* notations used by `gtools`.
Note that there is nothing special about the `.gr` filename suffix.
It is used here to tag grammar files as such, but filename suffixes have no special meaning in the `gtools` package.

Here are two small example grammars in the `gtools` notation, also included in the distribution:

    # bnf.gr -- an example classical *BNF* grammar for expressions
    > <expression>
    <expression> ::= <term> | <expression> + <term> | <expression> - <term>
    <term> ::= <factor> | <term> * <factor> | <term> / <factor>
    <factor> ::= <element> | - <element>
    <element> ::= <number> | <identifier> | ( <expression> )

    # ebnf.gr -- an example extended *BNF* grammar for expressions
    > expression
    / ''
    expression = term { ( '+' | '-' ) term }
    term = factor { ( '*' | '/' ) factor }
    factor = [ '-' ] ( number | identifier | '(' expression ')' )

They describe the same language using variant notations,
and will be used as the basis for demonstrations of the various tools in the following sections.

In classical *BNF*, angle brackets surround names of non-terminal symbols,
and quotation marks are not needed around terminal symbols,
except when the vertical bar is used as a terminal symbol.
For practical purposes, grammars used with `gtools` must explicitly declare the distinguished or head symbol of the grammar,
and comments may be included starting with a pound or hash sign.

The `gtools` notation considers any symbol on the left-hand side of a rule to be a non-terminal symbol,
and any symbol not found on the left-hand side of a rule to be terminal.
So, the angle brackets can be omitted from non-terminal symbols and terminal symbols can be quoted,
as in Wirth's (and many other variant *EBNF*) notations.
Furthermore, the classical `::=` of *BNF* may be abbreviated `:` or `:=` or `=`.

When the `gtools` *EBNF* notation is used, parentheses may be used to surround alternatives,
square brackets may be used to surround optional elements,
and curly braces may be used to surround elements that may be iterated zero or more times.
The `gtools` *EBNF* processor `deebnf` requires that the user name an empty symbol,
even if it does not occur in the grammar.
The empty string `''` was used for this in the above.

If, as in Wirth's notation, periods are to be used as terminators on rules,
the period must be offset from the rule by a blank and the period should be defined in an added rule as an empty symbol.
For example, part of the above file could have been written as follows:

    > expression
    / ''
    . ::= ''
    expression = term { ( '+' | '-' ) term } .

After applying the `deebnf` and `deempty` filters from `gtools`, the result will be pure *BNF* with none of Wirth's notation surviving
(although the form of the terminal and non-terminal symbols will be unchanged; no angle brackets will be added or quotation marks removed).

## The Tools
### gcopy — copy (and pretty print) a *BNF* grammar
Given the example *BNF* grammar given above, stored in a file named bnf.gr, type this shell command while in the `gtools` directory:

```bash
./gcopy < bnf.gr
```

This will produce the following output:

    > <expression>
    
    <expression> ::= <term>
                  |  <expression> + <term>
                  |  <expression> - <term>
    <term> ::= <factor>
            |  <term> * <factor>
            |  <term> / <factor>
    <factor> ::= <element>
              |  - <element>
    <element> ::= <number>
               |  <identifier>
               |  ( <expression> )
    
    # terminals:   + - * / <number> <identifier> ( )

Any comments appearing in the original are stripped out,
and the rules of the grammar are presented in an order determined by a depth-first traversal of the grammar starting from the distinguished symbol.
Each alternative is presented on its own line (or lines, if it is very long),
indented to the same depth as all other alternatives under the same non-terminal.
Finally, comments are appended documenting the terminal symbols of the grammar and any production rules and symbols that were not encountered during the traversal.

The output is fully compatible with the input format, so processing the output through `gcopy` should produce almost the same output,
although symbols may be slightly reordered.
Two applications of gcopy reaches the fixed point where further applications make no changes.

### gdeebnf — convert an *EBNF* grammar to *BNF*
Given the example *EBNF* grammar given above, stored in a file named `ebnf.gr`, type this shell command while in the `gtools` directory:

```bash
./gdeebnf < ebnf.gr
```

The output will be:

    > expression
    / ''
    
    expression ::= term expression-a
    term ::= factor term-a
    factor ::= factor-a factor-b
    factor-a ::= ''
    |  '-'
    factor-b ::= number
    |  identifier
    |  '(' expression ')'
    term-a ::= ''
    |  term-a-a factor term-a
    term-a-a ::= '*'
    |  '/'
    expression-a ::= ''
    |  expression-a-a term expression-a
    expression-a-a ::= '+'
    |  '-'
    
    # terminals:   '' '+' '-' '*' '/' number identifier '(' ')'

This output is fully compatible with the *BNF* expected by `gcopy` and the other tools,
and if input to `gcopy` it will be output with, at most,
a slight change in the order of the non-terminals listed at the end.

Note that the mechanical conversion from extended *BNF* to classical *BNF* involved the introduction of new symbols into the grammar.
The `gdeebnf` program generates new names for these by appending a suffix to the name of the symbol from which they were derived,
so all new non-terminals generated from processing the non-terminal term will have names like `term-a` or `term-b`.
The name term-a-a above was found while processing a rule hanging under term-a.

Automatically generated symbols never collide with symbols that were defined in the input.
If the name is derived from a symbol enclosed in quotes or angle brackets,
the name extension will be inside the enclosing marks, so for example,
`<term>` could give rise to `<term-a>` and `<term-b>`.

### gdeempty — eliminate empty symbols from a *BNF* grammar
The example *BNF* grammar derived by `gdeebnf` above is full of references to the empty symbol.
We can mechanically transform this to a form with no references to the empty symbol using the `gdeempty` tool, as follows:

```bash
./gdeebnf < ebnf.gr | ./gdeempty
```

The output will be:

    > expression
    
    expression ::= term expression-a
    |  term
    term ::= factor term-a
    |  factor
    factor ::= factor-a factor-b
    |  factor-b
    factor-a ::= '-'
    factor-b ::= number
    |  identifier
    |  '(' expression ')'
    term-a ::= term-a-a factor term-a
    |  term-a-a factor
    term-a-a ::= '*'
    |  '/'
    expression-a ::= expression-a-a term expression-a
    |  expression-a-a term
    expression-a-a ::= '+'
    |  '-'
    
    # terminals:   '-' number identifier '(' ')' '*' '/' '+'
    
    # unused terminals:  ''

Note that, after this transformation, the empty symbol `'''` is no-longer referenced anywhere in the grammar,
and the processor reports this with a comment at the end.

### gsqueeze — eliminate redundancy from a *BNF* grammar
The example *BNF* grammar derived by `gdeempty` above contains a redundant element,
but reading through the grammar to find and eliminate it is a tedious job.
We can mechanically remove such redundancies using the `gsqueeze` tool, as follows:

```bash
./gdeebnf < ebnf.gr | ./gdeempty | ./gsqueeze
```

The output will be:

    > expression
    
    expression ::= term expression-a
    |  term
    term ::= factor term-a
    |  factor
    factor ::= '-' factor-b
    |  factor-b
    factor-b ::= number
    |  identifier
    |  '(' expression ')'
    term-a ::= term-a-a factor term-a
    |  term-a-a factor
    term-a-a ::= '*'
    |  '/'
    expression-a ::= expression-a-a term expression-a
    |  expression-a-a term
    expression-a-a ::= '+'
    |  '-'
    
    # terminals:   '-' number identifier '(' ')' '*' '/' '+'
    
    # unused productions
    factor-a ::= '-'

The `gsqueeze` tool discovered that it could replace the non-terminal `factor-a` with the terminal symbol `'-'` wherever it occurred in the grammar.
It did this, and having done so, it found that one of the production rules in the grammar was no-longer needed.
This rule has been preserved, set off by a comment at the end.

Applying `gsqueeze` to grammars written with a pedagogical intent frequently finds large numbers of such substitutions.

### gsample — generate an example string from a *BNF* grammar
Sometimes, it is useful to see an example of a string in the language described by a grammar.
The `gsample` tool does this, generating a different random string each time it is run against a particular grammar.
In the current version, all the alternative rules associated with each symbol are equally weighted,
so where there is recursion in the grammar, it sometimes runs away, producing very long strings.
Consider the following command to apply `gsample` to the example *BNF* grammar:

```bash
./gsample < bnf.gr
```

Several runs of this produced the following outputs, as well as a number of much longer outputs,
many of which terminated with segmentation faults (stack overflows from runaway recursion).

    <number>
    - <number> * - <number>
    - <identifier> * ( - <number> )
    <number> / <number> * <identifier> / - <number>
    - ( - <identifier> / ( - <number> ) - <identifier> + <number> ) / <identifier>

Applying the same tool to an *EBNF* grammar requires first eliminating the use of *EBNF* notation:

```bash
./gdeebnf < ebnf.gr | ./gsample
```

Several runs of this produced the following outputs, as well as many longer ones.

    '-' identifier
    identifier '+' '(' '-' number ')' '/' identifier '/' number
    identifier '/' '-' '(' number ')' '*' '(' '(' '-' number ')' ')' '*' '-' number

Note that `gsample` does not include the empty symbol in its output,
so long as that symbol is properly declared in the grammar.
Note, also, that terminal symbols are output in exactly the form they were given in the grammar,
quotation marks and all.

The problem with very long strings in the output is caused by the fact that `gsample` equally weights all alternatives under each non-terminal.
If most of the alternatives under a non-terminal are recursive,
the result will produce long strings with a high probability.
If this is a problem, one option is to duplicate some of the non-recursive alternatives to increase their probability.
Consider this rule, with a 2/3 probability of recursion:

    <term> ::= <factor> | <term> * <factor> | <term> / <factor>

Replacing it with the following reduces the probability of recursion to 1/2:

    <term> ::= <factor> | <factor> | <term> * <factor> | <term> / <factor>

Applying `gsqueeze` to the modified rule above will restore it to its original form.

### gstartfollow — find the start and follow sets of all non-terminals
Many parsing algorithms require knowing two sets of symbols for each non-terminal in the grammar:

1. The start-set of a non-terminal contains all the terminal symbols that may occur as the first symbols of strings generated by that non-terminal.
2. The follow-set of a non-terminal contains all the terminal symbols that may occur immediately after a string generated by that non-terminal.

3. We can mechanically compute these sets with the `gstartfollow` tool, as follows:

```bash
./startfollow < bnf.gr
```

The output will be:

    > <expression>
    
    <expression> ::= <term>
                  |  <expression> + <term>
                  |  <expression> - <term>
    # start set:   - <number> <identifier> (
    # follow set:  + - )
    
    <term> ::= <factor>
            |  <term> * <factor>
            |  <term> / <factor>
    # start set:   - <number> <identifier> (
    # follow set:  * / + - )
    
    <factor> ::= <element>
              |  - <element>
    # start set:   - <number> <identifier> (
    # follow set:  * / + - )
    
    <element> ::= <number>
               |  <identifier>
               |  ( <expression> )
    # start set:   <number> <identifier> (
    # follow set:  * / + - )
    
    
    # terminals:   + - * / <number> <identifier> ( )

Note that the output is merely an annotation of the grammar;
aside from this, `gstartfollow` merely copies and pretty-prints the input grammar.

Note also that there is no explicit mention of the end-of-file as a terminal symbol following the `<expression>`.
If you want explicit consideration of the end-of-file character, it is best to explicitly include it in the grammar.
In the formal grammar literature, the symbols `⊢` and `⊣` (right tack and left tack) are frequently used for start and end of file.
For those who can't touch type Unicode, `/-` and `-/` might be preferable.
(Note, `|-` is not workable because of the special meaning for the unquoted vertical bar.)
The grammar can be rewritten to account for end of file by adding one new production rule at the start and declaring this to be the start symbol:

    > <file>
    <file> ::= /- <expression> -/

Note that `gstartfollow` treats the empty symbol as a normal terminal symbol.
As a result, the sets it identifies will be of limited use unless the grammar is first rewritten without the empty symbol,
for example, by applying `gdeempty`.
Also note that `gstartfollow` expects its input to be in simple *BNF*.
If you want to analyze an *EBNF* grammar, you must first convert it, for example:

```bash
./deebnf < ebnf.gr | ./gdeempty | ./gstartfollow
```

In the result, the start set and follow set for each non-terminal that was in the original *EBNF* grammar will be correct.
In addition, the start and follow sets for the added non-terminals will generally drive the logic within parsing of individual rules.

### gstats — count the rules and symbols of a *BNF* grammar
Particularly with large grammars, it is useful to get an overall view by simply counting the terminal symbols,
non-terminal symbols and production rules of a grammar.
The `gstats` tool does this, as follows:

```bash
./gstats < bnf.gr
```

The output will be:

    -- Total symbols:        12
    --   Terminal symbols:   8
    -- Production rules:     11

If there are symbols or rules that are not reachable from the distinguished or head symbol,
these are counted and identified as such.
When applied to an *EBNF* grammar, parentheses and brackets are not counted as introducing any new rules,
but they are counted as terminal symbols.
This may be misleading.

## Notes
### History
This software is descended from Pascal code I wrote back in the 1970s when I was involved in writing a Pascal compiler for the Modcomp IV computer.
I wrote it because I was disgusted with the extent to which Wirth's Pascal grammar was bloated by extra non-terminals introduced for what appeared to be pedagogical reasons.
At some point, after Pascal began to disappear from the computers available to me, I ran the code through `p2c` or some such converter.
The resulting C code was useful but annoying both because of conversion artifacts and because it was one big routine.
To change the filtering actions, you had to edit the main program, swapping in or out calls to various subroutines.

When I taught compiler construction in the Spring of 2013, I thought the code would be useful for my students, so I pared the mess down into a single utility to compute start and follow sets, and put that on the class web page.
The problem was, the code was basically embarrassing,
so I put some time in during the summer of 2013 to reconstruct bits and pieces of the original so I could release the components as stand-alone tools that could be strung together as filters,
in the long honored Unix shell scripting tradition.

### Links
[Robert Noonan](http://www.cs.wm.edu/~noonan/)
at the College of William and Mary has a set of
[grammar tools](http://www.cs.wm.edu/~noonan/Grammar/)
that primarily work using Wirth's extensions to *BNF*.
It should be easy to write a filter to convert Noonan's version of Wirth's format to be compatible with the tools here.