package user

import (
	"testing"
)

var testIDs = map[int64]string{
	0:    "AAAAAA",
	1:    "AAAAAB",
	2:    "AAAAAC",
	25:   "AAAAAZ",
	26:   "AAAAA0",
	35:   "AAAAA9",
	36:   "AAAABA",
	1295: "AAAA99",
	1296: "AAABAA",
}

func TestGenIDs(t *testing.T) {
	gotCount := 0
	for i := int64(0); gotCount < len(testIDs); i++ {
		got := <-userIDs
		if want, ok := testIDs[i]; ok {
			gotCount++
			if got != want {
				t.Errorf("id[%d] = %q, got %q", i, got, want)
			}
		}
	}
}

func BenchmarkGenIDs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		<-userIDs
	}
}
