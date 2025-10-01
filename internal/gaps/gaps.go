package gaps

import (
	"math"
	"strings"
)

// ---------------- utilidades de posiciones ----------------

func positionsByUpper(s string) map[byte][]int {
	pos := make(map[byte][]int, 26)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			pos[c] = append(pos[c], i)
		}
	}
	return pos
}

func indicesOfInRange(s string, ch byte, lo, hi int) []int {
	if lo < 0 {
		lo = 0
	}
	if hi >= len(s) {
		hi = len(s) - 1
	}
	if lo > hi {
		return nil
	}
	out := make([]int, 0)
	i := lo
	for {
		k := strings.IndexByte(s[i:hi+1], ch)
		if k < 0 {
			break
		}
		idx := i + k
		out = append(out, idx)
		i = idx + 1
	}
	return out
}

// ---------------- bandas de viabilidad ----------------

func earliestPositions(pattern, s string) ([]int, bool) {
	pos := make([]int, len(pattern))
	start := 0
	for i := 0; i < len(pattern); i++ {
		ch := pattern[i]
		off := strings.IndexByte(s[start:], ch)
		if off < 0 {
			return nil, false
		}
		pos[i] = start + off
		start = pos[i] + 1
	}
	return pos, true
}

func latestPositions(pattern, s string) ([]int, bool) {
	pos := make([]int, len(pattern))
	end := len(s)
	for i := len(pattern) - 1; i >= 0; i-- {
		ch := pattern[i]
		last := -1
		for j := end - 1; j >= 0; j-- {
			if s[j] == ch {
				last = j
				break
			}
		}
		if last < 0 {
			return nil, false
		}
		pos[i] = last
		end = last
	}
	return pos, true
}

// ---------------- mínimos (por si los sigues usando en algún modo) ----------------

func minGapForPairDistanceTotal(pattern string, s string, i int, ear, lat []int) (int, bool) {
	loA, hiA := ear[i], lat[i]
	loB, hiB := ear[i+1], lat[i+1]

	occA := indicesOfInRange(s, pattern[i], loA, hiA)
	occB := indicesOfInRange(s, pattern[i+1], loB, hiB)
	if len(occA) == 0 || len(occB) == 0 {
		return 0, false
	}

	minGap := math.MaxInt
	j := 0
	for _, a := range occA {
		for j < len(occB) && occB[j] <= a {
			j++
		}
		if j == len(occB) {
			break
		}
		b := occB[j]
		val := b - a - 1 // distancia total
		if val < minGap {
			minGap = val
		}
	}
	if minGap == math.MaxInt {
		return 0, false
	}
	return minGap, true
}

// MinGapsDistanceTotalViable: conserva tu modo "mínimo por secuencia".
func MinGapsDistanceTotalViable(seq string, pattern string) ([]int, bool) {
	if len(pattern) <= 1 {
		return []int{}, true
	}
	ear, ok := earliestPositions(pattern, seq)
	if !ok {
		return nil, false
	}
	lat, ok := latestPositions(pattern, seq)
	if !ok {
		return nil, false
	}
	mins := make([]int, 0, len(pattern)-1)
	for i := 0; i+1 < len(pattern); i++ {
		g, ok := minGapForPairDistanceTotal(pattern, seq, i, ear, lat)
		if !ok {
			return nil, false
		}
		mins = append(mins, g)
	}
	return mins, true
}

// ---------------- NUEVO: TODOS los gaps posibles por par (distancia total) ----------------

// AllGapValuesDistanceTotalViable devuelve, para UNA secuencia y un patrón,
// el conjunto de TODOS los valores posibles de b-a-1 por cada par consecutivo,
// considerando SOLO incrustaciones viables (bandas earliest/latest).
// Retorna un slice de sets (map[int]struct{}) de longitud len(pattern)-1.
func AllGapValuesDistanceTotalViable(seq string, pattern string) ([]map[int]struct{}, bool) {
	L := len(pattern)
	if L <= 1 {
		return make([]map[int]struct{}, 0), true
	}
	ear, ok := earliestPositions(pattern, seq)
	if !ok {
		return nil, false
	}
	lat, ok := latestPositions(pattern, seq)
	if !ok {
		return nil, false
	}

	// Para cada letra del patrón, precalculamos las ocurrencias viables en su banda [ear, lat].
	occ := make([][]int, L)
	for k := 0; k < L; k++ {
		occ[k] = indicesOfInRange(seq, pattern[k], ear[k], lat[k])
		if len(occ[k]) == 0 {
			return nil, false
		}
	}

	// sets[k] acumula TODOS los valores de gap entre pattern[k] -> pattern[k+1]
	sets := make([]map[int]struct{}, L-1)
	for i := range sets {
		sets[i] = make(map[int]struct{})
	}

	path := make([]int, L) // índices elegidos para cada letra del patrón

	var dfs func(k int, prev int) bool
	dfs = func(k int, prev int) bool {
		// k: posición en el patrón a elegir
		// prev: índice elegido del carácter anterior en seq (o -1 si k==0)
		if k == L { // tenemos una incrustación completa
			// acumular gaps de toda la ruta
			for i := 0; i+1 < L; i++ {
				a, b := path[i], path[i+1]
				sets[i][b-a-1] = struct{}{}
			}
			return true
		}
		// iterar ocurrencias viables para pattern[k], respetando orden creciente > prev
		v := occ[k]
		// avance inicial para mantener a > prev
		i0 := 0
		if prev >= 0 {
			// buscar primer v[i] > prev (lineal; si quieres, haz binaria)
			for i0 < len(v) && v[i0] <= prev {
				i0++
			}
		}
		okAny := false
		for i := i0; i < len(v); i++ {
			path[k] = v[i]
			// poda de viabilidad futura: aún podemos alcanzar latest[k+1], etc.
			// (ear/latest ya garantizan que hay al menos una completación posible)
			if dfs(k+1, v[i]) {
				okAny = true
			}
		}
		return okAny
	}

	if !dfs(0, -1) {
		return nil, false
	}
	return sets, true
}
