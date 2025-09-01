package pattern

import (
	"math/rand"
	"sort"
	"time"
)

type TopKOptions struct {
	NumOrders   int
	PerPairAlt  int
	BeamWidth   int
	DeleteDepth int // NUEVO: cuántas rondas de borrado explorar (1 o 2 recomendado)
	RandomSeed  int64
}

func DefaultTopKOptions() TopKOptions {
	return TopKOptions{
		NumOrders:   12,
		PerPairAlt:  3,
		BeamWidth:   20,
		DeleteDepth: 2, // probar subsecuencias quitando 1 y 2 letras
		RandomSeed:  time.Now().UnixNano(),
	}
}

func TopKCommonPatterns(seqs []string, K int, opt TopKOptions) ([]AggregatedPattern, bool) {
	if len(seqs) == 0 || K <= 0 {
		return nil, false
	}
	// Proyección
	projs := make([]UpperProjection, 0, len(seqs))
	uppers := make([]string, 0, len(seqs))
	for _, s := range seqs {
		p := ProjectUpper(s)
		projs = append(projs, p)
		uppers = append(uppers, p.UpperSeq)
	}

	// Órdenes
	baseIdx := make([]int, len(uppers))
	for i := range baseIdx {
		baseIdx[i] = i
	}
	sort.Slice(baseIdx, func(i, j int) bool {
		return len(uppers[baseIdx[i]]) < len(uppers[baseIdx[j]])
	})

	orders := make([][]int, 0, opt.NumOrders)
	orders = append(orders, append([]int(nil), baseIdx...))

	rng := rand.New(rand.NewSource(opt.RandomSeed))
	for t := 1; t < opt.NumOrders; t++ {
		perm := append([]int(nil), baseIdx...)
		for i := len(perm) - 1; i > 1; i-- {
			j := rng.Intn(i)
			if j == 0 {
				j = 1
			}
			perm[i], perm[j] = perm[j], perm[i]
		}
		orders = append(orders, perm)
	}

	// Beam de consensos
	type beamItem struct{ pat string }
	seenCand := map[string]struct{}{}
	candidates := make([]string, 0, 256)

	for _, ord := range orders {
		beam := []beamItem{{pat: uppers[ord[0]]}}
		for t := 1; t < len(ord); t++ {
			nextU := uppers[ord[t]]
			newBeam := make([]beamItem, 0, opt.BeamWidth*opt.PerPairAlt)
			for _, it := range beam {
				alts := LCSKBest(it.pat, nextU, opt.PerPairAlt)
				for _, c := range alts {
					if c == "" {
						continue
					}
					newBeam = append(newBeam, beamItem{pat: c})
				}
			}
			sort.Slice(newBeam, func(i, j int) bool {
				if len(newBeam[i].pat) != len(newBeam[j].pat) {
					return len(newBeam[i].pat) > len(newBeam[j].pat)
				}
				return newBeam[i].pat < newBeam[j].pat
			})
			if len(newBeam) > opt.BeamWidth {
				newBeam = newBeam[:opt.BeamWidth]
			}
			beam = newBeam
			if len(beam) == 0 {
				break
			}
		}
		for _, it := range beam {
			if _, ok := seenCand[it.pat]; !ok && len(it.pat) > 0 {
				seenCand[it.pat] = struct{}{}
				candidates = append(candidates, it.pat)
			}
		}
	}

	// NUEVO: expandir por borrado para obtener subsecuencias más cortas comunes
	if opt.DeleteDepth > 0 {
		candidates = expandByDeletions(candidates, uppers, opt.DeleteDepth)
	}

	// Agregar (min-min) y ordenar
	aggs := make([]AggregatedPattern, 0, len(candidates))
	seenAgg := map[string]struct{}{}
	for _, pat := range candidates {
		if _, ok := seenAgg[pat]; ok {
			continue
		}
		agg, ok := AggregateOverSequences(pat, projs)
		if ok {
			seenAgg[pat] = struct{}{}
			aggs = append(aggs, agg)
		}
	}
	if len(aggs) == 0 {
		return nil, false
	}

	sort.Slice(aggs, func(i, j int) bool {
		if aggs[i].ScoreUpper != aggs[j].ScoreUpper {
			return aggs[i].ScoreUpper > aggs[j].ScoreUpper
		}
		if aggs[i].ScoreGapsSum != aggs[j].ScoreGapsSum {
			return aggs[i].ScoreGapsSum > aggs[j].ScoreGapsSum
		}
		return aggs[i].Pattern < aggs[j].Pattern
	})
	if len(aggs) > K {
		aggs = aggs[:K]
	}
	return aggs, true
}
