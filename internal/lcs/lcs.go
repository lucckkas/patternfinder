package lcs

import "fmt"

// DPTable construye la tabla para las longitudes de LCS
func DPTable(sec1, sec2 string) [][]int {
	n, m := len(sec1), len(sec2)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if sec1[i-1] == sec2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1 // coincidencia: diagonal + 1
			} else {
				if dp[i-1][j] >= dp[i][j-1] { // tomar máximo
					dp[i][j] = dp[i-1][j]
				} else {
					dp[i][j] = dp[i][j-1]
				}
			}
		}
	}
	return dp
}

// AllLCS hace el backtracking para encontrar todas las LCS posibles
func AllLCS(sec1, sec2 string, dp [][]int) []string {
	type key struct{ i, j int }
	memo := make(map[key]map[string]struct{}) // para simplificar calculos repetidos

	var dfs func(i, j int) map[string]struct{}
	dfs = func(i, j int) map[string]struct{} {
		if dp[i][j] == 0 {
			return map[string]struct{}{"": {}}
		}
		k := key{i, j}
		if got, ok := memo[k]; ok {
			return got
		}
		res := make(map[string]struct{})

		// coincidencia: ir diagonal
		if i > 0 && j > 0 && sec1[i-1] == sec2[j-1] && dp[i][j] == dp[i-1][j-1]+1 {
			for s := range dfs(i-1, j-1) {
				res[s+string(sec1[i-1])] = struct{}{}
			}
			memo[k] = res
			return res
		}
		// no coincidencia: ir arriba o izquierda según convenga
		if i > 0 && dp[i-1][j] == dp[i][j] {
			for s := range dfs(i-1, j) {
				res[s] = struct{}{}
			}
		}
		if j > 0 && dp[i][j-1] == dp[i][j] {
			for s := range dfs(i, j-1) {
				res[s] = struct{}{}
			}
		}
		memo[k] = res
		return res
	}

	set := dfs(len(sec1), len(sec2))
	out := make([]string, 0, len(set))
	for s := range set {
		out = append(out, s)
	}
	return out
}

func PrintDP(sec1, sec2 string, dp [][]int) {
	fmt.Printf("     ")
	for j := 0; j < len(sec2); j++ {
		fmt.Printf("  %c", sec2[j])
	}
	fmt.Println()
	for i := 0; i <= len(sec1); i++ {
		if i == 0 {
			fmt.Printf("  ")
		} else {
			fmt.Printf("%c ", sec1[i-1])
		}
		for j := 0; j <= len(sec2); j++ {
			fmt.Printf("%2d", dp[i][j])
			if j < len(sec2) {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}
