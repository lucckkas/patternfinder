package lcs_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/lucckkas/patternfinder/internal/lcs"
)

// generateRandomSequence genera una secuencia aleatoria de letras mayúsculas
func generateRandomSequence(length int, seed int64) string {
	r := rand.New(rand.NewSource(seed))
	letters := "ABCDEFGHIKLMNPQRSTVWY" // Aminoácidos comunes
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = letters[r.Intn(len(letters))]
	}
	return string(result)
}

// BenchmarkSequentialVsParallel compara el rendimiento de las versiones secuencial y paralela
// con secuencias de longitud creciente
func BenchmarkSequentialVsParallel(b *testing.B) {
	// Configuración de longitudes a probar
	lengths := []int{20, 30, 40, 50, 60, 70, 80, 90, 100, 120, 140, 160, 180, 200}

	for _, length := range lengths {
		// Generar secuencias aleatorias para esta longitud
		seq1 := generateRandomSequence(length, 12345)
		seq2 := generateRandomSequence(length, 67890)

		b.Run(fmt.Sprintf("Sequential_len%d", length), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				dp := lcs.DPTable(seq1, seq2)
				_ = lcs.Backtracking(seq1, seq2, dp)
			}
		})

		b.Run(fmt.Sprintf("Parallel_len%d", length), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				dp := lcs.DPTableParallel(seq1, seq2)
				_ = lcs.BacktrackingParallel(seq1, seq2, dp)
			}
		})
	}
}

// TestSequentialVsParallelComparison ejecuta una comparación detallada con métricas
func TestSequentialVsParallelComparison(t *testing.T) {
	lengths := []int{20, 30, 40, 50, 60, 70, 80, 90, 100, 120, 140, 160, 180, 200}

	type result struct {
		length       int
		seqDP        time.Duration
		seqBT        time.Duration
		seqTotal     time.Duration
		parDP        time.Duration
		parBT        time.Duration
		parTotal     time.Duration
		speedupDP    float64
		speedupBT    float64
		speedupTotal float64
		lcsCount     int
	}

	var results []result

	t.Log("Iniciando comparación secuencial vs paralelo...")
	t.Log("========================================")

	for _, length := range lengths {
		// Generar secuencias aleatorias
		seq1 := generateRandomSequence(length, 12345)
		seq2 := generateRandomSequence(length, 67890)

		r := result{length: length}

		// Versión secuencial
		startDP := time.Now()
		dpSeq := lcs.DPTable(seq1, seq2)
		r.seqDP = time.Since(startDP)

		startBT := time.Now()
		lcsSeq := lcs.Backtracking(seq1, seq2, dpSeq)
		r.seqBT = time.Since(startBT)
		r.seqTotal = r.seqDP + r.seqBT
		r.lcsCount = len(lcsSeq)

		// Versión paralela
		startDP = time.Now()
		dpPar := lcs.DPTableParallel(seq1, seq2)
		r.parDP = time.Since(startDP)

		startBT = time.Now()
		lcsPar := lcs.BacktrackingParallel(seq1, seq2, dpPar)
		r.parBT = time.Since(startBT)
		r.parTotal = r.parDP + r.parBT

		// Verificar que ambos encuentran el mismo número de LCS
		if len(lcsSeq) != len(lcsPar) {
			t.Errorf("Longitud %d: diferencia en número de LCS (seq=%d, par=%d)",
				length, len(lcsSeq), len(lcsPar))
		}

		// Calcular speedup
		r.speedupDP = float64(r.seqDP) / float64(r.parDP)
		r.speedupBT = float64(r.seqBT) / float64(r.parBT)
		r.speedupTotal = float64(r.seqTotal) / float64(r.parTotal)

		results = append(results, r)
	}

	// Imprimir resultados en formato tabla
	t.Log("\nResultados detallados:")
	t.Log("========================================")
	t.Logf("%-6s | %-12s | %-12s | %-12s | %-8s | %-8s | %-8s | %-6s",
		"Len", "Seq Total", "Par Total", "Speedup", "Seq DP", "Par DP", "Seq BT", "Par BT")
	t.Log("-------|--------------|--------------|--------------|----------|----------|----------|--------")

	for _, r := range results {
		t.Logf("%-6d | %-12s | %-12s | %-12.2fx | %-8s | %-8s | %-8s | %-8s",
			r.length,
			r.seqTotal,
			r.parTotal,
			r.speedupTotal,
			r.seqDP,
			r.parDP,
			r.seqBT,
			r.parBT,
		)
	}

	t.Log("\nSpeedup por componente:")
	t.Log("========================================")
	t.Logf("%-6s | %-12s | %-12s | %-12s | %-6s",
		"Len", "DP Speedup", "BT Speedup", "Total Speedup", "LCS#")
	t.Log("-------|--------------|--------------|--------------|------")

	for _, r := range results {
		t.Logf("%-6d | %-12.2fx | %-12.2fx | %-12.2fx | %-6d",
			r.length,
			r.speedupDP,
			r.speedupBT,
			r.speedupTotal,
			r.lcsCount,
		)
	}

	// Calcular promedios
	var avgSpeedupDP, avgSpeedupBT, avgSpeedupTotal float64
	for _, r := range results {
		avgSpeedupDP += r.speedupDP
		avgSpeedupBT += r.speedupBT
		avgSpeedupTotal += r.speedupTotal
	}
	n := float64(len(results))
	avgSpeedupDP /= n
	avgSpeedupBT /= n
	avgSpeedupTotal /= n

	t.Log("\nPromedios:")
	t.Log("========================================")
	t.Logf("Speedup promedio DP:    %.2fx", avgSpeedupDP)
	t.Logf("Speedup promedio BT:    %.2fx", avgSpeedupBT)
	t.Logf("Speedup promedio Total: %.2fx", avgSpeedupTotal)
}

// BenchmarkDPTableOnly compara solo la construcción de la tabla DP
func BenchmarkDPTableOnly(b *testing.B) {
	lengths := []int{50, 100, 150, 200}

	for _, length := range lengths {
		seq1 := generateRandomSequence(length, 12345)
		seq2 := generateRandomSequence(length, 67890)

		b.Run(fmt.Sprintf("Sequential_len%d", length), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = lcs.DPTable(seq1, seq2)
			}
		})

		b.Run(fmt.Sprintf("Parallel_len%d", length), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = lcs.DPTableParallel(seq1, seq2)
			}
		})
	}
}

// BenchmarkBacktrackingOnly compara solo el backtracking
func BenchmarkBacktrackingOnly(b *testing.B) {
	lengths := []int{50, 100, 150, 200}

	for _, length := range lengths {
		seq1 := generateRandomSequence(length, 12345)
		seq2 := generateRandomSequence(length, 67890)

		// Pre-calcular la tabla DP
		dp := lcs.DPTable(seq1, seq2)

		b.Run(fmt.Sprintf("Sequential_len%d", length), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = lcs.Backtracking(seq1, seq2, dp)
			}
		})

		b.Run(fmt.Sprintf("Parallel_len%d", length), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = lcs.BacktrackingParallel(seq1, seq2, dp)
			}
		})
	}
}
