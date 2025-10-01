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
func AllLCS(sec1, sec2 string, matriz [][]int) []string {
	type key struct{ i, j int }
	memo := make(map[key]map[string]struct{}) // evita recalcular estados

	type frame struct {
		i, j int
		done bool
	}

	stack := []frame{{i: len(sec1), j: len(sec2)}}

	for len(stack) > 0 {
		fmt.Println("Stack size:", len(stack))
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		k := key{current.i, current.j}
		fmt.Println("Processing:", k)

		if _, ok := memo[k]; ok {
			fmt.Println("  Already computed, skipping")
			continue
		}

		if matriz[current.i][current.j] == 0 {
			memo[k] = map[string]struct{}{"": {}}
			continue
		}

		if !current.done {
			// reinsertar el estado indicando que ya se ha procesado
			stack = append(stack, frame{i: current.i, j: current.j, done: true})

			// si hay coincidencia, ir diagonal
			if sec1[current.i-1] == sec2[current.j-1] {
				stack = append(stack, frame{i: current.i - 1, j: current.j - 1})
				continue
			}

			// si no hay coincidencia, ir a los estados que mantienen la longitud
			if matriz[current.i-1][current.j] == matriz[current.i][current.j] {
				stack = append(stack, frame{i: current.i - 1, j: current.j})
			}
			if matriz[current.i][current.j-1] == matriz[current.i][current.j] {
				stack = append(stack, frame{i: current.i, j: current.j - 1})
			}
			continue
		}

		res := make(map[string]struct{})

		if sec1[current.i-1] == sec2[current.j-1] {
			for s := range memo[key{current.i - 1, current.j - 1}] {
				res[s+string(sec1[current.i-1])] = struct{}{}
			}
		} else {
			if matriz[current.i-1][current.j] == matriz[current.i][current.j] {
				for s := range memo[key{current.i - 1, current.j}] {
					res[s] = struct{}{}
				}
			}
			if matriz[current.i][current.j-1] == matriz[current.i][current.j] {
				for s := range memo[key{current.i, current.j - 1}] {
					res[s] = struct{}{}
				}
			}
		}

		memo[k] = res
	}

	set := memo[key{len(sec1), len(sec2)}]
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
