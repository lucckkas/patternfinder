package pattern

// ProjectUpper genera la proyección a mayúsculas y el prefijo de minúsculas.
// Asume entradas ASCII [A-Z][a-z]; si necesitas Unicode pleno, conviene trabajar con runas.
func ProjectUpper(s string) UpperProjection {
	upperPos := make([]int, 0, len(s))
	upperSeq := make([]byte, 0, len(s))
	prefLower := make([]int, len(s)+1)

	for i := 0; i < len(s); i++ {
		c := s[i]
		isLower := (c >= 'a' && c <= 'z')
		if isLower {
			prefLower[i+1] = prefLower[i] + 1
		} else {
			prefLower[i+1] = prefLower[i]
		}
		if c >= 'A' && c <= 'Z' {
			upperPos = append(upperPos, i)
			upperSeq = append(upperSeq, c)
		}
	}
	return UpperProjection{
		Original:  s,
		UpperSeq:  string(upperSeq),
		UpperPos:  upperPos,
		PrefLower: prefLower,
	}
}
