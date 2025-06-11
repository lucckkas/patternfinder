package discovery

import (
	"strconv"
	"strings"
)

// CompareSubsequences toma dos listas de subsecuencias y devuelve patrones únicos.
func CompareSubsequences(a, b []string) []string {
	setB := make(map[string]struct{}, len(b))
	for _, s := range b {
		setB[s] = struct{}{}
	}

	patterns := make(map[string]struct{})
	for _, s1 := range a {
		if _, ok := setB[s1]; ok {
			if pat := formatIdentical(s1); pat != "" {
				patterns[pat] = struct{}{}
			}
		} else {
			for _, s2 := range b {
				if stripDigits(s1) == stripDigits(s2) {
					if pat := formatVariant(s1, s2); pat != "" {
						patterns[pat] = struct{}{}
					}
				}
			}
		}
	}

	out := make([]string, 0, len(patterns))
	for p := range patterns {
		out = append(out, p)
	}
	return out
}

// stripDigits elimina dígitos para comparar solo letras.
func stripDigits(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if r < '0' || r > '9' {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

// formatIdentical formatea una subsecuencia idéntica al estilo ProSite.
func formatIdentical(seq string) string {
	var sb strings.Builder
	first := true
	for i := 0; i < len(seq); i++ {
		ch := seq[i]
		switch {
		case ch >= 'A' && ch <= 'Z':
			if !first {
				sb.WriteString("-")
			}
			sb.WriteByte(ch)
			first = false
		case ch >= '0' && ch <= '9':
			sb.WriteString("-x(")
			sb.WriteByte(ch)
			sb.WriteString(")")
			first = false
		}
	}
	return sb.String()
}

// formatVariant compara dos secuencias y genera el patrón con rangos.
func formatVariant(s1, s2 string) string {
	var sb strings.Builder
	first := true
	countX := 0
	length := min(len(s1), len(s2))

	flushX := func() {
		if countX > 0 {
			if !first {
				sb.WriteString("-")
			}
			sb.WriteString("x(")
			sb.WriteString(strconv.Itoa(countX))
			sb.WriteString(")")
			countX = 0
			first = false
		}
	}

	for i := 0; i < length; i++ {
		ch1, ch2 := s1[i], s2[i]
		switch {
		case ch1 >= 'A' && ch1 <= 'Z' && ch1 == ch2:
			flushX()
			if !first {
				sb.WriteString("-")
			}
			sb.WriteByte(ch1)
			first = false

		case isDigit(ch1) && isDigit(ch2):
			num1 := int(ch1 - '0')
			num2 := int(ch2 - '0')
			low, high := min(num1, num2), max(num1, num2)
			flushX()
			if !first {
				sb.WriteString("-")
			}
			sb.WriteString("x(")
			sb.WriteString(strconv.Itoa(low))
			sb.WriteString(",")
			sb.WriteString(strconv.Itoa(high))
			sb.WriteString(")")
			first = false

		default:
			countX++
		}
	}
	flushX()
	return sb.String()
}

func isDigit(ch byte) bool { return ch >= '0' && ch <= '9' }
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
