package parser

import (
	"testing"
)

var numericTests = []struct {
	Numeric      string
	ExpectString string
}{
	{
		Numeric:      RPL_WELCOME,
		ExpectString: "RPL_WELCOME",
	},
}

func TestNumeric(t *testing.T) {
	for idx, test := range numericTests {
		n := NewNumeric(test.Numeric)
		if got, want := n.Error(), test.ExpectString; got != want {
			t.Errorf("#%d: Error() = %s, want %s", idx, got, want)
		}
	}
}
