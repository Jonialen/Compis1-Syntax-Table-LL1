package grammar

import (
	"fmt"
	"strings"
)

const (
	Epsilon   = "ε"
	EndMarker = "$"
)

// Production represents a single grammar rule: Head -> Body.
type Production struct {
	Head string
	Body []string // ["ε"] means an epsilon production
}

func (p Production) String() string {
	if len(p.Body) == 1 && p.Body[0] == Epsilon {
		return p.Head + " → ε"
	}
	return p.Head + " → " + strings.Join(p.Body, " ")
}

// Grammar holds a complete context-free grammar.
type Grammar struct {
	NonTerminals map[string]bool
	Terminals    map[string]bool
	Productions  []Production
	Start        string
	NTOrder      []string // non-terminals in declaration order (for stable display)
}

// Parse builds a Grammar from a multi-line string.
//
// Format: one rule per line, alternatives separated by |
//
//	E  -> T E'
//	E' -> + T E' | ε
//	F  -> ( E ) | id
//
// Epsilon may be written as "ε", "eps", or "epsilon".
// Lines starting with '#' are ignored.
func Parse(input string) (*Grammar, error) {
	g := &Grammar{
		NonTerminals: make(map[string]bool),
		Terminals:    make(map[string]bool),
	}

	seenNT := make(map[string]bool)

	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "->", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid rule %q: expected '->'", line)
		}

		head := strings.TrimSpace(parts[0])
		if !seenNT[head] {
			seenNT[head] = true
			g.NTOrder = append(g.NTOrder, head)
			g.NonTerminals[head] = true
			if g.Start == "" {
				g.Start = head
			}
		}

		for _, alt := range strings.Split(parts[1], "|") {
			alt = strings.TrimSpace(alt)
			var body []string
			switch alt {
			case "ε", "eps", "epsilon", "":
				body = []string{Epsilon}
			default:
				body = strings.Fields(alt)
			}
			g.Productions = append(g.Productions, Production{Head: head, Body: body})
		}
	}

	if g.Start == "" {
		return nil, fmt.Errorf("empty grammar")
	}

	// Every symbol that appears in a rule body and is not a declared non-terminal is a terminal.
	for _, prod := range g.Productions {
		for _, sym := range prod.Body {
			if sym != Epsilon && !g.NonTerminals[sym] {
				g.Terminals[sym] = true
			}
		}
	}

	return g, nil
}
