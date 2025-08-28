package pattern

func AggregateOverSequences(pattern string, projs []UpperProjection) (AggregatedPattern, bool) {
	perSeqMins := make([][]int, 0, len(projs))
	for _, p := range projs {
		mins, ok := minGapsForPatternPerSequence(pattern, p)
		if !ok {
			return AggregatedPattern{}, false
		}
		perSeqMins = append(perSeqMins, mins)
	}

	L := len(pattern)
	ranges := make([]IntRange, 0, max(0, L-1))
	avgs := make([]float64, 0, max(0, L-1))

	for gapIdx := 0; gapIdx+1 <= L-1 && L >= 2; gapIdx++ {
		minv := int(^uint(0) >> 1)
		maxv := -minv - 1
		sum := 0
		for _, mins := range perSeqMins {
			val := mins[gapIdx]
			if val < minv {
				minv = val
			}
			if val > maxv {
				maxv = val
			}
			sum += val
		}
		ranges = append(ranges, IntRange{Min: minv, Max: maxv})
		avgs = append(avgs, float64(sum)/float64(len(perSeqMins)))
	}

	scoreGaps := 0
	for _, r := range ranges {
		scoreGaps += r.Min // conservador
	}

	return AggregatedPattern{
		Pattern:      pattern,
		GapRanges:    ranges,
		GapAvg:       avgs,
		ScoreUpper:   len(pattern),
		ScoreGapsSum: scoreGaps,
	}, true
}
