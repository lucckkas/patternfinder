package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/lucckkas/patternfinder/internal/lcs"
)

// generateRandomSequence genera una secuencia aleatoria de letras mayúsculas
func generateRandomSequence(length int, seed int64) string {
	r := rand.New(rand.NewSource(seed))
	letters := "ABCDEFGHIKLMNPQRSTVWY"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = letters[r.Intn(len(letters))]
	}
	return string(result)
}

func main() {
	lengths := []int{20, 30, 40, 50, 60, 70, 80, 90, 100, 120, 140, 160, 180, 200}

	// Crear archivo CSV
	file, err := os.Create("benchmark_results/results.csv")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creando archivo: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezado
	writer.Write([]string{
		"Length",
		"Seq_DP_ms", "Seq_BT_ms", "Seq_Total_ms",
		"Par_DP_ms", "Par_BT_ms", "Par_Total_ms",
		"Speedup_DP", "Speedup_BT", "Speedup_Total",
		"LCS_Count",
	})

	fmt.Println("Generando resultados de benchmark...")
	fmt.Println("=====================================")
	fmt.Printf("%-6s | %-10s | %-10s | %-10s | %-8s\n",
		"Length", "Seq (ms)", "Par (ms)", "Speedup", "LCS#")
	fmt.Println("-------|------------|------------|------------|--------")

	for _, length := range lengths {
		// Generar secuencias aleatorias
		seq1 := generateRandomSequence(length, 12345)
		seq2 := generateRandomSequence(length, 67890)

		// Versión secuencial
		startDP := time.Now()
		dpSeq := lcs.DPTable(seq1, seq2)
		seqDPTime := time.Since(startDP)

		startBT := time.Now()
		lcsSeq := lcs.Backtracking(seq1, seq2, dpSeq)
		seqBTTime := time.Since(startBT)
		seqTotal := seqDPTime + seqBTTime

		// Versión paralela
		startDP = time.Now()
		dpPar := lcs.DPTableParallel(seq1, seq2)
		parDPTime := time.Since(startDP)

		startBT = time.Now()
		_ = lcs.BacktrackingParallel(seq1, seq2, dpPar)
		parBTTime := time.Since(startBT)
		parTotal := parDPTime + parBTTime

		// Calcular speedup
		speedupDP := float64(seqDPTime) / float64(parDPTime)
		speedupBT := float64(seqBTTime) / float64(parBTTime)
		speedupTotal := float64(seqTotal) / float64(parTotal)

		// Escribir fila CSV
		writer.Write([]string{
			fmt.Sprintf("%d", length),
			fmt.Sprintf("%.6f", float64(seqDPTime.Microseconds())/1000.0),
			fmt.Sprintf("%.6f", float64(seqBTTime.Microseconds())/1000.0),
			fmt.Sprintf("%.6f", float64(seqTotal.Microseconds())/1000.0),
			fmt.Sprintf("%.6f", float64(parDPTime.Microseconds())/1000.0),
			fmt.Sprintf("%.6f", float64(parBTTime.Microseconds())/1000.0),
			fmt.Sprintf("%.6f", float64(parTotal.Microseconds())/1000.0),
			fmt.Sprintf("%.4f", speedupDP),
			fmt.Sprintf("%.4f", speedupBT),
			fmt.Sprintf("%.4f", speedupTotal),
			fmt.Sprintf("%d", len(lcsSeq)),
		})

		// Imprimir en consola
		fmt.Printf("%-6d | %-10.2f | %-10.2f | %-10.2fx | %-8d\n",
			length,
			float64(seqTotal.Microseconds())/1000.0,
			float64(parTotal.Microseconds())/1000.0,
			speedupTotal,
			len(lcsSeq),
		)
	}

	fmt.Println("\nResultados guardados en: benchmark_results/results.csv")
	fmt.Println("Puedes usar este archivo para generar gráficos en Excel, Python, etc.")
}
