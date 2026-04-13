package sets

import "syntaxtable/internal/grammar"

// FollowSets maps each non-terminal to its FOLLOW set.
type FollowSets map[string]map[string]bool

// ComputeFollow returns the FOLLOW sets for all non-terminals using a fixed-point iteration.
//
// Rules applied for each production A → α:
//  1. $ ∈ FOLLOW(Start)
//  2. For every non-terminal B in α at position i with β = α[i+1:]:
//     FIRST(β) – {ε} ⊆ FOLLOW(B)
//  3. If ε ∈ FIRST(β):  FOLLOW(A) ⊆ FOLLOW(B)
func ComputeFollow(g *grammar.Grammar, first FirstSets) FollowSets {
	follow := make(FollowSets)
	for nt := range g.NonTerminals {
		follow[nt] = make(map[string]bool)
	}

	// Rule 1
	follow[g.Start][grammar.EndMarker] = true

	for {
		changed := false
		for _, prod := range g.Productions {
			for i, sym := range prod.Body {
				if !g.NonTerminals[sym] {
					continue
				}

				// β = everything after position i
				beta := prod.Body[i+1:]
				betaFirst := OfBody(first, beta)

				// Rule 2: FIRST(β) – {ε} ⊆ FOLLOW(sym)
				for s := range betaFirst {
					if s != grammar.Epsilon && !follow[sym][s] {
						follow[sym][s] = true
						changed = true
					}
				}

				// Rule 3: if ε ∈ FIRST(β), FOLLOW(Head) ⊆ FOLLOW(sym)
				if betaFirst[grammar.Epsilon] {
					for s := range follow[prod.Head] {
						if !follow[sym][s] {
							follow[sym][s] = true
							changed = true
						}
					}
				}
			}
		}
		if !changed {
			break
		}
	}

	return follow
}
