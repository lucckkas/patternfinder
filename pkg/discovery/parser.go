package discovery

import (
	"strconv"
	"unicode"
)

// parseCollapsed: "A10C5D" → [A][10][C][5][D]
func parseCollapsed(s string) []token {
	var toks []token
	for i := 0; i < len(s); {
		r := rune(s[i])
		if unicode.IsDigit(r) {
			j := i
			for j < len(s) && unicode.IsDigit(rune(s[j])) {
				j++
			}
			n, _ := strconv.Atoi(s[i:j])
			toks = append(toks, token{isNum: true, num: n})
			i = j
		} else {
			toks = append(toks, token{isNum: false, letter: r})
			i++
		}
	}
	return toks
}
