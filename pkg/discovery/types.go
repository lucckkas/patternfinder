package discovery

// token representa ya sea un bloque de gap o una letra.
type token struct {
	isNum  bool
	num    int
	letter rune
}
