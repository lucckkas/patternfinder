package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/lucckkas/patternfinder/internal/aggregate"
	"github.com/lucckkas/patternfinder/internal/gaps"
	"github.com/lucckkas/patternfinder/internal/lcs"
	"github.com/lucckkas/patternfinder/internal/utils"
)

func main() {
	showDP := flag.Bool("dp", false, "imprimir matriz LCS (longitudes)")
	seq := flag.Bool("seq", false, "usar versión secuencial del LCS")
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Uso: %s <seq1> <seq2>\n", os.Args[0])
		os.Exit(2)
	}

	seqX := args[0]
	seqY := args[1]

	Ux := utils.UpperOnly(seqX)
	Uy := utils.UpperOnly(seqY)

	if len(Ux) == 0 || len(Uy) == 0 {
		fmt.Println("No hay mayúsculas en alguna secuencia; no existe LCS.")
		return
	}

	var (
		dp  [][]int
		all []string
	)
	if *seq {
		dp = lcs.DPTable(Ux, Uy)
		all = lcs.AllLCS(Ux, Uy, dp)
	} else {
		dp = lcs.DPTableParallel(Ux, Uy)
		all = lcs.AllLCSParallel(Ux, Uy, dp)
	}
	if *showDP {
		fmt.Println("Matriz LCS (longitudes):")
		lcs.PrintDP(Ux, Uy, dp)
	}

	if len(all) == 0 {
		fmt.Println("No se encontraron LCS.")
		return
	}
	sort.Slice(all, func(i, j int) bool {
		if len(all[i]) != len(all[j]) {
			return len(all[i]) > len(all[j])
		}
		return all[i] < all[j]
	})

	fmt.Printf("Secuencia 1 (original): %s\n", seqX)
	fmt.Printf("Secuencia 2 (original): %s\n", seqY)
	fmt.Printf("Mayúsculas 1: %s\n", Ux)
	fmt.Printf("Mayúsculas 2: %s\n\n", Uy)
	fmt.Printf("LCS: %v\n\n", all)

	for idx, pat := range all {
		setsX, okX := gaps.AllGapValuesDistanceTotalViable(seqX, pat)
		setsY, okY := gaps.AllGapValuesDistanceTotalViable(seqY, pat)
		if !okX || !okY {
			fmt.Printf("[%d] %s -> (no se pudo calcular gaps)\n", idx+1, pat)
			continue
		}
		union := aggregate.PairUnionSets(setsX, setsY)
		formatted := aggregate.FormatPatternWithValues(pat, union)
		fmt.Printf("[%d] %s | vals=%v | %s\n", idx+1, pat, union, formatted)
	}
}
