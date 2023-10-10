// gtools - a collection of grammar manipulation tools
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package gtools

// Interface

// bool samerule( PPRODUCTION p, PPRODUCTION q )
// compare rules p,q
func samerule(p PPRODUCTION, q PPRODUCTION) bool {
	var pe, qe PELEMENT /* pointers to elements of p and q */

	/* for successive elements */
	pe = p.data
	qe = q.data
	for (pe != nil) && (qe != nil) {
		if pe.data != qe.data {
			return false /* they differ */
		}
		pe = pe.next
		qe = qe.next
	}
	return ((pe == nil) && (qe == nil)) /* identical and same length */
}
