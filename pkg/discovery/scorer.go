package discovery

import (
	"strconv"
	"unicode"
)

// ScorePatterns asigna puntuaciones a cada patrón según tu fórmula.
func ScorePatterns(patterns []string, maxLen int) map[string]int {
	scores := make(map[string]int, len(patterns))
	for _, p := range patterns {
		scores[p] = scoreSingle(p, maxLen)
	}
	return scores
}

// scoreSingle parsea p buscando x(n) y letras, y aplica tu fórmula.
func scoreSingle(p string, maxLen int) int {
	score := 0
	numBuf := ""

	flushNum := func() {
		if numBuf != "" {
			if n, err := strconv.Atoi(numBuf); err == nil {
				// sumar el mayor entre n y maxLen (o tu regla)
				if n > maxLen {
					score += n
				} else {
					score += maxLen
				}
			}
			numBuf = ""
		}
	}

	for _, r := range p {
		switch {
		case unicode.IsDigit(r):
			numBuf += string(r)
		case r == 'x' || r == '(' || r == ')' || r == ',':
			// delimitadores de rango o placeholder
		case unicode.IsUpper(r):
			flushNum()
			score += maxLen * maxLen
		case unicode.IsLower(r):
			flushNum()
			score += maxLen
		default:
			flushNum()
		}
	}
	flushNum()
	return score
}
