# Syntax Table LL(1)

Implementación de los conjuntos **FIRST** y **FOLLOW** y la construcción automática de la **tabla de análisis sintáctico predictivo LL(1)** para gramáticas libres de contexto arbitrarias.

## Video explicativo

[![Video explicativo](https://img.youtube.com/vi/bSVQWcwBks0/0.jpg)](https://youtu.be/bSVQWcwBks0)

## Estructura del proyecto

```
SyntaxTableLL1/
├── cmd/
│   └── main.go                 # Punto de entrada; ejecuta las tres gramáticas de prueba
├── internal/
│   ├── grammar/
│   │   └── grammar.go          # Tipos Production y Grammar; parser de texto plano
│   ├── sets/
│   │   ├── first.go            # Cálculo de conjuntos FIRST (punto fijo)
│   │   └── follow.go           # Cálculo de conjuntos FOLLOW (punto fijo)
│   └── table/
│       └── table.go            # Construcción de la tabla LL(1) y detección de conflictos
└── go.mod
```

## Uso

```bash
go run ./cmd/main.go
```

### Formato de entrada

Las gramáticas se definen como cadenas de texto. Cada línea es una regla; las alternativas se separan con `|`; la épsilon se escribe `ε` o `eps`.

```
E  -> T E'
E' -> + T E' | ε
T  -> F T'
T' -> * F T' | ε
F  -> ( E ) | id
```

Para usar el parser desde otro código:

```go
g, err := grammar.Parse(input)
first := sets.Compute(g)
follow := sets.ComputeFollow(g, first)
pt := table.Build(g, first, follow)
fmt.Println(pt.Print())
```

## Algoritmos

### FIRST

Se calcula por iteración hasta punto fijo. Para una secuencia de símbolos `X1 X2 … Xn`:

1. Agregar `FIRST(X1) – {ε}` al resultado.
2. Si `ε ∈ FIRST(X1)`, continuar con `X2`; de lo contrario, detener.
3. Si todos los símbolos pueden derivar `ε`, agregar `ε` al resultado.

### FOLLOW

También por iteración hasta punto fijo. Para cada producción `A → α`:

- **Regla 1:** `$ ∈ FOLLOW(S)` donde S es el símbolo inicial.
- **Regla 2:** Para cada no-terminal B en posición i con sufijo `β = α[i+1:]`:  
  `FIRST(β) – {ε} ⊆ FOLLOW(B)`
- **Regla 3:** Si `ε ∈ FIRST(β)`:  
  `FOLLOW(A) ⊆ FOLLOW(B)`

### Tabla LL(1)

Para cada producción `A → α`:

- Para cada `a ∈ FIRST(α)`: `M[A, a] = A → α`
- Si `ε ∈ FIRST(α)`: para cada `b ∈ FOLLOW(A)`: `M[A, b] = A → α`

Si alguna celda recibe más de una producción, existe un conflicto y la gramática **no es LL(1)**.

## Gramáticas probadas

### 1. Expresiones aritméticas (requerida)

```
E  -> T E'
E' -> + T E' | ε
T  -> F T'
T' -> * F T' | ε
F  -> ( E ) | id
```

| No-terminal | FIRST | FOLLOW |
|-------------|-------|--------|
| E | `(, id` | `$, )` |
| E' | `+, ε` | `$, )` |
| T | `(, id` | `$, ), +` |
| T' | `*, ε` | `$, ), +` |
| F | `(, id` | `$, ), *, +` |

**¿Es LL(1)?** Sí — no hay conflictos.

---

### 2. Secuencia con no-terminales anulables

```
S -> A B
A -> a | ε
B -> b | ε
```

Elegida porque ejercita el caso más delicado del algoritmo de FOLLOW: cuando el sufijo `β` que sigue a un no-terminal puede derivar `ε` por completo, el no-terminal hereda el FOLLOW de la cabeza de la producción. En esta gramática, A está seguida de B (que es anulable), por lo que `FOLLOW(A) = {b, $}`.

| No-terminal | FIRST | FOLLOW |
|-------------|-------|--------|
| S | `a, b, ε` | `$` |
| A | `a, ε` | `b, $` |
| B | `b, ε` | `$` |

**¿Es LL(1)?** Sí — no hay conflictos.

---

### 3. Expresiones con recursión izquierda (no LL(1))

```
E -> E + T | T
T -> T * F | F
F -> id | ( E )
```

Elegida para demostrar por qué la recursión izquierda es incompatible con el análisis predictivo. Todas las producciones de E comparten el mismo FIRST `{id, (}`, lo que genera conflictos inevitables en la tabla: el parser no puede distinguir entre `E → E + T` y `E → T` con un solo token de anticipación.

| No-terminal | FIRST | FOLLOW |
|-------------|-------|--------|
| E | `(, id` | `$, ), +` |
| T | `(, id` | `$, ), *, +` |
| F | `(, id` | `$, ), *, +` |

**¿Es LL(1)?** No — conflictos en `M[E,(]`, `M[E,id]`, `M[T,(]`, `M[T,id]`.

La solución estándar es eliminar la recursión izquierda reescribiendo la gramática (como se hace en la gramática 1), lo que produce una gramática equivalente que sí es LL(1).
