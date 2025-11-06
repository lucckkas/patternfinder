package test

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lucckkas/patternfinder/internal/lcs"
)

// ======================== Config ========================

var (
	outPath    = getenv("BENCH_OUT", "length_sweep.csv")
	seed       = atoiEnv64("BENCH_SEED", 42)
	alphabet   = getenv("BENCH_ALPHABET", "ACDEFGHIKLMNPQRSTVWY")
	lengths    = []int{25, 50, 75, 100, 125, 150, 175}
	csvWriteMu sync.Mutex
)

func atoiEnv64(k string, def int64) int64 {
	if v := os.Getenv(k); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return def
}

// ==================== Generación de datos ====================

func makeSeq(rng *rand.Rand, n int, alphabet string) string {
	var b strings.Builder
	b.Grow(n)
	m := len(alphabet)
	for i := 0; i < n; i++ {
		b.WriteByte(alphabet[rng.Intn(m)])
	}
	return b.String()
}

// ===================== CSV (logging) =====================

func appendCSVRow(
	algo string,
	length int,
	iters int,
	total time.Duration,
) {
	// No guardar si solo hay 1 (o 0) iteración
	if iters <= 1 {
		return
	}

	csvWriteMu.Lock()
	defer csvWriteMu.Unlock()

	// Determinar tipo: "paralelo" o "secuencial"
	tipo := "desconocido"
	l := strings.ToLower(algo)
	switch {
	case strings.Contains(l, "parallel"):
		tipo = "paralelo"
	case strings.Contains(l, "sequential"):
		tipo = "secuencial"
	}

	avg := total / time.Duration(iters)

	f, err := os.OpenFile(outPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "no se pudo abrir %q: %v\n", outPath, err)
		return
	}
	defer f.Close()

	stat, _ := f.Stat()
	needHeader := stat.Size() == 0

	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		fmt.Fprintf(os.Stderr, "seek falló en %q: %v\n", outPath, err)
		return
	}

	w := csv.NewWriter(f)
	if needHeader {
		_ = w.Write([]string{
			"Tipo",
			"Longitud",
			"Tiempo promedio (ns)",
		})
	}

	_ = w.Write([]string{
		tipo,
		strconv.Itoa(length),
		strconv.FormatInt(avg.Nanoseconds(), 10),
	})
	w.Flush()
}

// ===================== Núcleo benchmark =====================

func runOne(b *testing.B, algoName string, seq1, seq2 string, builder dpBuilder, tracker backtracker) {
	b.Helper()
	b.ReportAllocs()

	// Corrida base (fuera de tiempo) para asegurar consistencia
	base := tracker(seq1, seq2, builder(seq1, seq2))
	if len(base) == 0 {
		b.Fatalf("no se obtuvo ningún patrón en la corrida base (algo=%s)", algoName)
	}

	// Tramo medido
	b.ResetTimer()
	start := time.Now()
	for i := 0; i < b.N; i++ {
		got := tracker(seq1, seq2, builder(seq1, seq2))
		if len(got) != len(base) {
			b.Fatalf("cantidad de patrones inesperada: got=%d want=%d (algo=%s)", len(got), len(base), algoName)
		}
	}
	elapsed := time.Since(start)
	b.StopTimer()

	// Log CSV (fuera de tiempo)
	appendCSVRow(algoName, len(seq1), b.N, elapsed)
}

// Variante que mide solo la construcción de la tabla DP (sin backtracking).
func runDPOnly(b *testing.B, algoName string, seq1, seq2 string, builder dpBuilder) {
	b.Helper()
	b.ReportAllocs()

	// Corrida base (fuera de tiempo) para obtener la longitud esperada de LCS
	base := builder(seq1, seq2)
	want := base[len(seq1)][len(seq2)]
	if want == 0 {
		// Puede ser 0 según las secuencias; no es error, pero lo dejamos explícito
	}

	// Tramo medido
	b.ResetTimer()
	start := time.Now()
	for i := 0; i < b.N; i++ {
		got := builder(seq1, seq2)
		if got[len(seq1)][len(seq2)] != want {
			b.Fatalf("longitud LCS inesperada: got=%d want=%d (algo=%s)", got[len(seq1)][len(seq2)], want, algoName)
		}
	}
	elapsed := time.Since(start)
	b.StopTimer()

	// Log CSV (fuera de tiempo)
	appendCSVRow(algoName, len(seq1), b.N, elapsed)
}

// =================== Bench entrypoints ===================

func BenchmarkLengthSweep_Sequential(b *testing.B) {
	rng := rand.New(rand.NewSource(seed))
	for _, L := range lengths {
		// Generar datos fuera del tramo medido
		s1 := makeSeq(rng, L, alphabet)
		// Para evitar casos triviales idénticos, genero otra con una semilla desplazada
		altRng := rand.New(rand.NewSource(seed + int64(L)*7919)) // 7919 primo
		s2 := makeSeq(altRng, L, alphabet)

		b.Run(fmt.Sprintf("Len_%03d", L), func(b *testing.B) {
			runOne(b, "Sequential(DPTable+Backtracking)", s1, s2, lcs.DPTable, lcs.Backtracking)
		})
	}
}

func BenchmarkLengthSweep_Parallel(b *testing.B) {
	rng := rand.New(rand.NewSource(seed + 1234567))
	for _, L := range lengths {
		s1 := makeSeq(rng, L, alphabet)
		altRng := rand.New(rand.NewSource(seed + int64(L)*104729)) // 104729 primo
		s2 := makeSeq(altRng, L, alphabet)

		b.Run(fmt.Sprintf("Len_%03d", L), func(b *testing.B) {
			runOne(b, "Parallel(DPTableParallel+BacktrackingParallel)", s1, s2, lcs.DPTableParallel, lcs.BacktrackingParallel)
		})
	}
}

// Benchmarks adicionales para aislar el costo de DP puro (sin enumerar todas las LCS),
// útil para verificar monotonía ~O(n*m) y descartar la variabilidad del backtracking.
func BenchmarkLengthSweep_DPOnly_Sequential(b *testing.B) {
	rng := rand.New(rand.NewSource(seed + 987654321))
	for _, L := range lengths {
		s1 := makeSeq(rng, L, alphabet)
		altRng := rand.New(rand.NewSource(seed + int64(L)*1223)) // 1223 primo
		s2 := makeSeq(altRng, L, alphabet)

		b.Run(fmt.Sprintf("Len_%03d", L), func(b *testing.B) {
			runDPOnly(b, "Sequential(DPTableOnly)", s1, s2, lcs.DPTable)
		})
	}
}

func BenchmarkLengthSweep_DPOnly_Parallel(b *testing.B) {
	rng := rand.New(rand.NewSource(seed + 192837465))
	for _, L := range lengths {
		s1 := makeSeq(rng, L, alphabet)
		altRng := rand.New(rand.NewSource(seed + int64(L)*7919)) // reutilizamos otro primo
		s2 := makeSeq(altRng, L, alphabet)

		b.Run(fmt.Sprintf("Len_%03d", L), func(b *testing.B) {
			runDPOnly(b, "Parallel(DPTableParallelOnly)", s1, s2, lcs.DPTableParallel)
		})
	}
}
