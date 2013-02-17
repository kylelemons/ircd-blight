package server

import (
	"errors"
	"reflect"
	"testing"
)

func TestSignon(t *testing.T) {
	s := NewServer("7ST")

	tests := []struct {
		Nick, Name string // input
		UID        string // output
		Error      error
	}{
		{
			Nick: "zaphod",
			Name: "Mr. President",
			UID:  "7STAAAAAA",
		},
		{
			Nick: "arthur",
			Name: "... what teacup?",
			UID:  "7STAAAAAB",
		},
		{
			Nick:  "zaphod",
			Name:  "impostor",
			Error: errors.New(`NICK "zaphod" already in use`),
		},
	}

	for _, test := range tests {
		uid, err := s.signon(test.Nick, test.Nick, test.Name)
		if !reflect.DeepEqual(err, test.Error) {
			t.Errorf("signon(%q, %q, %q): %v, want %v",
				test.Nick, test.Nick, test.Name,
				err, test.Error)
		}
		if err != nil {
			continue
		}
		if uid != test.UID {
			t.Errorf("signon(%q, %q, %q) = %q, want %q",
				test.Nick, test.Nick, test.Name,
				uid, test.UID)
		}
	}
}

func BenchmarkIDStr(b *testing.B) {
	N := uint64(b.N)
	for i := uint64(0); i < N; i++ {
		idstr(i, UserIDLen)
	}
}
