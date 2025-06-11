package discovery

// CompareSubsequences devuelve los patrones únicos de a vs. b,
// soportando rangos multi-dígito y varias letras.
func CompareSubsequences(a, b []string) []string {
	setB := make(map[string]struct{}, len(b))
	for _, subsequenceB := range b {
		setB[subsequenceB] = struct{}{}
	}

	result := make(map[string]struct{})
	for _, subsequenceA := range a {
		if _, ok := setB[subsequenceA]; ok { // si s1 está en b
			if pat := formatIdentical(subsequenceA); pat != "" {
				result[pat] = struct{}{}
			}
		} else {
			for _, s2 := range b { // buscamos variantes
				if stripDigits(subsequenceA) == stripDigits(s2) {
					if patron := formatVariant(subsequenceA, s2); patron != "" {
						result[patron] = struct{}{}
					}
				}
			}
		}
	}

	out := make([]string, 0, len(result))
	for p := range result {
		out = append(out, p)
	}
	return out
}
