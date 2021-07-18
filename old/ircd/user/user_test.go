package user

import (
	"testing"

	"github.com/kylelemons/ircd-blight/old/ircd/parser"
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

var dummyNick = "[Dummy]"
var dummyLower = "{dummy}"
var dummyMixed = "{dUMmy]"
var nickSetTests = []struct {
	Nick  string
	Error error
	Count int
	After string
}{
	{
		Nick:  "14029804",
		Error: parser.NewNumeric(parser.ERR_ERRONEUSNICKNAME, ""),
		Count: 1,
		After: "*",
	},
	{
		Nick:  dummyNick,
		Error: parser.NewNumeric(parser.ERR_NICKNAMEINUSE, ""),
		Count: 1,
		After: "*",
	},
	{
		Nick:  "Nickname",
		Error: nil,
		Count: 2,
		After: "Nickname",
	},
	{
		Nick:  dummyNick,
		Error: parser.NewNumeric(parser.ERR_NICKNAMEINUSE, ""),
		Count: 2,
		After: "Nickname",
	},
	{
		Nick:  dummyMixed,
		Error: parser.NewNumeric(parser.ERR_NICKNAMEINUSE, ""),
		Count: 2,
		After: "Nickname",
	},
	{
		Nick:  "@Nick",
		Error: parser.NewNumeric(parser.ERR_ERRONEUSNICKNAME, ""),
		Count: 2,
		After: "Nickname",
	},
	{
		Nick:  "NewNick",
		Error: nil,
		Count: 2,
		After: "NewNick",
	},
}

func TestSetNick(t *testing.T) {
	// Set up the dummy user
	dummy := Get(NextUserID())
	err := dummy.SetNick(dummyNick)
	if err != nil {
		t.Fatalf("dummy.SetNick(%s) returned %s", dummyNick, err)
	}
	if got, want := userNicks[dummyLower], dummy.ID(); got != want {
		t.Errorf("map[%q] = %q, want %q", dummyLower, got, want)
	}

	// Set up victim user
	victim := Get(NextUserID())

	for idx, test := range nickSetTests {
		err := victim.SetNick(test.Nick)
		if got, want := err, test.Error; got != want && got.Error() != want.Error() {
			t.Errorf("#%d: SetNick(%q) = %#v, want %#v", idx, test.Nick, got, want)
		}
		if got, want := len(userNicks), test.Count; got != want {
			t.Errorf("#%d: len(userNicks) = %d, want %d", idx, got, want)
		}
		if got, want := victim.Nick(), test.After; got != want {
			t.Errorf("#%d: nick after = %q, want %q", idx, got, want)
		}
	}
}

func BenchmarkGenIDs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		<-userIDs
	}
}
