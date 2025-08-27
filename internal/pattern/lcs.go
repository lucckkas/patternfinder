package pattern

// LCSClassic devuelve UNA LCS entre a y b (mayúsculas). Complejidad O(n*m).
// Elegimos implementación clara y determinista; sobre strings solo-mayúsculas suele ser suficiente.
func LCSClassic(a, b string) string {
	n, m := len(a), len(b)
	if n == 0 || m == 0 {
		return ""
	}
	// dp[i][j] = LCS length de a[i:] y b[j:]
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if a[i] == b[j] {
				dp[i][j] = 1 + dp[i+1][j+1]
			} else {
				if dp[i+1][j] >= dp[i][j+1] {
					dp[i][j] = dp[i+1][j]
				} else {
					dp[i][j] = dp[i][j+1]
				}
			}
		}
	}
	// backtracking
	out := make([]byte, 0, dp[0][0])
	i, j := 0, 0
	for i < n && j < m {
		if a[i] == b[j] {
			out = append(out, a[i])
			i++
			j++
		} else if dp[i+1][j] >= dp[i][j+1] {
			i++
		} else {
			j++
		}
	}
	return string(out)
}

// ProgressiveLCSUpper aplica LCS de manera progresiva sobre una lista de secuencias de mayúsculas.
// Estrategia: comenzar por la más corta y cruzar con el resto en orden creciente de longitud.
func ProgressiveLCSUpper(uppers []string) string {
	if len(uppers) == 0 {
		return ""
	}
	// ordenar índices por longitud
	idx := make([]int, len(uppers))
	for i := range uppers {
		idx[i] = i
	}
	// simple insertion sort por claridad (listas cortas). Puedes cambiar a sort.Slice si prefieres.
	for i := 1; i < len(idx); i++ {
		j := i
		for j > 0 && len(uppers[idx[j]]) < len(uppers[idx[j-1]]) {
			idx[j], idx[j-1] = idx[j-1], idx[j]
			j--
		}
	}
	consensus := uppers[idx[0]]
	for t := 1; t < len(idx); t++ {
		consensus = LCSClassic(consensus, uppers[idx[t]])
		if consensus == "" {
			break
		}
	}
	return consensus
}
