package aggregate

import (
	"fmt"
	"sort"
)

// GapValues guarda el conjunto de valores discretos observados para un gap.
type GapValues struct {
	Values []int // únicos y ordenados
}

// PairMinValues: toma los mínimos por gap de dos secuencias y devuelve los conjuntos.
// Con dos secuencias, cada gap tendrá 1 o 2 valores (si difieren).
func PairMinValues(minsX, minsY []int) []GapValues {
	n := len(minsX)
	if len(minsY) < n {
		n = len(minsY)
	}
	out := make([]GapValues, n)
	for i := 0; i < n; i++ {
		a, b := minsX[i], minsY[i]
		if a == b {
			out[i] = GapValues{Values: []int{a}}
		} else {
			if a > b {
				a, b = b, a
			}
			out[i] = GapValues{Values: []int{a, b}}
		}
	}
	return out
}

// helper: ¿la lista es un bloque contiguo de enteros?
func isContiguous(sorted []int) bool {
	if len(sorted) <= 1 {
		return true
	}
	for i := 1; i < len(sorted); i++ {
		if sorted[i] != sorted[i-1]+1 {
			return false
		}
	}
	return true
}

// FormatPatternWithValues imprime P-x(...)-Q-x(...)-…
// Regla:
// - si Values = {k}  => x(k)
// - si Values cubren todos los enteros entre min y max => x(min,max)
// - si no, x(v1|v2|...|vt)
func FormatPatternWithValues(pattern string, sets []GapValues) string {
	if len(pattern) == 0 {
		return ""
	}
	out := make([]byte, 0, len(pattern)*4)
	for i := 0; i < len(pattern); i++ {
		out = append(out, pattern[i])
		if i+1 < len(pattern) && i < len(sets) {
			vals := append([]int(nil), sets[i].Values...)
			sort.Ints(vals)
			switch len(vals) {
			case 0:
				out = append(out, []byte("-")...)
			case 1:
				out = append(out, []byte(fmt.Sprintf("-x(%d)-", vals[0]))...)
			default:
				if isContiguous(vals) {
					out = append(out, []byte(fmt.Sprintf("-x(%d,%d)-", vals[0], vals[len(vals)-1]))...)
				} else {
					// listado explícito
					str := "x("
					for j, v := range vals {
						if j > 0 {
							str += "|"
						}
						str += fmt.Sprintf("%d", v)
					}
					str += ")"
					out = append(out, []byte("-"+str+"-")...)
				}
			}
		}
	}
	return string(out)
}
func PairUnionSets(setsX, setsY []map[int]struct{}) []GapValues {
	n := len(setsX)
	if len(setsY) < n {
		n = len(setsY)
	}
	out := make([]GapValues, n)
	for i := 0; i < n; i++ {
		union := make(map[int]struct{})
		for v := range setsX[i] {
			union[v] = struct{}{}
		}
		for v := range setsY[i] {
			union[v] = struct{}{}
		}
		// pasar a slice ordenado
		vals := make([]int, 0, len(union))
		for v := range union {
			vals = append(vals, v)
		}
		sort.Ints(vals)
		// si solo hay 0 no tiene sentido guardarlo
		if len(vals) == 1 && vals[0] == 0 {
			vals = nil
		}
		out[i] = GapValues{Values: vals}
	}
	return out
}
