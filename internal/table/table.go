package table

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"syntaxtable/internal/grammar"
	"syntaxtable/internal/sets"
)

// Conflict records a cell in the parse table that holds more than one production.
type Conflict struct {
	NonTerminal string
	Terminal    string
	Productions []int // indices into Grammar.Productions
}

// ParseTable is the LL(1) predictive parsing table.
type ParseTable struct {
	// Table[nt][terminal] = list of production indices (>1 means a conflict)
	Table     map[string]map[string][]int
	Conflicts []Conflict
	IsLL1     bool
	g         *grammar.Grammar
}

// Build constructs the LL(1) parse table for g.
//
// For each production A → α:
//   - For each terminal a ∈ FIRST(α): M[A, a] += (A → α)
//   - If ε ∈ FIRST(α):  for each b ∈ FOLLOW(A): M[A, b] += (A → α)
func Build(g *grammar.Grammar, first sets.FirstSets, follow sets.FollowSets) *ParseTable {
	pt := &ParseTable{
		Table: make(map[string]map[string][]int),
		IsLL1: true,
		g:     g,
	}
	for nt := range g.NonTerminals {
		pt.Table[nt] = make(map[string][]int)
	}

	for i, prod := range g.Productions {
		bodyFirst := sets.OfBody(first, prod.Body)

		for a := range bodyFirst {
			if a != grammar.Epsilon {
				pt.add(prod.Head, a, i)
			}
		}

		if bodyFirst[grammar.Epsilon] {
			for b := range follow[prod.Head] {
				pt.add(prod.Head, b, i)
			}
		}
	}

	// Collect conflicts (deterministic order)
	for nt, row := range pt.Table {
		for terminal, prods := range row {
			if len(prods) > 1 {
				pt.IsLL1 = false
				pt.Conflicts = append(pt.Conflicts, Conflict{
					NonTerminal: nt,
					Terminal:    terminal,
					Productions: prods,
				})
			}
		}
	}
	sort.Slice(pt.Conflicts, func(i, j int) bool {
		if pt.Conflicts[i].NonTerminal != pt.Conflicts[j].NonTerminal {
			return pt.Conflicts[i].NonTerminal < pt.Conflicts[j].NonTerminal
		}
		return pt.Conflicts[i].Terminal < pt.Conflicts[j].Terminal
	})

	return pt
}

// add inserts production index idx into cell M[nt][terminal], avoiding duplicates.
func (pt *ParseTable) add(nt, terminal string, idx int) {
	for _, existing := range pt.Table[nt][terminal] {
		if existing == idx {
			return
		}
	}
	pt.Table[nt][terminal] = append(pt.Table[nt][terminal], idx)
}

// Print returns the parse table formatted as a tab-aligned grid.
func (pt *ParseTable) Print() string {
	g := pt.g

	// Collect terminals from the grammar; append $ at the end.
	var terminals []string
	for t := range g.Terminals {
		terminals = append(terminals, t)
	}
	sort.Strings(terminals)
	terminals = append(terminals, grammar.EndMarker)

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 1, 1, 2, ' ', 0)

	// Header row
	fmt.Fprint(w, " ")
	for _, t := range terminals {
		fmt.Fprintf(w, "\t%s", t)
	}
	fmt.Fprintln(w)

	// One row per non-terminal (in declaration order)
	for _, nt := range g.NTOrder {
		fmt.Fprint(w, nt)
		for _, t := range terminals {
			cell := ""
			if prods, ok := pt.Table[nt][t]; ok && len(prods) > 0 {
				var parts []string
				for _, idx := range prods {
					parts = append(parts, g.Productions[idx].String())
				}
				cell = strings.Join(parts, " / ")
			}
			fmt.Fprintf(w, "\t%s", cell)
		}
		fmt.Fprintln(w)
	}

	w.Flush()
	return buf.String()
}
