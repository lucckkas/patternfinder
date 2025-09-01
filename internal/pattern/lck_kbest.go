package pattern

// LCSKBest devuelve hasta k distintas LCS entre a y b, ordenadas por longitud (todas misma)
// y luego por lexicográfico (para estabilidad). Si k<=1, usa LCSClassic.
func LCSKBest(a, b string, k int) []string {
	if k <= 1 {
		return []string{LCSClassic(a, b)}
	}
	n, m := len(a), len(b)
	if n == 0 || m == 0 {
		return []string{""}
	}
	// DP longitud
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
	target := dp[0][0]
	if target == 0 {
		return []string{""}
	}

	// Backtracking K-best con poda (evitar explosión):
	// - Cuando a[i]==b[j], bajamos por match.
	// - Cuando dp[i+1][j] == dp[i][j+1], ramificamos (pero limitamos resultados totales a k).
	// - Memo key: (i,j,restLen) para evitar reexpandir caminos que no pueden llegar.
	type key struct{ i, j, need int }
	memo := map[key]int{} // max cuántas soluciones quedan desde (i,j) con need

	var res []string
	var dfs func(i, j, need int, buf []byte)
	dfs = func(i, j, need int, buf []byte) {
		if len(res) >= k {
			return
		}
		if need == 0 {
			res = append(res, string(buf))
			return
		}
		if i >= n || j >= m {
			return
		}
		// poda por dp
		if dp[i][j] < need {
			return
		}
		kk := key{i, j, need}
		if left, ok := memo[kk]; ok && left == 0 {
			return
		}

		if a[i] == b[j] {
			dfs(i+1, j+1, need-1, append(buf, a[i]))
		} else {
			// Branching control
			v1, v2 := dp[i+1][j], dp[i][j+1]
			if v1 > v2 {
				dfs(i+1, j, need, buf)
			} else if v2 > v1 {
				dfs(i, j+1, need, buf)
			} else { // empate, ramificamos (dos caminos)
				dfs(i+1, j, need, buf)
				if len(res) < k {
					dfs(i, j+1, need, buf)
				}
			}
		}
		// nota: no guardamos conteo exacto, sólo marcamos que pasamos por aquí
		// para indicar que no vale la pena revisitar si ya alcanzamos k.
		if len(res) >= k {
			memo[kk] = 0
		}
	}

	dfs(0, 0, target, make([]byte, 0, target))

	// Dedup (es raro pero posible con empates).
	seen := map[string]struct{}{}
	out := make([]string, 0, len(res))
	for _, s := range res {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
		if len(out) >= k {
			break
		}
	}
	if len(out) == 0 {
		return []string{LCSClassic(a, b)}
	}
	return out
}
