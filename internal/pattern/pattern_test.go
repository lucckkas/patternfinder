package pattern

import "testing"

func TestBestCommonPattern(t *testing.T) {
	seqs := []string{
		"asAfdBasdAdsC",
		"AsdGsBC",
		"AbsdfBdsBasdC",
		"AsdfBsadC",
	}
	agg, ok := BestCommonPattern(seqs)
	if !ok {
		t.Fatalf("se esperaba patrón común")
	}
	if agg.Pattern != "ABC" {
		t.Fatalf("patrón mayúsculas esperado ABC, got %s", agg.Pattern)
	}
	// comprobación básica del formateo (rangos específicos dependen de gaps exactos)
	form := FormatPattern(agg)
	if len(form) == 0 {
		t.Fatalf("formato vacío")
	}
}
