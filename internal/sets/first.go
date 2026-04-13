package sets

import "syntaxtable/internal/grammar"

// FirstSets maps every symbol to its FIRST set.
type FirstSets map[string]map[string]bool

// Compute returns the FIRST sets for all symbols in g using a fixed-point iteration.
//
// Algorithm:
//  1. FIRST(terminal) = {terminal}
//  2. For A → ε: ε ∈ FIRST(A)
//  3. For A → X1 X2 … Xn:
//     - Add FIRST(Xi) – {ε} to FIRST(A)
//     - If ε ∈ FIRST(Xi), continue with Xi+1; otherwise stop
//     - If all Xi derive ε, add ε to FIRST(A)
//  4. Repeat until no change (fixed point).
func Compute(g *grammar.Grammar) FirstSets {
	first := make(FirstSets)

	// Seed: non-terminals start empty, terminals map to themselves.
	for nt := range g.NonTerminals {
		first[nt] = make(map[string]bool)
	}
	for t := range g.Terminals {
		first[t] = map[string]bool{t: true}
	}
	first[grammar.Epsilon] = map[string]bool{grammar.Epsilon: true}

	for {
		changed := false
		for _, prod := range g.Productions {
			for sym := range OfBody(first, prod.Body) {
				if !first[prod.Head][sym] {
					first[prod.Head][sym] = true
					changed = true
				}
			}
		}
		if !changed {
			break
		}
	}

	return first
}

// OfBody computes FIRST of an arbitrary sequence of grammar symbols.
// Exported so both the FOLLOW computation and the table builder can reuse it.
func OfBody(first FirstSets, body []string) map[string]bool {
	result := make(map[string]bool)

	// Empty body or explicit epsilon production → derives ε
	if len(body) == 0 || (len(body) == 1 && body[0] == grammar.Epsilon) {
		result[grammar.Epsilon] = true
		return result
	}

	allDeriveEpsilon := true
	for _, sym := range body {
		symFirst := first[sym]
		if symFirst == nil {
			// Unknown symbol: treat as a terminal (can happen during construction)
			result[sym] = true
			allDeriveEpsilon = false
			break
		}
		for s := range symFirst {
			if s != grammar.Epsilon {
				result[s] = true
			}
		}
		if !symFirst[grammar.Epsilon] {
			allDeriveEpsilon = false
			break
		}
	}

	if allDeriveEpsilon {
		result[grammar.Epsilon] = true
	}

	return result
}
