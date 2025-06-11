package discovery

import (
	"strconv"
	"strings"
	"unicode"
)

// ScorePatterns asigna a cada patrón su puntuación según:
//
//	Score = X + m*n + y*n^2
//	donde:
//	  X  = suma de los gaps (para x(i,j) toma el mayor: max(i,j))
//	  m  = nº de letras minúsculas
//	  y  = nº de letras mayúsculas
//	  n  = longitud original (maxLen)
func ScorePatterns(patterns []string, maxLen int) map[string]int {
	scores := make(map[string]int, len(patterns))
	for _, p := range patterns {
		scores[p] = scoreSingle(p, maxLen)
	}
	return scores
}

func scoreSingle(p string, maxLen int) int {
	parts := strings.Split(p, "-")
	sumX := 0
	lowerCount := 0
	upperCount := 0

	for _, tok := range parts {
		// Gaps: "x(n)" o "x(i,j)"
		if strings.HasPrefix(tok, "x(") && strings.HasSuffix(tok, ")") {
			inner := tok[2 : len(tok)-1] // entre paréntesis
			if comma := strings.Index(inner, ","); comma >= 0 {
				// rango: tomamos el mayor
				n1, _ := strconv.Atoi(inner[:comma])
				n2, _ := strconv.Atoi(inner[comma+1:])
				if n1 > n2 {
					sumX += n1
				} else {
					sumX += n2
				}
			} else {
				// único valor
				n, _ := strconv.Atoi(inner)
				sumX += n
			}
			continue
		}

		// Letras
		if len(tok) == 1 {
			r := rune(tok[0])
			switch {
			case unicode.IsUpper(r):
				upperCount++
			case unicode.IsLower(r):
				lowerCount++
			}
		}
	}

	// aplicamos la fórmula
	return sumX + lowerCount*maxLen + upperCount*(maxLen*maxLen)
}
