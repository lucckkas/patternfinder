package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/lucckkas/Memoria/pkg/discovery"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		fmt.Fprintln(os.Stderr, "Usage: discovery <sequence1> <sequence2>")
		os.Exit(1)
	}
	seq1, seq2 := flag.Arg(0), flag.Arg(1)

	// Iniciamos el cronómetro justo antes de la fase principal
	start := time.Now()

	// 1) Generar subsecuencias
	sub1 := discovery.GenerateSubsequences(seq1)
	sub2 := discovery.GenerateSubsequences(seq2)

	// 2) Comparar y obtener patrones
	pats := discovery.CompareSubsequences(sub1, sub2)

	// 3) Puntuar patrones
	maxLen := max(len(seq1), len(seq2))
	scores := discovery.ScorePatterns(pats, maxLen)

	// 4) Ordenar patrones por puntuación descendente
	sort.Slice(pats, func(i, j int) bool {
		return scores[pats[i]] > scores[pats[j]]
	})

	// Calculamos y mostramos el tiempo transcurrido
	elapsed := time.Since(start)
	fmt.Printf("Duration: %s\n", elapsed)

	// 5) Imprimir resultados
	fmt.Println("Patterns found and scored:")
	for _, p := range pats {
		fmt.Printf("%s \t(%d)\n", p, scores[p])
	}

}
