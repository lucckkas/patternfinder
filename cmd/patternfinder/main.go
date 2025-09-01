package main

import (
	"flag"
	"fmt"

	"github.com/lucckkas/patternfinder/internal/pattern"
)

func main() {
	topK := flag.Int("top", 1, "cuántos patrones devolver (top-K)")
	flag.Parse()
	seqs := flag.Args()
	if len(seqs) < 2 {
		fmt.Println("Se requieren al menos 2 secuencias.")
		return
	}

	if *topK <= 1 {
		agg, ok := pattern.BestCommonPattern(seqs)
		if !ok {
			fmt.Println("No hay patrón común.")
			return
		}
		fmt.Println("Patrón (mayúsculas):", agg.Pattern)
		fmt.Println("Gaps (min,max):      ", agg.GapRanges)
		fmt.Println("Formateado:          ", pattern.FormatPattern(agg))
		return
	}

	opts := pattern.DefaultTopKOptions()
	aggs, ok := pattern.TopKCommonPatterns(seqs, *topK, opts)
	if !ok {
		fmt.Println("No hay patrones.")
		return
	}
	for i, agg := range aggs {
		fmt.Printf("[%d] %s |\t %s\n",
			i+1, agg.Pattern, pattern.FormatPattern(agg))
	}
}
