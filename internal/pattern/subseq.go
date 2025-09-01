package pattern

// isSubsequence verifica si pat es subsecuencia de s.
func isSubsequence(pat, s string) bool {
	if len(pat) == 0 {
		return true
	}
	i := 0
	for j := 0; j < len(s) && i < len(pat); j++ {
		if s[j] == pat[i] {
			i++
		}
	}
	return i == len(pat)
}

// isSubsequenceAll: pat subsecuencia de TODOS los uppers.
func isSubsequenceAll(pat string, uppers []string) bool {
	for _, u := range uppers {
		if !isSubsequence(pat, u) {
			return false
		}
	}
	return true
}

// deleteOneAll genera todas las subsecuencias de pat eliminando 1 letra.
func deleteOneAll(pat string) []string {
	if len(pat) <= 1 {
		return nil
	}
	out := make([]string, 0, len(pat))
	for i := 0; i < len(pat); i++ {
		out = append(out, pat[:i]+pat[i+1:])
	}
	return out
}

// expandByDeletions toma un conjunto de candidatos y agrega subsecuencias
// eliminando hasta depth letras (p.ej. depth=2), manteniendo sólo las que
// son subsecuencia de TODAS las uppers.
func expandByDeletions(cands []string, uppers []string, depth int) []string {
	seen := map[string]struct{}{}
	queue := make([]string, 0, len(cands))
	for _, c := range cands {
		if _, ok := seen[c]; !ok {
			seen[c] = struct{}{}
			queue = append(queue, c)
		}
	}
	layer := cands
	for d := 0; d < depth; d++ {
		next := []string{}
		for _, c := range layer {
			for _, s := range deleteOneAll(c) {
				if _, ok := seen[s]; ok {
					continue
				}
				if len(s) == 0 {
					continue
				}
				if isSubsequenceAll(s, uppers) {
					seen[s] = struct{}{}
					next = append(next, s)
					queue = append(queue, s)
				}
			}
		}
		if len(next) == 0 {
			break
		}
		layer = next
	}
	return queue
}
