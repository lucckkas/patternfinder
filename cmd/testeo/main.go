// cmd/discovery/main.go
package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type Anchor struct {
	Pos    int
	Letter rune
}

type Pattern struct {
	Seq   string
	Score int
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Uso: %s SEQ1 SEQ2 [SEQ3 ...]\n", os.Args[0])
		os.Exit(1)
	}
	seqs := os.Args[1:]
	pats := generatePatterns(seqs)
	fmt.Println("Mejores patrones:")
	for _, p := range pats {
		fmt.Printf("%s (%d)\n", p.Seq, p.Score)
	}
}

// ——————— Extracción de anclas ———————

func extractAnchors(s string) []Anchor {
	var a []Anchor
	for i, r := range s {
		if unicode.IsUpper(r) {
			a = append(a, Anchor{i, r})
		}
	}
	return a
}

// ——————— Índice de posiciones por letra ———————

func buildIndex(seqs []string) []map[rune][]int {
	idx := make([]map[rune][]int, len(seqs))
	for i, s := range seqs {
		m := make(map[rune][]int)
		for pos, r := range s {
			if unicode.IsLetter(r) {
				m[r] = append(m[r], pos)
			}
		}
		idx[i] = m
	}
	return idx
}

// ——————— Calcular min y max en slice ———————

func minMax(xs []int) (int, int) {
	lo, hi := math.MaxInt32, -1
	for _, v := range xs {
		if v < lo {
			lo = v
		}
		if v > hi {
			hi = v
		}
	}
	return lo, hi
}

// ——————— Generación recursiva de patrones ———————

func generatePatterns(seqs []string) []Pattern {
	ref := seqs[0]
	n := len(ref)
	anchors := extractAnchors(ref)
	indexes := buildIndex(seqs)

	seen := make(map[string]struct{})
	var out []Pattern

	var dfs func(start int, current []Anchor)
	dfs = func(start int, current []Anchor) {
		for i := start; i < len(anchors); i++ {
			next := append(current, anchors[i])
			if len(next) >= 2 {
				if pat, ok := computePattern(next, indexes, n); ok {
					if _, ex := seen[pat.Seq]; !ex {
						seen[pat.Seq] = struct{}{}
						out = append(out, pat)
					}
				}
			}
			dfs(i+1, next)
		}
	}
	dfs(0, nil)

	// Orden descendente por Score
	sort.Slice(out, func(i, j int) bool {
		return out[i].Score > out[j].Score
	})
	return out
}

// ——————— Cálculo de un patrón dado un slice de anclas ———————

func computePattern(anchors []Anchor, idx []map[rune][]int, n int) (Pattern, bool) {
	segments := len(anchors) - 1
	// Para cada segmento guardamos minGlobal y maxGlobal
	minGlob := make([]int, segments)
	maxGlob := make([]int, segments)
	for s := 0; s < segments; s++ {
		minGlob[s] = math.MaxInt32
		maxGlob[s] = -1
	}

	// Para cada secuencia, calculamos min/max gaps en cada segmento
	for si := 0; si < len(idx); si++ {
		for s := 0; s < segments; s++ {
			letterA := anchors[s].Letter
			letterB := anchors[s+1].Letter
			posA := idx[si][letterA]
			posB := idx[si][letterB]

			var local []int
			for _, pA := range posA {
				for _, pB := range posB {
					if pB > pA {
						local = append(local, pB-pA-1)
					}
				}
			}
			if len(local) == 0 {
				return Pattern{}, false
			}
			lMin, lMax := minMax(local)
			if lMin < minGlob[s] {
				minGlob[s] = lMin
			}
			if lMax > maxGlob[s] {
				maxGlob[s] = lMax
			}
		}
	}

	// Montamos el string ProSite
	parts := []string{}
	for i, a := range anchors {
		parts = append(parts, string(a.Letter))
		if i < segments {
			lo, hi := minGlob[i], maxGlob[i]
			if lo == hi {
				parts = append(parts, fmt.Sprintf("x(%d)", lo))
			} else {
				parts = append(parts, fmt.Sprintf("x(%d,%d)", lo, hi))
			}
		}
	}
	pat := strings.Join(parts, "-")
	score := scorePattern(pat, n)
	return Pattern{pat, score}, true
}

// ——————— Scoring: x + m·n + y·n² ———————

func scorePattern(p string, n int) int {
	parts := strings.Split(p, "-")
	sumX, m, y := 0, 0, 0

	for _, tok := range parts {
		if strings.HasPrefix(tok, "x(") {
			inner := tok[2 : len(tok)-1]
			if comma := strings.Index(inner, ","); comma >= 0 {
				// Rangos: tomar el mayor
				a, _ := strconv.Atoi(inner[:comma])
				b, _ := strconv.Atoi(inner[comma+1:])
				if a > b {
					sumX += a
				} else {
					sumX += b
				}
			} else {
				v, _ := strconv.Atoi(inner)
				sumX += v
			}
		} else if len(tok) == 1 {
			r := rune(tok[0])
			if unicode.IsUpper(r) {
				y++
			} else {
				m++
			}
		}
	}

	return sumX + m*n + y*(n*n)
}
