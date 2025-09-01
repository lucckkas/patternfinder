package pattern

import "strings"

func earliestPositions(pattern, s string) ([]int, bool) {
	pos := make([]int, len(pattern))
	start := 0
	for i := 0; i < len(pattern); i++ {
		ch := pattern[i]
		off := strings.IndexByte(s[start:], ch)
		if off < 0 {
			return nil, false
		}
		pos[i] = start + off
		start = pos[i] + 1
	}
	return pos, true
}

func latestPositions(pattern, s string) ([]int, bool) {
	pos := make([]int, len(pattern))
	end := len(s)
	for i := len(pattern) - 1; i >= 0; i-- {
		ch := pattern[i]
		last := -1
		for j := end - 1; j >= 0; j-- {
			if s[j] == ch {
				last = j
				break
			}
		}
		if last < 0 {
			return nil, false
		}
		pos[i] = last
		end = last
	}
	return pos, true
}

func indicesOf(s string, ch byte, lo, hi int) []int {
	if lo < 0 {
		lo = 0
	}
	if hi >= len(s) {
		hi = len(s) - 1
	}
	if lo > hi {
		return nil
	}
	out := make([]int, 0)
	i := lo
	for {
		k := strings.IndexByte(s[i:hi+1], ch)
		if k < 0 {
			break
		}
		idx := i + k
		out = append(out, idx)
		i = idx + 1
	}
	return out
}

// minGapForPair: MÍNIMO número de minúsculas entre pattern[i] y pattern[i+1]
// en UNA secuencia, considerando todas las incrustaciones viables (a<b).
func minGapForPair(pattern string, proj UpperProjection, i int, earliest, latest []int) (int, bool) {
	s := proj.Original

	loA, hiA := earliest[i], latest[i]
	loB, hiB := earliest[i+1], latest[i+1]

	occA := indicesOf(s, pattern[i], loA, hiA)
	occB := indicesOf(s, pattern[i+1], loB, hiB)
	if len(occA) == 0 || len(occB) == 0 {
		return 0, false
	}

	minGap := int(^uint(0) >> 1) // +Inf
	j := 0
	for _, a := range occA {
		for j < len(occB) && occB[j] <= a {
			j++
		}
		if j == len(occB) {
			break
		}
		b := occB[j]
		val := b - a - 1
		if val < minGap {
			minGap = val
		}
	}
	if minGap == int(^uint(0)>>1) {
		return 0, false
	}
	return minGap, true
}

// minGapsForPatternPerSequence: para UNA secuencia, devuelve el MINIMO gap por cada par del patrón.
func minGapsForPatternPerSequence(pattern string, proj UpperProjection) ([]int, bool) {
	if len(pattern) <= 1 {
		return []int{}, true
	}
	ear, ok := earliestPositions(pattern, proj.Original)
	if !ok {
		return nil, false
	}
	lat, ok := latestPositions(pattern, proj.Original)
	if !ok {
		return nil, false
	}
	out := make([]int, 0, len(pattern)-1)
	for i := 0; i+1 < len(pattern); i++ {
		g, ok := minGapForPair(pattern, proj, i, ear, lat)
		if !ok {
			return nil, false
		}
		out = append(out, g)
	}
	return out, true
}
