package discovery

// CompareSubsequences devuelve los patrones únicos de a vs. b,
// soportando rangos multi-dígito y varias letras.
func CompareSubsequences(a, b []string) []string {
	setB := make(map[string]struct{}, len(b))
	for _, s := range b {
		setB[s] = struct{}{}
	}

	result := make(map[string]struct{})
	for _, s1 := range a {
		if _, ok := setB[s1]; ok {
			if pat := formatIdentical(s1); pat != "" {
				result[pat] = struct{}{}
			}
		} else {
			for _, s2 := range b {
				if stripDigits(s1) == stripDigits(s2) {
					if pat := formatVariant(s1, s2); pat != "" {
						result[pat] = struct{}{}
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
