package lcs

import (
	"testing"

	"github.com/lucckkas/patternfinder/internal/gaps"
)

func TestAllGapsDebug(t *testing.T) {
	seq := "ABACBACCBACCCBACCCBABASBBBABAAABABABAA"
	pattern := "AAAAAAABBBBBBB"

	t.Logf("Testing pattern=%s in seq=%s", pattern, seq)

	// Contar caracteres
	countA := 0
	countB := 0
	for _, ch := range seq {
		if ch == 'A' {
			countA++
		} else if ch == 'B' {
			countB++
		}
	}
	t.Logf("Sequence has %d 'A's and %d 'B's", countA, countB)
	t.Logf("Pattern needs 7 'A's and 7 'B's")

	sets, ok := gaps.AllGapValuesDistanceTotalViable(seq, pattern)
	if !ok {
		t.Logf("Failed to calculate gaps - this is the bug we need to fix")
		return
	}

	t.Logf("Success! Found %d gap sets", len(sets))
	for i, s := range sets {
		t.Logf("Gap set %d has %d values", i, len(s))
	}
}
