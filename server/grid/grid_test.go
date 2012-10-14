package grid

import (
	"testing"
	"math/rand"
	"reflect"
)

func TestGridOps(t *testing.T) {
	type insertion struct {
		first, second string
		inserted bool
	}
	type deletion struct {
		first, second string
		deleted bool
	}

	tests := []struct{
		desc   string
		ops    []interface{}
		result [2]map[string][]string
	}{
		{
			desc: "inserts",
			ops: []interface{}{
				insertion{"#chat", "oper", true},
				insertion{"#opers", "oper", true},
				insertion{"#chat", "user", true},
				insertion{"#acro", "user", true},
			},
			result: [2]map[string][]string{
				{
					"#chat": []string{"oper", "user"},
					"#opers": []string{"oper"},
					"#acro": []string{"user"},
				},
				{
					"oper": []string{"#chat", "#opers"},
					"user": []string{"#acro", "#chat"},
				},
			},
		},
		{
			desc: "dup insert",
			ops: []interface{}{
				insertion{"#chat", "oper", true},
				insertion{"#opers", "oper", true},
				insertion{"#chat", "oper", false},
			},
			result: [2]map[string][]string{
				{
					"#chat": []string{"oper"},
					"#opers": []string{"oper"},
				},
				{
					"oper": []string{"#chat", "#opers"},
				},
			},
		},
		{
			desc: "basic delete",
			ops: []interface{}{
				insertion{"#chat", "user", true},
				deletion{"#chat", "user", true},
				insertion{"#chat", "user", true},
			},
			result: [2]map[string][]string{
				{
					"#chat": []string{"user"},
				},
				{
					"user": []string{"#chat"},
				},
			},
		},
		{
			desc: "deletes",
			ops: []interface{}{
				insertion{"#chat", "oper", true},
				deletion{"#chat", "fake", false},
				insertion{"#opers", "user", true},
				deletion{"#fake", "user", false},
				insertion{"#opers", "oper", true},
				deletion{"#opers", "user", true},
				insertion{"#chat", "user", true},
				insertion{"#acro", "user", true},
				deletion{"#fake", "fake", false},
			},
			result: [2]map[string][]string{
				{
					"#chat": []string{"oper", "user"},
					"#opers": []string{"oper"},
					"#acro": []string{"user"},
				},
				{
					"oper": []string{"#chat", "#opers"},
					"user": []string{"#acro", "#chat"},
				},
			},
		},
	}

	for _, test := range tests {
		var g Grid
		for _, op := range test.ops {
			switch op := op.(type) {
			case insertion:
				pair := [2]string{op.first, op.second}
				if got, want := g.Insert(pair), op.inserted; got != want {
					t.Errorf("%s: insert(%q) = %v, want %v", test.desc, pair, got, want)
				}
			case deletion:
				pair := [2]string{op.first, op.second}
				if got, want := g.Delete(pair), op.deleted; got != want {
					t.Errorf("%s: delete(%q) = %v, want %v", test.desc, pair, got, want)
				}
			}
		}
		if got, want := g.dump(), test.result; !reflect.DeepEqual(got, want) {
			t.Errorf("%s: got  %q", test.desc, got)
			t.Errorf("%s: want %q", test.desc, want)
		}
	}
}

// TODO(kevlar): Examples for Grid

var (
	channels = [...]string{
		"#opers", "#chat", "#users", "#irc", "#help", "#golang",
		"##programming", "##IRC", "#testing", "#42",
	}
	users = [...]string{
		"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "golf",
		"hotel", "india", "juliet", "kilo", "lima", "mike", "november",
		"oscar", "papa", "quebec", "romeo", "sierra", "tango", "uniform",
		"victor", "whiskey", "x-ray", "yankee", "zulu",
	}
)

func insertCount(cnt int) {
	prng := rand.New(rand.NewSource(int64(cnt)))

	var g Grid
	for i := 0; i < cnt; i++ {
		c, u := prng.Intn(len(channels)), prng.Intn(len(users))
		g.Insert([2]string{channels[c], users[u]})
	}
}

func BenchmarkGridInsert(b *testing.B) {
	const batch = 10000
	for i := 0; i < b.N; {
		insertCount(batch)
		i += batch
	}
}
