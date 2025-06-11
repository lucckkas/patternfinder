package discovery

import (
	"strings"
	"unicode"
)

// normalizeTokens intercalza letra-gap-letra… e inserta gap=0 si falta.
// stripDigits elimina dígitos para comparar solo letras.
func normalizeTokens(raw []token) []token {
	// descartamos posibles bloques iniciales de número
	i := 0
	for i < len(raw) && raw[i].isNum {
		i++
	}
	raw = raw[i:]

	// contamos cuántas letras hay
	letCount := 0
	for _, t := range raw {
		if !t.isNum {
			letCount++
		}
	}
	if letCount == 0 {
		return nil
	}

	// construimos la lista normalizada: L N L N ... L (2*letCount-1 tokens)
	norm := make([]token, 0, 2*letCount-1)
	idxRaw := 0
	for li := 0; li < letCount; li++ {
		// letra
		if idxRaw < len(raw) && !raw[idxRaw].isNum {
			norm = append(norm, raw[idxRaw])
			idxRaw++
		} else {
			return nil
		}
		// espacio de número (salvo tras la última letra)
		if li < letCount-1 {
			if idxRaw < len(raw) && raw[idxRaw].isNum {
				norm = append(norm, raw[idxRaw])
				idxRaw++
			} else {
				// insertamos gap=0 si no había número
				norm = append(norm, token{isNum: true, num: 0})
			}
		}
	}
	return norm
}

func stripDigits(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if !unicode.IsDigit(r) {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
