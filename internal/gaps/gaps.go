package gaps

import (
	"math"
	"strings"
)

// ---------------- utilidades de posiciones ----------------

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
// Optimizado para usar menos memoria: procesa gaps incrementalmente sin almacenar rutas completas.
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

	// Búsqueda binaria para encontrar el primer elemento > minVal
	findStart := func(arr []int, minVal int) int {
		left, right := 0, len(arr)
		for left < right {
			mid := (left + right) / 2
			if arr[mid] <= minVal {
				left = mid + 1
			} else {
				right = mid
			}
		}
		return left
	}

	// DFS optimizado: registra gaps inmediatamente, sin almacenar toda la ruta
	var dfs func(k int, prevPos int) bool
	dfs = func(k int, prevPos int) bool {
		// k: posición actual en el patrón
		// prevPos: posición en seq del carácter anterior (o -1 si k==0)
		
		if k == L { // Llegamos al final de una incrustación válida
			return true
		}

		v := occ[k]
		
		// Búsqueda binaria para encontrar primer v[i] > prevPos
		i0 := 0
		if prevPos >= 0 {
			i0 = findStart(v, prevPos)
		}
		
		// Poda: si no hay ocurrencias válidas, retornar false
		if i0 >= len(v) {
			return false
		}
		
		// Poda temprana: verificar que podemos alcanzar el último carácter
		if k+1 < L && v[len(v)-1] >= lat[k+1] {
			// Encontrar la última posición válida
			lastValid := len(v) - 1
			for lastValid >= i0 && v[lastValid] >= lat[k+1] {
				lastValid--
			}
			if lastValid < i0 {
				return false
			}
		}
		
		okAny := false
		for i := i0; i < len(v); i++ {
			currentPos := v[i]
			
			// Poda: si la posición actual ya supera o iguala latest del siguiente carácter
			if k+1 < L && currentPos >= lat[k+1] {
				break
			}
			
			// Si no es el primer carácter, registrar el gap inmediatamente
			if k > 0 {
				gap := currentPos - prevPos - 1
				sets[k-1][gap] = struct{}{}
			}
			
			// Continuar DFS
			if dfs(k+1, currentPos) {
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
