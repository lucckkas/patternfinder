package lcs

import (
	"fmt"
	"sync"
)

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

// DPTableParallel construye la tabla LCS evaluando cada diagonal en paralelo.
func DPTableParallel(sec1, sec2 string) [][]int {
	n, m := len(sec1), len(sec2)
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, m+1)
	}

	for diag := 2; diag <= n+m; diag++ {
		start := 1
		if diag-m > start {
			start = diag - m
		}
		if start < 1 {
			start = 1
		}

		end := diag - 1
		if end > n {
			end = n
		}

		var wg sync.WaitGroup
		for i := start; i <= end; i++ {
			j := diag - i
			if j < 1 || j > m {
				continue
			}

			wg.Add(1)
			go func(i, j int) {
				defer wg.Done()
				if sec1[i-1] == sec2[j-1] {
					dp[i][j] = dp[i-1][j-1] + 1
					return
				}
				if dp[i-1][j] >= dp[i][j-1] {
					dp[i][j] = dp[i-1][j]
					return
				}
				dp[i][j] = dp[i][j-1]
			}(i, j)
		}
		wg.Wait()
	}

	return dp
}

// Backtracking hace el backtracking para encontrar todas las LCS posibles
func Backtracking(sec1, sec2 string, matriz [][]int) []string {
	type key struct{ i, j int }
	memo := make(map[key]map[string]struct{}) // evita recalcular estados

	type frame struct {
		i, j int
		done bool
	}

	stack := []frame{{i: len(sec1), j: len(sec2)}}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		k := key{current.i, current.j}

		if _, ok := memo[k]; ok {
			continue
		}

		if matriz[current.i][current.j] == 0 {
			memo[k] = map[string]struct{}{"": {}}
			continue
		}

		if !current.done {
			// reinsertar el estado indicando que los hijos ya fueron explorados
			stack = append(stack, frame{i: current.i, j: current.j, done: true})

			if sec1[current.i-1] == sec2[current.j-1] {
				stack = append(stack, frame{i: current.i - 1, j: current.j - 1})
				continue
			}

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

// BacktrackingParallel explora el backtracking de forma concurrente cuando existen ramas independientes.
func BacktrackingParallel(sec1, sec2 string, matriz [][]int) []string {
	type key struct{ i, j int }

	type entry struct {
		once sync.Once
		res  map[string]struct{}
	}

	var (
		memo   = make(map[key]*entry)
		memoMu sync.Mutex
	)

	getEntry := func(i, j int) *entry {
		k := key{i: i, j: j}
		memoMu.Lock()
		e, ok := memo[k]
		if !ok {
			e = &entry{}
			memo[k] = e
		}
		memoMu.Unlock()
		return e
	}

	var compute func(i, j int) map[string]struct{}
	compute = func(i, j int) map[string]struct{} {
		e := getEntry(i, j)
		e.once.Do(func() {
			if matriz[i][j] == 0 {
				e.res = map[string]struct{}{"": {}}
				return
			}

			if i > 0 && j > 0 && sec1[i-1] == sec2[j-1] {
				prev := compute(i-1, j-1)
				res := make(map[string]struct{}, len(prev))
				for s := range prev {
					res[s+string(sec1[i-1])] = struct{}{}
				}
				e.res = res
				return
			}

			branches := make([]struct{ i, j int }, 0, 2)
			if i > 0 && matriz[i-1][j] == matriz[i][j] {
				branches = append(branches, struct{ i, j int }{i - 1, j})
			}
			if j > 0 && matriz[i][j-1] == matriz[i][j] {
				branches = append(branches, struct{ i, j int }{i, j - 1})
			}

			res := make(map[string]struct{})
			switch len(branches) {
			case 0:
				// sin ramas viables, mantenemos conjunto vacío
			case 1:
				sub := compute(branches[0].i, branches[0].j)
				for s := range sub {
					res[s] = struct{}{}
				}
			default:
				var (
					wg sync.WaitGroup
					mu sync.Mutex
				)
				for _, br := range branches {
					br := br
					wg.Add(1)
					go func() {
						defer wg.Done()
						sub := compute(br.i, br.j)
						mu.Lock()
						for s := range sub {
							res[s] = struct{}{}
						}
						mu.Unlock()
					}()
				}
				wg.Wait()
			}

			e.res = res
		})
		return e.res
	}

	set := compute(len(sec1), len(sec2))
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
