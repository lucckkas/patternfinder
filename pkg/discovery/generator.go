package discovery

import (
	"math/big"
	"strconv"
	"strings"
)

// GenerateSubsequences genera todas las subsecuencias ProSite-like de seq.
// Omite aquellas que no tengan al menos una letra.
func GenerateSubsequences(seq string) []string {
	n := len(seq)
	total := new(big.Int).Lsh(big.NewInt(1), uint(n))
	var out []string
	tmp := make([]rune, n)

	for mask := big.NewInt(0); mask.Cmp(total) < 0; mask.Add(mask, big.NewInt(1)) {
		for i, r := range seq {
			// bit más significativo primero
			if mask.Bit(n-1-i) == 1 {
				tmp[i] = 'x'
			} else {
				tmp[i] = r
			}
		}
		collapsed := collapseX(tmp)
		if containsLetter(collapsed) {
			out = append(out, collapsed)
		}
	}
	return out
}

// collapseX convierte ['x','x','A','x','B'] → "2A1B"
func collapseX(runes []rune) string {
	var sb strings.Builder
	count := 0
	flush := func() {
		if count > 0 {
			sb.WriteString(strconv.Itoa(count))
			count = 0
		}
	}

	for _, r := range runes {
		if r == 'x' {
			count++
		} else {
			flush()
			sb.WriteRune(r)
		}
	}
	flush()
	return sb.String()
}

// containsLetter devuelve true si s contiene al menos una A–Z.
func containsLetter(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}
