package pattern

import "strings"

// alignPatternIndices encuentra índices en s para cada letra de pattern,
// usando greedy izquierda->derecha y luego un "empuje a la derecha" local
// para maximizar gaps sin romper el orden. Devuelve los índices de coincidencia.
func alignPatternIndices(pattern string, s string) ([]int, bool) {
	matches := make([]int, 0, len(pattern))
	start := 0
	for i := 0; i < len(pattern); i++ {
		ch := pattern[i]
		pos := strings.IndexByte(s[start:], ch)
		if pos < 0 {
			return nil, false
		}
		pos += start
		matches = append(matches, pos)
		start = pos + 1
	}

	// Empuje a la derecha local: para cada match (menos el último),
	// intenta moverlo a la última aparición posible antes del siguiente.
	for i := 0; i+1 < len(matches); i++ {
		ch := pattern[i]
		cur := matches[i]
		bound := matches[i+1]
		best := cur
		seek := cur + 1
		for {
			// buscar próxima ocurrencia del mismo ch
			nextRel := strings.IndexByte(s[seek:bound], ch)
			if nextRel < 0 {
				break
			}
			nextAbs := seek + nextRel
			best = nextAbs
			seek = nextAbs + 1
		}
		matches[i] = best
	}

	return matches, true
}

// lowercaseBetween cuenta minúsculas estrictamente entre (i, j) en s usando prefijo.
func lowercaseBetween(prefLower []int, i, j int) int {
	if j <= i {
		return 0
	}
	// contamos minúsculas en (i, j) => prefLower[j] - prefLower[i+1]
	return prefLower[j] - prefLower[i+1]
}

// gapsForPattern devuelve, para un pattern incrustado en s, las cantidades de minúsculas
// entre letras consecutivas. Si el pattern no se incrusta, retorna false.
func gapsForPattern(pattern string, proj UpperProjection) ([]int, bool) {
	indices, ok := alignPatternIndices(pattern, proj.Original)
	if !ok {
		return nil, false
	}
	if len(indices) <= 1 {
		return []int{}, true
	}
	gaps := make([]int, 0, len(indices)-1)
	for i := 0; i+1 < len(indices); i++ {
		a, b := indices[i], indices[i+1]
		gaps = append(gaps, lowercaseBetween(proj.PrefLower, a, b))
	}
	return gaps, true
}
