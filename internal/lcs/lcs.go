package lcs

import (
	"fmt"
	"runtime"
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

// DPTableParallel construye la tabla LCS usando paralelismo por diagonales
// con un pool de workers fijo (bounded goroutines).
func DPTableParallel(sec1, sec2 string) [][]int {
    n, m := len(sec1), len(sec2)
    dp := make([][]int, n+1)
    for i := range dp {
        dp[i] = make([]int, m+1)
    }

    type cell struct {
        i, j int
    }

    numWorkers := runtime.GOMAXPROCS(0)
    if numWorkers < 1 {
        numWorkers = 1
    }

    // Recorremos todas las diagonales i+j = const
    for diag := 2; diag <= n+m; diag++ {
        // Cálculo de rango [start, end] de i en esta diagonal
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

        // Canal de trabajos para esta diagonal
        jobs := make(chan cell)
        var wg sync.WaitGroup

        // Lanzamos un número fijo de workers para esta diagonal
        for w := 0; w < numWorkers; w++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                for c := range jobs {
                    i, j := c.i, c.j

                    if sec1[i-1] == sec2[j-1] {
                        dp[i][j] = dp[i-1][j-1] + 1
                    } else if dp[i-1][j] >= dp[i][j-1] {
                        dp[i][j] = dp[i-1][j]
                    } else {
                        dp[i][j] = dp[i][j-1]
                    }
                }
            }()
        }

        // Enviamos todos los trabajos de la diagonal
        for i := start; i <= end; i++ {
            j := diag - i
            if j < 1 || j > m {
                continue
            }
            jobs <- cell{i: i, j: j}
        }

        // Cerramos el canal y esperamos a que todos los workers terminen
        close(jobs)
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

// BacktrackingParallel encuentra todas las LCS usando la matriz de DP,
// con concurrencia limitada mediante un semáforo (pool acotado de goroutines).
func BacktrackingParallel(sec1, sec2 string, matriz [][]int) []string {
    type cellKey struct{ i, j int }

    // Caminos visitados: para cada celda, qué patrones ya se han explorado
    visited := make(map[cellKey]map[string]struct{})
    var visitedMu sync.Mutex

    // Resultados finales (LCS distintas)
    results := make(map[string]struct{})
    var resultsMu sync.Mutex

    var wg sync.WaitGroup

    // Semáforo para limitar número máximo de goroutines concurrentes
    maxGoroutines := runtime.GOMAXPROCS(0)
    if maxGoroutines < 1 {
        maxGoroutines = 1
    }
    // Puedes multiplicar por un factor si quieres más paralelismo profundidad-abajo
    sem := make(chan struct{}, maxGoroutines*4)

    // Función auxiliar para registrar si ya se visitó (i,j) con cierto patrón
    registerPath := func(i, j int, pattern string) bool {
        k := cellKey{i: i, j: j}
        visitedMu.Lock()
        defer visitedMu.Unlock()

        m, ok := visited[k]
        if !ok {
            m = make(map[string]struct{})
            visited[k] = m
        }
        if _, exists := m[pattern]; exists {
            return false
        }
        m[pattern] = struct{}{}
        return true
    }

    var compute func(i, j int, pattern string)

    compute = func(i, j int, pattern string) {
        defer wg.Done()

        // Caso base: borde de la tabla → patrón completo
        if i == 0 || j == 0 {
            resultsMu.Lock()
            results[pattern] = struct{}{}
            resultsMu.Unlock()
            return
        }

        curr := matriz[i][j]

        // Coincidencia de caracteres ⇒ movemos en diagonal y agregamos caracter
        if sec1[i-1] == sec2[j-1] {
            newPattern := string(sec1[i-1]) + pattern
            if !registerPath(i-1, j-1, newPattern) {
                return
            }

            wg.Add(1)
            // Intentar ejecutar en goroutine si hay cupo en el semáforo
            select {
            case sem <- struct{}{}:
                go func() {
                    defer func() { <-sem }()
                    compute(i-1, j-1, newPattern)
                }()
            default:
                // Sin cupo: seguir de forma secuencial
                compute(i-1, j-1, newPattern)
            }
            return
        }

        // En caso de no coincidencia, pueden existir hasta dos ramas:
        // arriba (i-1, j) y/o izquierda (i, j-1), manteniendo el valor curr.

        // Helper para lanzar una rama con patrón actual
        runBranch := func(ni, nj int) {
            if !registerPath(ni, nj, pattern) {
                return
            }
            wg.Add(1)
            select {
            case sem <- struct{}{}:
                go func() {
                    defer func() { <-sem }()
                    compute(ni, nj, pattern)
                }()
            default:
                compute(ni, nj, pattern)
            }
        }

        if i > 0 && matriz[i-1][j] == curr {
            runBranch(i-1, j)
        }
        if j > 0 && matriz[i][j-1] == curr {
            runBranch(i, j-1)
        }
    }

    // Lanzar el backtracking desde la esquina inferior derecha
    // con patrón vacío
    if !registerPath(len(sec1), len(sec2), "") {
        // Técnicamente no debería pasar, pero por si acaso
        return nil
    }

    wg.Add(1)
    compute(len(sec1), len(sec2), "")
    wg.Wait()

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
