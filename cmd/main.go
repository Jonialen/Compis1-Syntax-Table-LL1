package main

import (
	"fmt"
	"sort"
	"strings"

	"syntaxtable/internal/grammar"
	"syntaxtable/internal/sets"
	"syntaxtable/internal/table"
)

// testGrammars are the grammars analyzed by this program.
// The first is required by the assignment; the remaining two are free choice.
var testGrammars = []struct {
	name  string
	input string
}{
	{
		// Required grammar from the assignment.
		name: "Expresiones Aritméticas",
		input: `
E  -> T E'
E' -> + T E' | ε
T  -> F T'
T' -> * F T' | ε
F  -> ( E ) | id`,
	},
	{
		// Chosen grammar #1: demonstrates how epsilon productions interact
		// when multiple nullable non-terminals follow each other.
		name: "Nullable Sequence  (S -> A B, con ε en A y B)",
		input: `
S -> A B
A -> a | ε
B -> b | ε`,
	},
	{
		// Chosen grammar #2: classic left-recursive grammar — demonstrates
		// why left recursion breaks LL(1) and forces conflicts in the table.
		name: "Expresiones con Recursión Izquierda (no LL(1))",
		input: `
E -> E + T | T
T -> T * F | F
F -> id | ( E )`,
	},
}

func main() {
	for _, tg := range testGrammars {
		printSeparator()
		fmt.Printf("Gramática: %s\n", tg.name)
		printSeparator()

		g, err := grammar.Parse(tg.input)
		if err != nil {
			fmt.Printf("Error al parsear la gramática: %v\n\n", err)
			continue
		}

		printGrammar(g)

		first := sets.Compute(g)
		printFirst(g, first)

		follow := sets.ComputeFollow(g, first)
		printFollow(g, follow)

		pt := table.Build(g, first, follow)
		fmt.Println("Tabla de Análisis Sintáctico Predictivo:")
		fmt.Println(pt.Print())

		if pt.IsLL1 {
			fmt.Println("RESULTADO: La gramática ES LL(1) — no hay conflictos en la tabla.")
		} else {
			fmt.Println("RESULTADO: La gramática NO es LL(1). Conflictos encontrados:")
			for _, c := range pt.Conflicts {
				fmt.Printf("  M[%s, %s] tiene %d producciones:\n",
					c.NonTerminal, c.Terminal, len(c.Productions))
				for _, idx := range c.Productions {
					fmt.Printf("    %s\n", g.Productions[idx])
				}
			}
		}
		fmt.Println()
	}
}

func printGrammar(g *grammar.Grammar) {
	fmt.Println("Producciones:")
	for _, p := range g.Productions {
		fmt.Printf("  %s\n", p)
	}
	fmt.Println()
}

func printFirst(g *grammar.Grammar, first sets.FirstSets) {
	fmt.Println("Conjuntos FIRST:")
	for _, nt := range g.NTOrder {
		fmt.Printf("  FIRST(%-4s) = { %s }\n", nt, sortedJoin(first[nt]))
	}
	fmt.Println()
}

func printFollow(g *grammar.Grammar, follow sets.FollowSets) {
	fmt.Println("Conjuntos FOLLOW:")
	for _, nt := range g.NTOrder {
		fmt.Printf("  FOLLOW(%-4s) = { %s }\n", nt, sortedJoin(follow[nt]))
	}
	fmt.Println()
}

func sortedJoin(s map[string]bool) string {
	elems := make([]string, 0, len(s))
	for k := range s {
		elems = append(elems, k)
	}
	sort.Strings(elems)
	return strings.Join(elems, ", ")
}

func printSeparator() {
	fmt.Println("═══════════════════════════════════════════════════════")
}
