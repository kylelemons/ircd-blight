package server

import (
	"errors"
	"reflect"
	"testing"

	"github.com/kylelemons/ircd-blight/server/data"
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
		u, err := s.signon(test.Nick, test.Nick, test.Name)
		if !reflect.DeepEqual(err, test.Error) {
			t.Errorf("signon(%q, %q, %q): %v, want %v",
				test.Nick, test.Nick, test.Name,
				err, test.Error)
		}
		if err != nil {
			continue
		}
		if u.UID != test.UID {
			t.Errorf("signon(%q, %q, %q) = %q, want %q",
				test.Nick, test.Nick, test.Name,
				u.UID, test.UID)
		}
	}
}

func TestJoin(t *testing.T) {
	s := NewServer("7ST")

	user := func(name string) string {
		u, err := s.signon(name, name, name)
		if err != nil {
			t.Fatalf("signon(%q): %s", err)
		}
		return u.UID
	}

	var (
		zaphod = user("zaphod")
		ford   = user("ford")
	)

	tests := []struct {
		UID, Channel string
		Mode         data.MemberMode
		Error        error
	}{
		{
			UID:     zaphod,
			Channel: "#HoG",
			Mode:    data.MemberOp | data.MemberAdmin,
		},
		{
			UID:     ford,
			Channel: "#HoG",
		},
		{
			UID:     ford,
			Channel: "#HoG",
			Error:   errors.New(`UID "7STAAAAAB" is already on #HoG`),
		},
	}

	for _, test := range tests {
		m, err := s.join(test.UID, test.Channel)
		if !reflect.DeepEqual(err, test.Error) {
			t.Errorf("join(%q, %q): %v, want %v",
				test.UID, test.Channel,
				err, test.Error)
		}
		if err != nil {
			continue
		}
		if got, want := m.Mode, test.Mode; got != want {
			t.Errorf("join(%q, %q).mode = %08b, want %08b",
				test.UID, test.Channel,
				got, want)
		}
	}
}

func BenchmarkIDStr(b *testing.B) {
	N := uint64(b.N)
	for i := uint64(0); i < N; i++ {
		idstr(i, UserIDLen)
	}
}
