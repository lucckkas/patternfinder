package test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/lucckkas/patternfinder/internal/aggregate"
	"github.com/lucckkas/patternfinder/internal/gaps"
	"github.com/lucckkas/patternfinder/internal/lcs"
	"github.com/lucckkas/patternfinder/internal/utils"
)

type patternResult struct {
	pattern   string
	formatted string
}

type patternCase struct {
	name     string
	seq1     string
	seq2     string
	expected []patternResult
	wantErr  string
}

var testCases = []patternCase{
	{
		name: "NoPattern",
		seq1: "AB",
		seq2: "CD",
		expected: []patternResult{
			{pattern: "", formatted: ""},
		},
	},
	{
		name: "SingleLetterPattern",
		seq1: "AB",
		seq2: "CA",
		expected: []patternResult{
			{pattern: "A", formatted: "A"},
		},
	},
	{
		name: "SimplePatterns",
		seq1: "AB",
		seq2: "AB",
		expected: []patternResult{
			{pattern: "AB", formatted: "A-B"},
		},
	},
	{
		name: "ComplexPatterns",
		seq1: "DxxBxxxxDxxCxxxAxxBxxxA",
		seq2: "BxxAxxxBxxCxxxxBxxDxxxAxxBxxxxB",
		expected: []patternResult{
			{pattern: "BABA", formatted: "B-x(2|11)-A-x(2|3|11)-B-x(3|6|14)-A"},
			{pattern: "BCAB", formatted: "B-x(2|7|9)-C-x(3|11)-A-x(2|7)-B"},
			{pattern: "BCBA", formatted: "B-x(2|7|9)-C-x(4|6)-B-x(3|6)-A"},
			{pattern: "BDAB", formatted: "B-x(2|4|10|17)-D-x(3|6)-A-x(2|7)-B"},
		},
	},
}

// formatPatternResults ordena los patrones y construye la salida final con gaps.
func formatPatternResults(patterns []string, seqX, seqY string) ([]patternResult, error) {
	sort.Slice(patterns, func(i, j int) bool {
		if len(patterns[i]) != len(patterns[j]) {
			return len(patterns[i]) > len(patterns[j])
		}
		return patterns[i] < patterns[j]
	})

	results := make([]patternResult, 0, len(patterns))

	for _, pat := range patterns {
		gapsX, ok := gaps.AllGapValuesDistanceTotalViable(seqX, pat)
		if !ok {
			return nil, fmt.Errorf("no se pudieron calcular los gaps para %q en la secuencia 1", pat)
		}
		gapsY, ok := gaps.AllGapValuesDistanceTotalViable(seqY, pat)
		if !ok {
			return nil, fmt.Errorf("no se pudieron calcular los gaps para %q en la secuencia 2", pat)
		}

		union := aggregate.PairUnionSets(gapsX, gapsY)
		formatted := aggregate.FormatPatternWithValues(pat, union)
		results = append(results, patternResult{pattern: pat, formatted: formatted})
	}

	return results, nil
}

func TestPatternfinderScenarios(t *testing.T) {
	// modes ejecuta los casos con la variante secuencial y la paralela para comparar resultados.
	modes := []struct {
		name string
		run  func(seq1, seq2 string) ([]string, error)
	}{
		{
			name: "Sequential",
			run: func(seq1, seq2 string) ([]string, error) {
				upperX := utils.UpperOnly(seq1)
				upperY := utils.UpperOnly(seq2)
				matriz := lcs.DPTable(upperX, upperY)
				return lcs.Backtracking(upperX, upperY, matriz), nil
			},
		},
		{
			name: "Parallel",
			run: func(seq1, seq2 string) ([]string, error) {
				upperX := utils.UpperOnly(seq1)
				upperY := utils.UpperOnly(seq2)
				matriz := lcs.DPTableParallel(upperX, upperY)
				return lcs.BacktrackingParallel(upperX, upperY, matriz), nil
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		for _, mode := range modes {
			mode := mode
			t.Run(tc.name+"/"+mode.name, func(t *testing.T) {
				patterns, err := mode.run(tc.seq1, tc.seq2)
				if err != nil {
					t.Fatalf("error inesperado en modo %s: %v", mode.name, err)
				}
				got, err := formatPatternResults(patterns, tc.seq1, tc.seq2)
				if tc.wantErr != "" {
					if err == nil {
						t.Fatalf("se esperaba error que contuviera %q", tc.wantErr)
					}
					if !strings.Contains(err.Error(), tc.wantErr) {
						t.Fatalf("error esperado que contuviera %q, se obtuvo %q", tc.wantErr, err)
					}
					return
				}

				if err != nil {
					t.Fatalf("no se esperaba error: %v", err)
				}

				if len(got) != len(tc.expected) {
					t.Fatalf("se esperaban %d patrones, se obtuvieron %d", len(tc.expected), len(got))
				}

				for i := range tc.expected {
					if got[i] != tc.expected[i] {
						t.Fatalf("[%d] esperado %+v, obtenido %+v", i, tc.expected[i], got[i])
					}
				}
			})
		}
	}
}
