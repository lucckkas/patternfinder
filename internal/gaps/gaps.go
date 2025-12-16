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

// ---------------- CONSOLIDACIÓN DE PATRONES ----------------

// PatternStat almacena estadísticas de un patrón
type PatternStat struct {
	Pattern         string
	UppercaseCount  int
	SequenceIndices map[int]bool // Índices de secuencias que contienen este patrón
}

// patternGroup agrupa patrones con la misma base
type patternGroup struct {
	letters    string
	gapCount   int
	gapValues  [][]int      // gapValues[i] = valores del gap i
	seqIndices map[int]bool
	patterns   []string // patrones originales
}

// parsePattern extrae la base del patrón (letras mayúsculas) y los valores de gaps
// Ahora devuelve un slice de gaps donde gaps[i] es el gap DESPUÉS de la letra i
// Si no hay gap después de una letra, el valor es -1
// Ejemplo: "C-x(2)-H-C" -> letters="CHC", gaps=[2, -1] (gap después de C, nada después de H)
func parsePattern(pattern string) (string, []int) {
	var letters strings.Builder
	gapsAfter := []int{} // gap después de cada letra (-1 si no hay)

	i := 0
	for i < len(pattern) {
		r := rune(pattern[i])

		// Si es una letra mayúscula, agregarla
		if r >= 'A' && r <= 'Z' {
			letters.WriteRune(r)
			// Por defecto, no hay gap después de esta letra
			gapsAfter = append(gapsAfter, -1)
			i++
			continue
		}

		// Si encontramos x(, extraer el número y asociarlo a la última letra
		if i+2 < len(pattern) && pattern[i] == 'x' && pattern[i+1] == '(' {
			i += 2 // saltar "x("
			numStart := i
			for i < len(pattern) && pattern[i] >= '0' && pattern[i] <= '9' {
				i++
			}
			if numStart < i && len(gapsAfter) > 0 {
				num := 0
				for j := numStart; j < i; j++ {
					num = num*10 + int(pattern[j]-'0')
				}
				// Asignar el gap a la última letra agregada
				gapsAfter[len(gapsAfter)-1] = num
			}
			// Saltar hasta después del ')'
			for i < len(pattern) && pattern[i] != ')' {
				i++
			}
			if i < len(pattern) {
				i++ // saltar ')'
			}
			continue
		}

		i++
	}

	// Quitar el último elemento de gapsAfter (no hay gap después de la última letra)
	if len(gapsAfter) > 0 {
		gapsAfter = gapsAfter[:len(gapsAfter)-1]
	}

	return letters.String(), gapsAfter
}

// countUppercaseInPattern cuenta las letras mayúsculas en una cadena
func countUppercaseInPattern(s string) int {
	count := 0
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			count++
		}
	}
	return count
}

// sortInts ordena un slice de enteros in-place
func sortInts(arr []int) {
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
}

// formatGapRange formatea un rango de gaps consecutivos
// Asume que los gaps ya son consecutivos
func formatGapRange(gaps []int) string {
	if len(gaps) == 0 {
		return ""
	}
	if len(gaps) == 1 {
		return "(" + itoa(gaps[0]) + ")"
	}
	// Formato (min,max) para rangos consecutivos
	return "(" + itoa(gaps[0]) + "," + itoa(gaps[len(gaps)-1]) + ")"
}

// findConsecutiveRanges divide una lista de enteros en rangos consecutivos
// Ejemplo: [1, 3, 11, 13, 15, 16, 18] -> [[1], [3], [11], [13], [15,16], [18]]
func findConsecutiveRanges(gaps []int) [][]int {
	if len(gaps) == 0 {
		return [][]int{}
	}
	if len(gaps) == 1 {
		return [][]int{{gaps[0]}}
	}

	ranges := [][]int{}
	currentRange := []int{gaps[0]}

	for i := 1; i < len(gaps); i++ {
		if gaps[i] == gaps[i-1]+1 {
			// Es consecutivo, agregar al rango actual
			currentRange = append(currentRange, gaps[i])
		} else {
			// No es consecutivo, guardar rango actual y empezar uno nuevo
			ranges = append(ranges, currentRange)
			currentRange = []int{gaps[i]}
		}
	}
	// Agregar el último rango
	ranges = append(ranges, currentRange)

	return ranges
}

// itoa convierte un entero a string (sin dependencias adicionales)
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	digits := []byte{}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}

// buildConsolidatedPattern construye el patrón final con gaps consolidados
// gapValues[i] contiene los valores del gap DESPUÉS de la letra i
// Si gapValues[i] está vacío, no hay gap después de esa letra
func buildConsolidatedPattern(letters string, gapValues [][]int) string {
	if len(letters) == 0 {
		return ""
	}

	var result strings.Builder
	letterRunes := []rune(letters)

	for i, letter := range letterRunes {
		if i > 0 {
			result.WriteString("-")
		}
		result.WriteRune(letter)

		// Agregar gap si existe y no es la última letra
		if i < len(gapValues) {
			gaps := gapValues[i]
			if len(gaps) > 0 {
				result.WriteString("-x")
				result.WriteString(formatGapRange(gaps))
			}
		}
	}

	return result.String()
}

// ConsolidatePatterns agrupa patrones con gaps consecutivos
// Ejemplo: C-x(2)-H y C-x(3)-H -> C-x(2,3)-H
// Si los gaps no son consecutivos, genera múltiples patrones
func ConsolidatePatterns(stats map[string]*PatternStat) map[string]*PatternStat {
	groups := make(map[string]*patternGroup)

	for pattern, stat := range stats {
		letters, gaps := parsePattern(pattern)
		
		// La clave incluye las letras y las posiciones de los gaps (cuáles tienen gap y cuáles no)
		// Esto asegura que solo se agrupen patrones con la misma estructura
		gapPositions := ""
		for i, g := range gaps {
			if g >= 0 {
				gapPositions += itoa(i) + ","
			}
		}
		key := letters + "|" + gapPositions

		if _, exists := groups[key]; !exists {
			gapValues := make([][]int, len(gaps))
			for i := range gapValues {
				gapValues[i] = []int{}
			}
			groups[key] = &patternGroup{
				letters:    letters,
				gapCount:   len(gaps),
				gapValues:  gapValues,
				seqIndices: make(map[int]bool),
				patterns:   []string{},
			}
		}

		g := groups[key]
		g.patterns = append(g.patterns, pattern)

		// Agregar valores de gaps (solo si el gap existe, es decir >= 0)
		for i, gapVal := range gaps {
			if gapVal < 0 {
				continue // No hay gap en esta posición
			}
			// Verificar si ya existe
			found := false
			for _, v := range g.gapValues[i] {
				if v == gapVal {
					found = true
					break
				}
			}
			if !found {
				g.gapValues[i] = append(g.gapValues[i], gapVal)
			}
		}

		// Unir secuencias
		for seqIdx := range stat.SequenceIndices {
			g.seqIndices[seqIdx] = true
		}
	}

	// Construir patrones consolidados
	result := make(map[string]*PatternStat)

	for _, g := range groups {
		// Ordenar cada lista de gaps
		for i := range g.gapValues {
			sortInts(g.gapValues[i])
		}

		// Encontrar rangos consecutivos para cada posición de gap
		allRanges := make([][][]int, len(g.gapValues))
		for i, gapList := range g.gapValues {
			allRanges[i] = findConsecutiveRanges(gapList)
		}

		// Generar todas las combinaciones de rangos
		combinations := generateRangeCombinations(allRanges)

		// Crear un patrón para cada combinación
		for _, combo := range combinations {
			consolidatedPattern := buildConsolidatedPattern(g.letters, combo)

			result[consolidatedPattern] = &PatternStat{
				Pattern:         consolidatedPattern,
				UppercaseCount:  countUppercaseInPattern(consolidatedPattern),
				SequenceIndices: g.seqIndices,
			}
		}
	}

	return result
}

// generateRangeCombinations genera todas las combinaciones de rangos
// Ejemplo: [[[1], [3]], [[11], [15,16]]] -> [[[1], [11]], [[1], [15,16]], [[3], [11]], [[3], [15,16]]]
func generateRangeCombinations(allRanges [][][]int) [][][]int {
	if len(allRanges) == 0 {
		return [][][]int{{}}
	}

	// Calcular el número total de combinaciones
	totalCombos := 1
	for _, ranges := range allRanges {
		if len(ranges) > 0 {
			totalCombos *= len(ranges)
		}
	}

	result := make([][][]int, 0, totalCombos)

	// Generar combinaciones recursivamente
	var generate func(pos int, current [][]int)
	generate = func(pos int, current [][]int) {
		if pos == len(allRanges) {
			// Hacer una copia de current
			combo := make([][]int, len(current))
			copy(combo, current)
			result = append(result, combo)
			return
		}

		ranges := allRanges[pos]
		if len(ranges) == 0 {
			// No hay rangos para esta posición, usar vacío
			generate(pos+1, append(current, []int{}))
		} else {
			for _, r := range ranges {
				generate(pos+1, append(current, r))
			}
		}
	}

	generate(0, [][]int{})
	return result
}
