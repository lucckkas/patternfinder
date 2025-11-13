package lcs_test

import (
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/lucckkas/patternfinder/internal/aggregate"
	// "github.com/lucckkas/patternfinder/internal/gaps"
	"github.com/lucckkas/patternfinder/internal/lcs"
	"github.com/lucckkas/patternfinder/internal/utils"
)

// TestSequentialHotspots perfila cada etapa del algoritmo secuencial
// para evidenciar qué secciones concentran la mayor carga computacional.
func TestSequentialHotspots(t *testing.T) {
	seq1Raw := strings.Repeat("acbdefgHIKLMNPQRSTVWY", 400) + "ACDEFGHIKLMNPQRSTVWY"
	seq2Raw := strings.Repeat("ghIKLMNPQRSTVWYacbdef", 200) + "ACDFGHIKLMNPQRSTVWY"

	seq1 := utils.UpperOnly(seq1Raw)
	seq2 := utils.UpperOnly(seq2Raw)
	if len(seq1) == 0 || len(seq2) == 0 {
		t.Fatalf("las secuencias de prueba no contienen mayúsculas suficientes")
	}

	type measurement struct {
		name string
		dur  time.Duration
	}
	var (
		measures []measurement
		dp       [][]int
		lcsList  []string
		bestLCS  string
		setsX    []map[int]struct{}
		setsY    []map[int]struct{}
		union    []aggregate.GapValues
	)

	record := func(name string, fn func()) {
		start := time.Now()
		fn()
		measures = append(measures, measurement{name: name, dur: time.Since(start)})
	}

	record("dp_table", func() {
		dp = lcs.DPTable(seq1, seq2)
	})

	record("backtracking", func() {
		lcsList = lcs.Backtracking(seq1, seq2, dp)
	})
	if len(lcsList) == 0 {
		t.Fatalf("no se generaron LCS para las secuencias de prueba")
	}
	sort.Slice(lcsList, func(i, j int) bool {
		if len(lcsList[i]) != len(lcsList[j]) {
			return len(lcsList[i]) > len(lcsList[j])
		}
		return lcsList[i] < lcsList[j]
	})
	bestLCS = lcsList[0]

	// record("gaps_seq1", func() {
	// 	var ok bool
	// 	setsX, ok = gaps.AllGapValuesDistanceTotalViable(seq1Raw, bestLCS)
	// 	if !ok {
	// 		t.Fatalf("no se pudieron obtener gaps para seq1")
	// 	}
	// })

	// record("gaps_seq2", func() {
	// 	var ok bool
	// 	setsY, ok = gaps.AllGapValuesDistanceTotalViable(seq2Raw, bestLCS)
	// 	if !ok {
	// 		t.Fatalf("no se pudieron obtener gaps para seq2")
	// 	}
	// })

	record("gap_union", func() {
		union = aggregate.PairUnionSets(setsX, setsY)
	})

	record("pattern_format", func() {
		_ = aggregate.FormatPatternWithValues(bestLCS, union)
	})

	var total time.Duration
	for _, m := range measures {
		total += m.dur
	}
	sort.Slice(measures, func(i, j int) bool {
		return measures[i].dur > measures[j].dur
	})
	t.Logf("Carga secuencial para |seq1|=%d y |seq2|=%d (total %s):", len(seq1), len(seq2), total)
	for _, m := range measures {
		pct := 0.0
		if total > 0 {
			pct = float64(m.dur) / float64(total) * 100
		}
		t.Logf(" - %s: %s (%.1f%%)", m.name, m.dur, pct)
	}
	if len(measures) > 0 {
		t.Logf("Sección más costosa: %s", measures[0].name)
	}
}
