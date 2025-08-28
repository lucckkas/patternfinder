package main

import (
	"fmt"
	"os"

	"github.com/lucckkas/patternfinder/internal/pattern"
)

func main() {
	seqs := os.Args[1:]
	if len(seqs) == 0 {
		fmt.Println("No se proporcionaron secuencias.")
		fmt.Println("Uso: patternfinder SEQ1 SEQ2 SEQ3 ...")
		return
	}

	agg, ok := pattern.BestCommonPattern(seqs)
	if !ok {
		fmt.Println("No hay patrón común.")
		return
	}

	fmt.Println("Patrón (mayúsculas):", agg.Pattern)
	fmt.Println("Gaps (min,max):      ", agg.GapRanges)
	fmt.Println("Formateado:          ", pattern.FormatPattern(agg))
}
