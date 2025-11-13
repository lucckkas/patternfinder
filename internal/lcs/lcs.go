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

// Backtracking hace el backtracking para encontrar todas las LCS posibles.
// Utiliza detección de caminos duplicados para evitar explorar el mismo camino
// con el mismo patrón parcial.
func Backtracking(sec1, sec2 string, matriz [][]int) []string {
	type cellKey struct{ i, j int }
	
	// Registro de caminos visitados: para cada celda, guarda los patrones con los que se ha llegado
	visited := make(map[cellKey]map[string]bool)
	
	// Resultados finales
	results := make(map[string]struct{})

	type frame struct {
		i, j    int
		pattern string
	}

	stack := []frame{{i: len(sec1), j: len(sec2), pattern: ""}}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		
		k := cellKey{i: current.i, j: current.j}

		// Verificar si este camino ya fue explorado
		if visited[k] == nil {
			visited[k] = make(map[string]bool)
		}
		if visited[k][current.pattern] {
			continue // Este camino ya fue explorado con este patrón
		}
		visited[k][current.pattern] = true

		// Caso base: llegamos al origen
		if matriz[current.i][current.j] == 0 {
			results[current.pattern] = struct{}{}
			continue
		}

		// Si hay coincidencia de caracteres, avanzamos en diagonal
		if sec1[current.i-1] == sec2[current.j-1] {
			newPattern := string(sec1[current.i-1]) + current.pattern
			stack = append(stack, frame{
				i:       current.i - 1,
				j:       current.j - 1,
				pattern: newPattern,
			})
			continue
		}

		// Sin coincidencia: explorar ramas válidas
		if matriz[current.i-1][current.j] == matriz[current.i][current.j] {
			stack = append(stack, frame{
				i:       current.i - 1,
				j:       current.j,
				pattern: current.pattern,
			})
		}
		if matriz[current.i][current.j-1] == matriz[current.i][current.j] {
			stack = append(stack, frame{
				i:       current.i,
				j:       current.j - 1,
				pattern: current.pattern,
			})
		}
	}

	// Convertir el conjunto de resultados a slice
	out := make([]string, 0, len(results))
	for s := range results {
		out = append(out, s)
	}
	return out
}

// BacktrackingParallel explora el backtracking de forma concurrente cuando existen ramas independientes.
// Utiliza detección de caminos duplicados para evitar que múltiples goroutines exploren
// el mismo camino con el mismo patrón parcial.
func BacktrackingParallel(sec1, sec2 string, matriz [][]int) []string {
	type cellKey struct{ i, j int }
	type pathState struct {
		cell    cellKey
		pattern string
	}

	var (
		// Registro de caminos visitados: para cada celda, guarda los patrones con los que se ha llegado
		visited   = make(map[cellKey]map[string]bool)
		visitedMu sync.Mutex
		
		// Resultados finales
		results   = make(map[string]struct{})
		resultsMu sync.Mutex
	)

	// Registra que una goroutine llegó a una celda con un patrón específico
	// Retorna true si el camino es nuevo, false si ya estaba siendo explorado
	registerPath := func(i, j int, pattern string) bool {
		k := cellKey{i: i, j: j}
		visitedMu.Lock()
		defer visitedMu.Unlock()
		
		if visited[k] == nil {
			visited[k] = make(map[string]bool)
		}
		
		// Si el patrón ya fue registrado, este camino ya está siendo explorado
		if visited[k][pattern] {
			return false
		}
		
		visited[k][pattern] = true
		return true
	}

	var compute func(i, j int, currentPattern string, wg *sync.WaitGroup)
	compute = func(i, j int, currentPattern string, wg *sync.WaitGroup) {
		if wg != nil {
			defer wg.Done()
		}

		// Verificar si este camino ya está siendo explorado
		if !registerPath(i, j, currentPattern) {
			// Otra goroutine ya llegó aquí con el mismo patrón, detener esta goroutine
			return
		}

		// Caso base: llegamos al origen
		if matriz[i][j] == 0 {
			resultsMu.Lock()
			results[currentPattern] = struct{}{}
			resultsMu.Unlock()
			return
		}

		// Si hay coincidencia de caracteres, avanzamos en diagonal
		if i > 0 && j > 0 && sec1[i-1] == sec2[j-1] {
			newPattern := string(sec1[i-1]) + currentPattern
			compute(i-1, j-1, newPattern, nil)
			return
		}

		// Determinar las ramas válidas
		branches := make([]struct{ i, j int }, 0, 2)
		if i > 0 && matriz[i-1][j] == matriz[i][j] {
			branches = append(branches, struct{ i, j int }{i - 1, j})
		}
		if j > 0 && matriz[i][j-1] == matriz[i][j] {
			branches = append(branches, struct{ i, j int }{i, j - 1})
		}

		switch len(branches) {
		case 0:
			// Sin ramas viables, este es un camino inválido
			return
		case 1:
			// Una sola rama, continuar en la misma goroutine
			compute(branches[0].i, branches[0].j, currentPattern, nil)
		default:
			// Múltiples ramas, explorar en paralelo
			var branchWg sync.WaitGroup
			for _, br := range branches {
				branchWg.Add(1)
				go compute(br.i, br.j, currentPattern, &branchWg)
			}
			branchWg.Wait()
		}
	}

	// Iniciar el backtracking desde la esquina inferior derecha
	compute(len(sec1), len(sec2), "", nil)

	// Convertir el conjunto de resultados a slice
	out := make([]string, 0, len(results))
	for s := range results {
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
