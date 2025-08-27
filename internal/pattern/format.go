package pattern

import "fmt"

// FormatPattern convierte el patrón y sus rangos a la notación A-x(n)-B-x(n,m)-C.
func FormatPattern(agg AggregatedPattern) string {
	if len(agg.Pattern) == 0 {
		return ""
	}
	out := make([]byte, 0, len(agg.Pattern)*4)
	for i := 0; i < len(agg.Pattern); i++ {
		out = append(out, agg.Pattern[i])
		if i+1 < len(agg.Pattern) {
			r := agg.GapRanges[i]
			if r.Min == r.Max {
				out = append(out, []byte(fmt.Sprintf("-x(%d)-", r.Min))...)
			} else {
				out = append(out, []byte(fmt.Sprintf("-x(%d,%d)-", r.Min, r.Max))...)
			}
		}
	}
	return string(out)
}
