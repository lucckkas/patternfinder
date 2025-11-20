package lcs_test

import (
	"fmt"
	"testing"

	"github.com/lucckkas/patternfinder/internal/lcs"
)

func TestDPTableParallelCorrectness(t *testing.T) {
	tests := []struct {
		name string
		sec1 string
		sec2 string
	}{
		{"Simple", "ABC", "AC"},
		{"Example from docs", "BABCBDABB", "DBDCABA"},
		{"Equal strings", "ABCD", "ABCD"},
		{"No common", "ABC", "DEF"},
		{"Empty sec2", "ABC", ""},
		{"Empty sec1", "", "ABC"},
		{"Both empty", "", ""},
		{"Longer sequences", "AGGTAB", "GXTXAYB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compute using sequential method
			dpSeq := lcs.DPTable(tt.sec1, tt.sec2)

			// Compute using parallel method
			dpPar := lcs.DPTableParallel(tt.sec1, tt.sec2)

			// Verify dimensions
			if len(dpSeq) != len(dpPar) {
				t.Fatalf("Different number of rows: seq=%d, par=%d", len(dpSeq), len(dpPar))
			}

			// Verify all values match
			for i := range dpSeq {
				if len(dpSeq[i]) != len(dpPar[i]) {
					t.Fatalf("Different number of cols at row %d: seq=%d, par=%d",
						i, len(dpSeq[i]), len(dpPar[i]))
				}
				for j := range dpSeq[i] {
					if dpSeq[i][j] != dpPar[i][j] {
						t.Errorf("Mismatch at dp[%d][%d]: seq=%d, par=%d",
							i, j, dpSeq[i][j], dpPar[i][j])
					}
				}
			}

			// Print the result for visual inspection
			if t.Failed() {
				fmt.Printf("\nSequential DP Table for %s vs %s:\n", tt.sec1, tt.sec2)
				lcs.PrintDP(tt.sec1, tt.sec2, dpSeq)
				fmt.Printf("\nParallel DP Table for %s vs %s:\n", tt.sec1, tt.sec2)
				lcs.PrintDP(tt.sec1, tt.sec2, dpPar)
			}
		})
	}
}

func TestDPTableParallelWithBacktracking(t *testing.T) {
	sec1 := "BABCBDABB"
	sec2 := "DBDCABA"

	// Compute DP tables
	dpSeq := lcs.DPTable(sec1, sec2)
	dpPar := lcs.DPTableParallel(sec1, sec2)

	// Get LCS using both DP tables
	lcsSeq := lcs.Backtracking(sec1, sec2, dpSeq)
	lcsPar := lcs.Backtracking(sec1, sec2, dpPar)

	// Convert to sets for comparison
	setSeq := make(map[string]bool)
	for _, s := range lcsSeq {
		setSeq[s] = true
	}

	setPar := make(map[string]bool)
	for _, s := range lcsPar {
		setPar[s] = true
	}

	// Verify same results
	if len(setSeq) != len(setPar) {
		t.Errorf("Different number of LCS: seq=%d, par=%d", len(setSeq), len(setPar))
	}

	for lcs := range setSeq {
		if !setPar[lcs] {
			t.Errorf("LCS %q found in sequential but not in parallel", lcs)
		}
	}

	for lcs := range setPar {
		if !setSeq[lcs] {
			t.Errorf("LCS %q found in parallel but not in sequential", lcs)
		}
	}

	t.Logf("Found %d LCS strings correctly", len(setSeq))
}
