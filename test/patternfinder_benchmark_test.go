package test

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/lucckkas/patternfinder/internal/aggregate"
	"github.com/lucckkas/patternfinder/internal/gaps"
	"github.com/lucckkas/patternfinder/internal/lcs"
	"github.com/lucckkas/patternfinder/internal/utils"
)

var (
	benchOutPath   = defaultOutPath()
	writeMu        sync.Mutex
	benchmarkPairs []seqPair // se llena en init() leyendo el CSV

	// Control de volcado de patrones
	// BENCH_PATTERNS=main (default) -> guarda solo el patrón principal (más largo)
	// BENCH_PATTERNS=all             -> guarda todos los patrones (unidos por '|')
	saveAllPatterns = strings.ToLower(getenv("BENCH_PATTERNS", "main")) == "all"
	// Recorte visual para no reventar el CSV si el campo de patrones es enorme
)

type dpBuilder func(string, string) [][]int
type backtracker func(string, string, [][]int) []string

type record struct {
	ID      string
	Chain   string
	Seq     string // versión UpperOnly (para LCS)
	Orig    string // versión original tal cual del CSV (para gaps/formateo)
	RawLine []string
}

type seqPair struct {
	Name       string
	S1, S2     string // UpperOnly (para LCS)
	Raw1, Raw2 string // originales (para gaps/formateo)
	// opcional: metadata para el nombre
	ID1, Chain1 string
	ID2, Chain2 string
}

// ---------- Config & init ----------

func defaultOutPath() string {
	if p := os.Getenv("BENCH_OUT"); p != "" {
		return p
	}
	return "bench_results.csv"
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func atoiEnv(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// init carga los pares desde el CSV antes de correr benchmarks
func init() {
	csvPath := getenv("DIVERSE_CSV", "diverse_selection.csv")
	pairing := strings.ToLower(getenv("DIVERSE_PAIRING", "adjacent")) // adjacent|roundrobin|all
	maxPairs := atoiEnv("DIVERSE_MAX_PAIRS", 50)

	recs, warn := loadRecordsFromCSV(csvPath)
	if warn != "" {
		fmt.Fprintf(os.Stderr, "[WARN] %s\n", warn)
	}
	if len(recs) == 0 {
		fmt.Fprintf(os.Stderr, "[WARN] no se cargaron secuencias desde %q; se usarán secuencias de fallback.\n", csvPath)
		// Fallback mínimo para no romper el benchmark si el CSV no existe:
		recs = []record{
			{ID: "fallback1", Chain: "A", Orig: "DBDCABADBDCABADBDCABADBDCABADBDCABADBDCABADBDCABADBDCABADBDCABADBDCABADBDCABAD"},
			{ID: "fallback2", Chain: "B", Orig: "BABCBDABBBABCBDABBBABCBDABBBABCBDABBBABCBDABBBABCBDABBBABCBDABBBABCBDABBBABCB"},
			{ID: "fallback3", Chain: "C", Orig: "DBDCABADBDCABADBDCABADBDCAB"},
		}
		for i := range recs {
			recs[i].Seq = utils.UpperOnly(recs[i].Orig)
		}
	}

	benchmarkPairs = buildPairs(recs, pairing, maxPairs)
	if len(benchmarkPairs) == 0 {
		fmt.Fprintf(os.Stderr, "[WARN] no se pudieron construir pares con la estrategia %q\n", pairing)
	}
}

// ---------- Lectura CSV ----------

// Aceptamos varios nombres de columna para la secuencia completa (idealmente no truncada)
var candidateSeqCols = []string{
	"sequence", "secuencia",
	"firma usada para diversidad (completa)",
	"signaturefull", "signature_full", "signature",
}

func loadRecordsFromCSV(path string) ([]record, string) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Sprintf("no pude abrir %q: %v", path, err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return nil, fmt.Sprintf("no pude leer cabecera de %q: %v", path, err)
	}

	normalize := func(s string) string {
		s = strings.ToLower(strings.TrimSpace(s))
		// quitar tildes comunes
		replacer := strings.NewReplacer("á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u")
		return replacer.Replace(s)
	}

	colIdx := func(names ...string) int {
		for i, h := range header {
			hn := normalize(h)
			for _, name := range names {
				if hn == normalize(name) {
					return i
				}
			}
		}
		return -1
	}

	// 1) candidatos “exactos”
	seqCol := -1
	for _, cand := range candidateSeqCols {
		seqCol = colIdx(cand)
		if seqCol != -1 {
			break
		}
	}

	warn := ""
	// 2) si no encontramos, aceptamos la columna “acortada” (o cualquiera que contenga "firma usada para diversidad"/"signature")
	if seqCol == -1 {
		for i, h := range header {
			hn := normalize(h)
			if strings.Contains(hn, "firma usada para diversidad") || strings.Contains(hn, "signature") {
				seqCol = i
				if strings.Contains(hn, "acortada") || strings.Contains(hn, "muestra") {
					warn = fmt.Sprintf("usando columna %q (ACORTADA). Idealmente exporta también una columna completa (p.ej. 'SignatureFull').", header[i])
				}
				break
			}
		}
	}

	// columnas opcionales para nombrar pares
	idCol := colIdx("pdb id", "pdb_id", "id", "id original")
	chainCol := colIdx("cadena", "chain")

	// Si no hay columna de secuencia, no podemos seguir
	if seqCol == -1 {
		return nil, "no encontré una columna de secuencia válida en el CSV (busco 'Sequence', 'SignatureFull' o 'Firma usada para diversidad (...)')."
	}

	var out []record
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			warn = warn + fmt.Sprintf(" (error de lectura: %v)", err)
			break
		}
		if seqCol >= len(row) || row[seqCol] == "" {
			continue
		}
		id := ""
		chain := ""
		if idCol >= 0 && idCol < len(row) {
			id = row[idCol]
		}
		if chainCol >= 0 && chainCol < len(row) {
			chain = row[chainCol]
		}
		raw := strings.TrimSpace(row[seqCol]) // original tal cual (para gaps/formateo)
		seq := utils.UpperOnly(raw)           // versión UpperOnly (para LCS)
		out = append(out, record{
			ID:      strings.TrimSpace(id),
			Chain:   strings.TrimSpace(chain),
			Seq:     seq,
			Orig:    raw,
			RawLine: row,
		})
	}
	return out, warn
}

// ---------- Construcción de pares ----------

func buildPairs(recs []record, strategy string, maxPairs int) []seqPair {
	var pairs []seqPair
	n := len(recs)
	nameOf := func(a, b record) string {
		left := a.ID
		right := b.ID
		if left == "" {
			left = "rec" // fallback
		}
		if right == "" {
			right = "rec"
		}
		return fmt.Sprintf("%s%s_vs_%s%s",
			left, suffixIfNotEmpty(a.Chain),
			right, suffixIfNotEmpty(b.Chain),
		)
	}

	switch strategy {
	case "all":
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				pairs = append(pairs, seqPair{
					Name: nameOf(recs[i], recs[j]),
					S1:   recs[i].Seq, S2: recs[j].Seq,
					Raw1: recs[i].Orig, Raw2: recs[j].Orig,
					ID1: recs[i].ID, Chain1: recs[i].Chain,
					ID2: recs[j].ID, Chain2: recs[j].Chain,
				})
				if maxPairs > 0 && len(pairs) >= maxPairs {
					return pairs
				}
			}
		}
	case "roundrobin":
		half := n / 2
		for i := 0; i < half && half+i < n; i++ {
			pairs = append(pairs, seqPair{
				Name: nameOf(recs[i], recs[half+i]),
				S1:   recs[i].Seq, S2: recs[half+i].Seq,
				Raw1: recs[i].Orig, Raw2: recs[half+i].Orig,
				ID1: recs[i].ID, Chain1: recs[i].Chain,
				ID2: recs[half+i].ID, Chain2: recs[half+i].Chain,
			})
			if maxPairs > 0 && len(pairs) >= maxPairs {
				return pairs
			}
		}
	default: // "adjacent"
		for i := 0; i+1 < n; i += 2 {
			pairs = append(pairs, seqPair{
				Name: nameOf(recs[i], recs[i+1]),
				S1:   recs[i].Seq, S2: recs[i+1].Seq,
				Raw1: recs[i].Orig, Raw2: recs[i+1].Orig,
				ID1: recs[i].ID, Chain1: recs[i].Chain,
				ID2: recs[i+1].ID, Chain2: recs[i+1].Chain,
			})
			if maxPairs > 0 && len(pairs) >= maxPairs {
				return pairs
			}
		}
	}
	return pairs
}

func suffixIfNotEmpty(s string) string {
	if s == "" {
		return ""
	}
	return "_" + s
}

// ---------- Utilidades de patrones ----------

func pickPatternFieldFormatted(raw1, raw2 string, patterns []string) (main string, field string) {
	// principal = el más largo (sobre el patrón crudo)
	for _, p := range patterns {
		fmt.Println("DEBUG largo patrón:", len(p))
		if len(p) > len(main) {
			main = p
		}
	}

	// helper para formatear (como en main.go)
	formatOne := func(p string) string {
		setsX, okX := gaps.AllGapValuesDistanceTotalViable(raw1, p)
		setsY, okY := gaps.AllGapValuesDistanceTotalViable(raw2, p)
		if !okX || !okY {
			return p // fallback si no se pudieron calcular gaps
		}
		union := aggregate.PairUnionSets(setsX, setsY)
		return aggregate.FormatPatternWithValues(p, union)
	}

	if saveAllPatterns {
		formatted := make([]string, 0, len(patterns))
		for _, p := range patterns {
			formatted = append(formatted, formatOne(p))
		}
		field = strings.Join(formatted, "|")
	} else {
		field = formatOne(main)
	}
	return main, field
}

// ---------- Bench core (parametrizado por secuencias) ----------

func benchmarkLCSWithSeqs(b *testing.B, raw1, raw2, seq1, seq2 string, builder dpBuilder, tracker backtracker) {
	// Aseguramos mayúsculas para LCS
	seq1 = utils.UpperOnly(seq1)
	seq2 = utils.UpperOnly(seq2)

	b.Helper()
	b.ReportAllocs()

	// Corrida base (no medida): usamos estos patrones (crudos) para el CSV
	patternsBase := runPipeline(seq1, seq2, builder, tracker)
	baseline := len(patternsBase)
	if baseline == 0 {
		b.Fatal("se esperaba al menos un patrón en la corrida base")
	}
	mainPat, patternField := pickPatternFieldFormatted(raw1, raw2, patternsBase)

	builderName := shortFuncName(builder)
	trackerName := shortFuncName(tracker)

	// ----- Sección medida -----
	b.ResetTimer()
	start := time.Now()

	for i := 0; i < b.N; i++ {
		patterns := runPipeline(seq1, seq2, builder, tracker)
		if len(patterns) != baseline {
			b.Fatalf("cantidad de patrones inesperada: got %d, want %d", len(patterns), baseline)
		}
	}

	elapsed := time.Since(start)
	b.StopTimer()
	// ----- Fin sección medida -----

	// Log a CSV (fuera de tiempo)
	writeCSVRow(
		b,
		b.Name(),          // Nombre del (sub)benchmark
		builderName,       // Constructor DP
		trackerName,       // Backtracking
		shorten(seq1, 80), // muestra seq1
		shorten(seq2, 80), // muestra seq2
		len(seq1),         // long seq1
		len(seq2),         // long seq2
		baseline,          // #patrones encontrados (crudos)
		len(mainPat),      // longitud del patrón principal (crudo)
		patternField,      // patrón(es) formateado(s) (como en main.go)
		b.N,               // iteraciones
		elapsed,           // total
	)
}

func runPipeline(seq1, seq2 string, builder dpBuilder, tracker backtracker) []string {
	dp := builder(seq1, seq2)
	return tracker(seq1, seq2, dp)
}

// ---------- Bench entrypoints ----------

func BenchmarkPatternfinderSequential(b *testing.B) {
	if len(benchmarkPairs) == 0 {
		b.Skip("no hay pares para ejecutar")
	}
	for _, p := range benchmarkPairs {
		p := p
		b.Run("Seq_"+sanitizeName(p.Name), func(b *testing.B) {
			benchmarkLCSWithSeqs(b, p.Raw1, p.Raw2, p.S1, p.S2, lcs.DPTable, lcs.Backtracking)
		})
	}
}

func BenchmarkPatternfinderParallel(b *testing.B) {
	if len(benchmarkPairs) == 0 {
		b.Skip("no hay pares para ejecutar")
	}
	for _, p := range benchmarkPairs {
		p := p
		b.Run("Seq_"+sanitizeName(p.Name), func(b *testing.B) {
			benchmarkLCSWithSeqs(b, p.Raw1, p.Raw2, p.S1, p.S2, lcs.DPTableParallel, lcs.BacktrackingParallel)
		})
	}
}

// ----------------- Helpers varios -----------------

func shortFuncName(fn any) string {
	full := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	if idx := strings.LastIndex(full, "."); idx != -1 {
		return full[idx+1:]
	}
	return full
}

func shorten(s string, n int) string {
	if len(s) <= n {
		return s
	}
	half := n / 2
	return s[:half] + "…" + s[len(s)-half:]
}

func sanitizeName(s string) string {
	s = strings.ReplaceAll(s, string(filepath.Separator), "_")
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, ":", "_")
	return s
}

func writeCSVRow(
	b *testing.B,
	benchName, builderName, trackerName string,
	seq1Sample, seq2Sample string,
	seq1Len, seq2Len int,
	numPatterns int, // # patrones crudos
	mainPatternLen int, // len del patrón principal crudo
	patternField string, // patrón(es) formateado(s)
	iters int,
	total time.Duration,
) {
	writeMu.Lock()
	defer writeMu.Unlock()

	f, err := os.OpenFile(benchOutPath, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		b.Logf("no se pudo abrir %q: %v", benchOutPath, err)
		return
	}
	defer f.Close()

	stat, _ := f.Stat()
	needHeader := stat.Size() == 0

	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		b.Logf("seek falló en %q: %v", benchOutPath, err)
		return
	}

	w := csv.NewWriter(f)
	if needHeader {
		header := []string{
			"Nombre del benchmark",
			"Constructor DP",
			"Backtracking",
			"Secuencia 1 (muestra)",
			"Secuencia 2 (muestra)",
			"Longitud secuencia 1",
			"Longitud secuencia 2",
			"Cantidad de patrones",
			"Longitud patrón principal",
			"Patrón(es) encontrado(s) (formateado)",
			"Numero de iteraciones",
			"Tiempo total (ns)",
			"Tiempo por operacion (ns)",
		}
		if err := w.Write(header); err != nil {
			b.Logf("no se pudo escribir cabecera CSV: %v", err)
			return
		}
	}

	perOp := time.Duration(0)
	if iters > 0 {
		perOp = total / time.Duration(iters)
	}

	record := []string{
		benchName,
		builderName,
		trackerName,
		seq1Sample,
		seq2Sample,
		fmt.Sprintf("%d", seq1Len),
		fmt.Sprintf("%d", seq2Len),
		fmt.Sprintf("%d", numPatterns),
		fmt.Sprintf("%d", mainPatternLen),
		patternField,
		fmt.Sprintf("%d", iters),
		fmt.Sprintf("%d", total.Nanoseconds()),
		fmt.Sprintf("%d", perOp.Nanoseconds()),
	}

	if err := w.Write(record); err != nil {
		b.Logf("no se pudo escribir fila CSV: %v", err)
		return
	}
	w.Flush()
	if err := w.Error(); err != nil {
		b.Logf("error al cerrar writer CSV: %v", err)
	}
}
