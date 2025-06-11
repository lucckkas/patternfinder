package discovery

// token representa ya sea un bloque de gap o una letra.
type token struct {
	isNum  bool
	num    int
	letter rune
}

// helpers genéricos
func isDigit(r rune) bool { return r >= '0' && r <= '9' }
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
