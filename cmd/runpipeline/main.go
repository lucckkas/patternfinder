package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Segmentos representa la estructura del JSON generado por Interactions.py
// map[protein_name]map[ligand_id][]string
type Segmentos map[string]map[string][]string

func main() {
	// Flags
	cifDir := flag.String("d", "./cifs", "directorio con archivos .cif")
	ligand := flag.String("l", "ZN", "código del ligando a analizar (ej: ZN, MG, CA)")
	pythonScript := flag.String("py", "./Interactions.py", "ruta a Interactions.py")
	batchcompare := flag.String("b", "./build/batchcompare", "ruta al ejecutable batchcompare")
	outputCSV := flag.String("o", "resultados.csv", "archivo CSV de salida")
	workers := flag.Int("w", 6, "número de workers para batchcompare")
	distance := flag.Float64("dist", 4.0, "distancia de interacción en Å")
	flag.Parse()

	fmt.Println("=== Pipeline Completo: CIF → Interactions → BatchCompare ===")

	// 1. Verificar que el directorio existe
	if _, err := os.Stat(*cifDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: directorio no encontrado: %s\n", *cifDir)
		os.Exit(1)
	}

	// 2. Obtener todos los archivos .cif
	cifFiles, err := filepath.Glob(filepath.Join(*cifDir, "*.cif"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al buscar archivos .cif: %v\n", err)
		os.Exit(1)
	}

	if len(cifFiles) == 0 {
		fmt.Fprintf(os.Stderr, "No se encontraron archivos .cif en %s\n", *cifDir)
		os.Exit(1)
	}

	fmt.Printf("Encontrados %d archivos .cif\n\n", len(cifFiles))

	// 3. Crear archivo temporal para segmentos JSON
	segmentosFile := "temp_segmentos.json"
	// defer os.Remove(segmentosFile)

	// 4. Procesar cada archivo .cif con Interactions.py
	fmt.Println("=== Paso 1: Extracción de segmentos con Interactions.py ===")
	for i, cifPath := range cifFiles {
		fmt.Printf("[%d/%d] Procesando %s...\n", i+1, len(cifFiles), filepath.Base(cifPath))

		cmd := exec.Command("python3", *pythonScript,
			cifPath,
			"-d", fmt.Sprintf("%.1f", *distance),
			"-o", segmentosFile,
		)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error al ejecutar Interactions.py en %s: %v\n", cifPath, err)
			fmt.Fprintf(os.Stderr, "Salida: %s\n", string(output))
			continue
		}
		// Mostrar resumen (última línea generalmente indica éxito)
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Segmento interactuante") {
				fmt.Printf("  %s\n", line)
			}
		}
	}

	// 5. Leer el JSON generado
	fmt.Println("\n=== Paso 2: Lectura de segmentos ===")
	jsonData, err := os.ReadFile(segmentosFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al leer %s: %v\n", segmentosFile, err)
		os.Exit(1)
	}

	var segmentos Segmentos
	if err := json.Unmarshal(jsonData, &segmentos); err != nil {
		fmt.Fprintf(os.Stderr, "Error al parsear JSON: %v\n", err)
		os.Exit(1)
	}

	// 6. Extraer secuencias del ligando especificado
	var sequences []string
	ligandUpper := strings.ToUpper(*ligand)

	for proteinName, ligands := range segmentos {
		for ligandID, seqs := range ligands {
			// Verificar si el ligandID contiene el código del ligando
			if strings.HasPrefix(ligandID, ligandUpper+"_") {
				fmt.Printf("Proteína %s, Ligando %s: %d segmentos\n", proteinName, ligandID, len(seqs))
				sequences = append(sequences, seqs...)
			}
		}
	}

	if len(sequences) == 0 {
		fmt.Fprintf(os.Stderr, "No se encontraron segmentos para el ligando %s\n", *ligand)
		os.Exit(1)
	}

	fmt.Printf("\nTotal de secuencias extraídas: %d\n", len(sequences))

	// 7. Guardar secuencias en archivo temporal
	seqFile := "temp_sequences.txt"
	// defer os.Remove(seqFile)

	f, err := os.Create(seqFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al crear archivo de secuencias: %v\n", err)
		os.Exit(1)
	}

	for _, seq := range sequences {
		fmt.Fprintln(f, seq)
	}
	f.Close()

	fmt.Printf("Secuencias guardadas en %s\n", seqFile)

	// 8. Ejecutar batchcompare
	fmt.Println("\n=== Paso 3: Comparación con BatchCompare ===")

	absPath, err := filepath.Abs(*batchcompare)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error al obtener ruta absoluta de batchcompare: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: batchcompare no encontrado en %s\n", absPath)
		os.Exit(1)
	}

	cmd := exec.Command(absPath,
		"-f", seqFile,
		"-csv", *outputCSV,
		"-w", fmt.Sprintf("%d", *workers),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al ejecutar batchcompare: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n=== Pipeline completado ===\n")
	fmt.Printf("Resultados guardados en: %s\n", *outputCSV)
}
