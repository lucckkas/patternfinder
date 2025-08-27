package pattern

// BestCommonPattern ejecuta el pipeline completo:
// - Proyección a mayúsculas
// - LCS progresivo sobre mayúsculas
// - Agregación de gaps y selección final (único consenso en esta versión)
func BestCommonPattern(seqs []string) (AggregatedPattern, bool) {
	if len(seqs) == 0 {
		return AggregatedPattern{}, false
	}
	projs := make([]UpperProjection, 0, len(seqs))
	uppers := make([]string, 0, len(seqs))
	for _, s := range seqs {
		p := ProjectUpper(s)
		projs = append(projs, p)
		uppers = append(uppers, p.UpperSeq)
	}
	consensus := ProgressiveLCSUpper(uppers)
	if consensus == "" {
		return AggregatedPattern{}, false
	}
	agg, ok := AggregateOverSequences(consensus, projs)
	if !ok {
		return AggregatedPattern{}, false
	}
	return agg, true
}
