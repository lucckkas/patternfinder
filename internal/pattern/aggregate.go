package pattern

// AggregateOverSequences computa rangos y promedios de gaps para un patrón dado
// sobre múltiples secuencias proyectadas. Si alguna no soporta el patrón, retorna false.
func AggregateOverSequences(pattern string, projs []UpperProjection) (AggregatedPattern, bool) {
	allGaps := make([][]int, 0, len(projs))
	for _, p := range projs {
		g, ok := gapsForPattern(pattern, p)
		if !ok {
			return AggregatedPattern{}, false
		}
		allGaps = append(allGaps, g)
	}
	L := len(pattern)
	ranges := make([]IntRange, 0, max(0, L-1))
	avgs := make([]float64, 0, max(0, L-1))
	for gapIdx := 0; gapIdx+1 <= L-1 && L >= 2; gapIdx++ {
		minv := int(^uint(0) >> 1) // max int
		maxv := -minv - 1          // min int
		sum := 0
		for _, g := range allGaps {
			val := g[gapIdx]
			if val < minv {
				minv = val
			}
			if val > maxv {
				maxv = val
			}
			sum += val
		}
		ranges = append(ranges, IntRange{Min: minv, Max: maxv})
		avgs = append(avgs, float64(sum)/float64(len(allGaps)))
	}
	scoreGaps := 0
	for _, r := range ranges {
		scoreGaps += r.Min // conservador: suma de mínimos
	}
	return AggregatedPattern{
		Pattern:      pattern,
		GapRanges:    ranges,
		GapAvg:       avgs,
		ScoreUpper:   len(pattern),
		ScoreGapsSum: scoreGaps,
	}, true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
