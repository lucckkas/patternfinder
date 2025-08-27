package pattern

// IntRange representa un rango [Min, Max] de minúsculas entre dos mayúsculas consecutivas.
type IntRange struct {
	Min int
	Max int
}

// UpperProjection contiene la proyección a mayúsculas y ayudas de conteo de minúsculas.
type UpperProjection struct {
	Original  string // string original
	UpperSeq  string // sólo mayúsculas extraídas del original
	UpperPos  []int  // posiciones (índices) en Original donde están las mayúsculas
	PrefLower []int  // prefijo: cantidad de minúsculas hasta cada índice (len= len(Original)+1)
}

// AggregatedPattern resume el patrón final y sus métricas agregadas.
type AggregatedPattern struct {
	Pattern      string     // p.ej. "ABC"
	GapRanges    []IntRange // para cada gap entre letras del patrón
	GapAvg       []float64  // promedio de minúsculas por gap (opcional, informativo)
	ScoreUpper   int        // cantidad de mayúsculas (prioridad 1)
	ScoreGapsSum int        // suma de mínimos por gap (prioridad 2)
}
