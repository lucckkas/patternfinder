package discovery

import (
	"fmt"
	"strings"
)

// formatIdentical: a tokens normalizados les aplica “x(n)” y letras.
func formatIdentical(seq string) string {
	raw := parseCollapsed(seq)
	toks := normalizeTokens(raw)
	if toks == nil {
		return ""
	}
	parts := make([]string, 0, len(toks))
	for _, t := range toks {
		if t.isNum {
			if t.num > 0 {
				parts = append(parts, fmt.Sprintf("x(%d)", t.num))
			}
		} else {
			parts = append(parts, string(t.letter))
		}
	}
	return strings.Join(parts, "-")
}

// formatVariant: dado seq1 y seq2 genera rangos “x(min,max)”.
func formatVariant(s1, s2 string) string {
	raw1 := parseCollapsed(s1)
	raw2 := parseCollapsed(s2)
	t1 := normalizeTokens(raw1)
	t2 := normalizeTokens(raw2)
	if t1 == nil || t2 == nil || len(t1) != len(t2) {
		return ""
	}

	parts := make([]string, 0, len(t1))
	for i := 0; i < len(t1); i++ {
		a, b := t1[i], t2[i]
		// misma letra
		if !a.isNum && !b.isNum && a.letter == b.letter {
			parts = append(parts, string(a.letter))

			// gap numérico en ambas
		} else if a.isNum && b.isNum {
			lo, hi := a.num, b.num
			if lo > hi {
				lo, hi = hi, lo
			}
			if lo == hi {
				parts = append(parts, fmt.Sprintf("x(%d)", lo))
			} else {
				parts = append(parts, fmt.Sprintf("x(%d,%d)", lo, hi))
			}
		} else {
			// caso inesperado
			return ""
		}
	}

	return strings.Join(parts, "-")
}
