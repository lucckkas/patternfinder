package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

func main() {
	inputFile := flag.String("f", "", "archivo de texto con las secuencias (una por línea)")
	patternfinderPath := flag.String("p", "./build/patternfinder", "ruta al ejecutable de patternfinder")
	showDP := flag.Bool("dp", false, "pasar flag -dp a patternfinder")
	seq := flag.Bool("seq", false, "pasar flag -seq a patternfinder")
	outputFile := flag.String("o", "", "archivo de salida para los resultados (opcional, por defecto stdout)")
	workers := flag.Int("w", 6, "número de workers paralelos para ejecutar comparaciones")
	csvFile := flag.String("csv", "", "archivo CSV para guardar estadísticas de patrones")
	flag.Parse()

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Uso: %s -f <archivo_secuencias> [-p <path_patternfinder>] [-dp] [-seq] [-w <workers>] [-o <archivo_salida>] [-csv <archivo_csv>]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nEl archivo de secuencias debe contener una secuencia por línea.\n")
		fmt.Fprintf(os.Stderr, "Por defecto usa 4 workers para ejecutar comparaciones en paralelo.\n")
		fmt.Fprintf(os.Stderr, "Use -csv para generar un archivo CSV con estadísticas de patrones.\n")
		os.Exit(2)
	}

	// Leer las secuencias del archivo
	sequences, err := readSequences(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al leer el archivo: %v\n", err)
		os.Exit(1)
	}

	if len(sequences) < 2 {
		fmt.Fprintf(os.Stderr, "Se necesitan al menos 2 secuencias en el archivo.\n")
		os.Exit(1)
	}

	fmt.Printf("Leyendo %d secuencias del archivo %s\n", len(sequences), *inputFile)
	fmt.Printf("Total de comparaciones: %d\n", (len(sequences)*(len(sequences)-1))/2)
	fmt.Printf("Ejecutando con %d workers en paralelo\n\n", *workers)

	// Configurar salida
	var output *os.File
	if *outputFile != "" {
		output, err = os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al crear archivo de salida: %v\n", err)
			os.Exit(1)
		}
		defer output.Close()
		fmt.Printf("Los resultados se guardarán en: %s\n\n", *outputFile)
	} else {
		output = os.Stdout
	}

	// Verificar que el ejecutable de patternfinder existe
	if _, err := os.Stat(*patternfinderPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "El ejecutable de patternfinder no existe en: %s\n", *patternfinderPath)
		os.Exit(1)
	}

	// Obtener la ruta absoluta del ejecutable
	absPath, err := filepath.Abs(*patternfinderPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al obtener ruta absoluta: %v\n", err)
		os.Exit(1)
	}

	// Estructura para almacenar los resultados de cada comparación
	type ComparisonResult struct {
		Index  int
		SeqI   int
		SeqJ   int
		Output string
		Error  error
	}

	// Crear lista de trabajos (pares de secuencias a comparar)
	type Job struct {
		Index int
		SeqI  int
		SeqJ  int
		Seq1  string
		Seq2  string
	}

	var jobs []Job
	comparisonCount := 0
	for i := 0; i < len(sequences); i++ {
		for j := i + 1; j < len(sequences); j++ {
			comparisonCount++
			jobs = append(jobs, Job{
				Index: comparisonCount,
				SeqI:  i + 1,
				SeqJ:  j + 1,
				Seq1:  sequences[i],
				Seq2:  sequences[j],
			})
		}
	}

	// Canal para enviar trabajos
	jobsChan := make(chan Job, len(jobs))
	// Canal para recibir resultados
	resultsChan := make(chan ComparisonResult, len(jobs))

	// Lanzar workers
	var wg sync.WaitGroup
	for w := 0; w < *workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobsChan {
				// Preparar los argumentos para patternfinder
				args := []string{}
				if *showDP {
					args = append(args, "-dp")
				}
				if *seq {
					args = append(args, "-seq")
				}
				args = append(args, job.Seq1, job.Seq2)

				// Ejecutar patternfinder
				cmd := exec.Command(absPath, args...)
				cmdOutput, err := cmd.CombinedOutput()

				result := ComparisonResult{
					Index:  job.Index,
					SeqI:   job.SeqI,
					SeqJ:   job.SeqJ,
					Output: string(cmdOutput),
					Error:  err,
				}
				resultsChan <- result
			}
		}(w)
	}

	// Enviar todos los trabajos al canal
	for _, job := range jobs {
		jobsChan <- job
	}
	close(jobsChan)

	// Esperar a que terminen todos los workers
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Recolectar todos los resultados
	results := make([]ComparisonResult, 0, len(jobs))
	for result := range resultsChan {
		results = append(results, result)
	}

	// Ordenar resultados por índice para mantener el orden original
	// Usamos un mapa para acceso rápido por índice
	resultMap := make(map[int]ComparisonResult)
	for _, r := range results {
		resultMap[r.Index] = r
	}

	// Mapa para recolectar estadísticas de patrones
	patternStats := make(map[string]*PatternStat)

	// Escribir resultados en orden y recolectar patrones
	for i := 1; i <= len(jobs); i++ {
		result := resultMap[i]
		fmt.Fprintf(output, "========================================\n")
		fmt.Fprintf(output, "Comparación %d: Secuencia %d vs Secuencia %d\n", result.Index, result.SeqI, result.SeqJ)
		fmt.Fprintf(output, "========================================\n")

		if result.Error != nil {
			fmt.Fprintf(output, "Error al ejecutar patternfinder: %v\n", result.Error)
			fmt.Fprintf(output, "Salida: %s\n", result.Output)
		} else {
			fmt.Fprintf(output, "%s", result.Output)
			// Extraer patrones de la salida
			extractPatterns(result.Output, patternStats, result.SeqI, result.SeqJ)
		}

		fmt.Fprintf(output, "\n")
	}

	fmt.Printf("\nComparaciones completadas: %d\n", comparisonCount)
	if *outputFile != "" {
		fmt.Printf("Resultados guardados en: %s\n", *outputFile)
	}

	// Generar CSV si se especificó
	if *csvFile != "" {
		err := generateCSV(*csvFile, patternStats, len(sequences))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al generar CSV: %v\n", err)
		} else {
			fmt.Printf("Estadísticas de patrones guardadas en: %s\n", *csvFile)
		}
	}
}

// readSequences lee un archivo de texto y retorna un slice con las secuencias
// Ignora líneas vacías y elimina espacios en blanco al inicio/final
func readSequences(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sequences []string
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Ignorar líneas vacías y comentarios
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		sequences = append(sequences, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return sequences, nil
}

// PatternStat almacena estadísticas de un patrón
type PatternStat struct {
	Pattern         string
	UppercaseCount  int
	SequenceIndices map[int]bool // Índices de secuencias que contienen este patrón
}

// countUppercase cuenta las letras mayúsculas en una cadena
func countUppercase(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsUpper(r) {
			count++
		}
	}
	return count
}

// extractPatterns extrae los patrones de la salida de patternfinder
func extractPatterns(output string, stats map[string]*PatternStat, seqI, seqJ int) {
	// Buscar líneas que contienen patrones con el formato [n.m] pattern
	// Ejemplo: [1.1] A-x(2)-B-x(3)-C
	re := regexp.MustCompile(`\[\d+\.\d+\]\s+(.+)`)
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			pattern := strings.TrimSpace(matches[1])
			// Ignorar patrones vacíos
			if pattern == "" {
				continue
			}

			// Si el patrón no existe, crearlo
			if _, exists := stats[pattern]; !exists {
				stats[pattern] = &PatternStat{
					Pattern:         pattern,
					UppercaseCount:  countUppercase(pattern),
					SequenceIndices: make(map[int]bool),
				}
			}

			// Agregar las secuencias que contienen este patrón
			stats[pattern].SequenceIndices[seqI] = true
			stats[pattern].SequenceIndices[seqJ] = true
		}
	}
}

// generateCSV genera un archivo CSV con las estadísticas de patrones
func generateCSV(filename string, stats map[string]*PatternStat, totalSequences int) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezado
	err = writer.Write([]string{"Patrón", "Cantidad de Mayúsculas", "Cantidad de Secuencias", "Porcentaje de Secuencias"})
	if err != nil {
		return err
	}

	// Ordenar patrones para salida consistente
	var patterns []string
	for pattern := range stats {
		patterns = append(patterns, pattern)
	}

	// Escribir datos
	for _, pattern := range patterns {
		stat := stats[pattern]
		seqCount := len(stat.SequenceIndices)
		percentage := float64(seqCount) / float64(totalSequences) * 100

		row := []string{
			stat.Pattern,
			strconv.Itoa(stat.UppercaseCount),
			strconv.Itoa(seqCount),
			fmt.Sprintf("%.2f%%", percentage),
		}

		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}
