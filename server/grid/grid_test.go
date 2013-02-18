package grid

import (
	"math/rand"
	"reflect"
	"testing"
)

func TestGridOps(t *testing.T) {
	type insertion struct {
		first, second string
		notify        [2][]string // first element will be the inserted one
		inserted      bool
	}
	type deletion struct {
		first, second string
		notify        [2][]string
		deleted       bool
	}
	type quit struct {
		user    string
		parted  []string
		notify  [][]string
		deleted bool
	}

	tests := []struct {
		desc   string
		ops    []interface{}
		result [2]map[string][]string
	}{
		{
			desc: "inserts",
			ops: []interface{}{
				insertion{"#chat", "oper", [2][]string{{"#chat"}, {"oper"}}, true},
				insertion{"#opers", "oper", [2][]string{{"#chat", "#opers"}, {"oper"}}, true},
				insertion{"#chat", "user", [2][]string{{"#chat"}, {"oper", "user"}}, true},
				insertion{"#acro", "user", [2][]string{{"#acro", "#chat"}, {"user"}}, true},
			},
			result: [2]map[string][]string{
				{
					"#chat":  {"oper", "user"},
					"#opers": {"oper"},
					"#acro":  {"user"},
				},
				{
					"oper": {"#chat", "#opers"},
					"user": {"#acro", "#chat"},
				},
			},
		},
		{
			desc: "dup insert",
			ops: []interface{}{
				insertion{"#chat", "oper", [2][]string{{"#chat"}, {"oper"}}, true},
				insertion{"#opers", "oper", [2][]string{{"#chat", "#opers"}, {"oper"}}, true},
				insertion{"#chat", "oper", [2][]string{}, false},
			},
			result: [2]map[string][]string{
				{
					"#chat":  {"oper"},
					"#opers": {"oper"},
				},
				{
					"oper": {"#chat", "#opers"},
				},
			},
		},
		{
			desc: "basic delete",
			ops: []interface{}{
				insertion{"#chat", "user", [2][]string{{"#chat"}, {"user"}}, true},
				deletion{"#chat", "user", [2][]string{{"#chat"}, {"user"}}, true},
				insertion{"#chat", "user", [2][]string{{"#chat"}, {"user"}}, true},
			},
			result: [2]map[string][]string{
				{
					"#chat": {"user"},
				},
				{
					"user": {"#chat"},
				},
			},
		},
		{
			desc: "deletes",
			ops: []interface{}{
				insertion{"#chat", "oper", [2][]string{{"#chat"}, {"oper"}}, true},
				deletion{"#chat", "fake", [2][]string{}, false},
				insertion{"#opers", "user", [2][]string{{"#opers"}, {"user"}}, true},
				deletion{"#fake", "user", [2][]string{}, false},
				insertion{"#opers", "oper", [2][]string{{"#chat", "#opers"}, {"oper", "user"}}, true},
				deletion{"#opers", "user", [2][]string{{"#opers"}, {"oper", "user"}}, true},
				insertion{"#chat", "user", [2][]string{{"#chat"}, {"oper", "user"}}, true},
				insertion{"#acro", "user", [2][]string{{"#acro", "#chat"}, {"user"}}, true},
				deletion{"#fake", "fake", [2][]string{}, false},
			},
			result: [2]map[string][]string{
				{
					"#chat":  {"oper", "user"},
					"#opers": {"oper"},
					"#acro":  {"user"},
				},
				{
					"oper": {"#chat", "#opers"},
					"user": {"#acro", "#chat"},
				},
			},
		},
		{
			desc: "row col delete",
			ops: []interface{}{
				insertion{"#chat", "oper", [2][]string{{"#chat"}, {"oper"}}, true},
				insertion{"#opers", "oper", [2][]string{{"#chat", "#opers"}, {"oper"}}, true},
				insertion{"#opers", "user", [2][]string{{"#opers"}, {"oper", "user"}}, true},
				insertion{"#chat", "user", [2][]string{{"#chat", "#opers"}, {"oper", "user"}}, true},
				insertion{"#acro", "user", [2][]string{{"#acro", "#chat", "#opers"}, {"user"}}, true},
				quit{
					"user",
					[]string{"#acro", "#chat", "#opers"},
					[][]string{{"user"}, {"oper", "user"}, {"oper", "user"}},
					true,
				},
				quit{"other", nil, nil, false},
			},
			result: [2]map[string][]string{
				{
					"#chat":  {"oper"},
					"#opers": {"oper"},
				},
				{
					"oper": {"#chat", "#opers"},
				},
			},
		},
	}

	for _, test := range tests {
		var g Grid
		for idx, op := range test.ops {
			switch op := op.(type) {
			case insertion:
				pair := [2]string{op.first, op.second}
				notify, ok := g.Insert(pair, nil)
				if got, want := ok, op.inserted; got != want {
					t.Errorf("%s.%d: insert(%q).ok = %v, want %v", test.desc, idx, pair, got, want)
				}
				if got, want := notify, op.notify; !reflect.DeepEqual(got, want) {
					t.Errorf("%s.%d: insert(%q).notify = %v, want %v", test.desc, idx, pair, got, want)
				}
			case deletion:
				pair := [2]string{op.first, op.second}
				notify, ok := g.Delete(pair)
				if got, want := ok, op.deleted; got != want {
					t.Errorf("%s.%d: delete(%q).ok = %v, want %v", test.desc, idx, pair, got, want)
				}
				if got, want := notify, op.notify; !reflect.DeepEqual(got, want) {
					t.Errorf("%s.%d: delete(%q).notify = %v, want %v", test.desc, idx, pair, got, want)
				}
			case quit:
				// users are the second edge
				deleted, notify, ok := g.DeleteAll(1, op.user)
				if got, want := ok, op.deleted; got != want {
					t.Errorf("%s.%d: deleteall(%q).deleted = %v, want %v", test.desc, idx, op.user, got, want)
				}
				if got, want := notify, op.notify; !reflect.DeepEqual(got, want) {
					t.Errorf("%s.%d: deleteall(%q).notify = %v, want %v", test.desc, idx, op.user, got, want)
				}
				if got, want := deleted, op.parted; !reflect.DeepEqual(got, want) {
					t.Errorf("%s.%d: deleteall(%q).list = %v, want %v", test.desc, idx, op.user, got, want)
				}
			default:
				t.Errorf("%s.%d: unknown %#v", test.desc, idx, op)
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

func insertCount(cnt int) Grid {
	prng := rand.New(rand.NewSource(int64(cnt)))

	var g Grid
	for i := 0; i < cnt; i++ {
		c, u := prng.Intn(len(channels)), prng.Intn(len(users))
		g.Insert([2]string{channels[c], users[u]}, nil)
	}
	return g
}

func BenchmarkGridInsert(b *testing.B) {
	const batch = 10000
	for i := 0; i < b.N; {
		insertCount(batch)
		i += batch
	}
}

func getCount(g Grid, cnt int) {
	prng := rand.New(rand.NewSource(int64(cnt)))

	for i := 0; i < cnt; i++ {
		c, u := prng.Intn(len(channels)), prng.Intn(len(users))
		g.Get([2]string{channels[c], users[u]})
	}
}

func BenchmarkGridGet(b *testing.B) {
	const start = 5000
	const batch = 10000

	g := insertCount(start)

	b.ResetTimer()
	for i := 0; i < b.N; {
		getCount(g, batch)
		i += batch
	}
}

func insDelCount(g Grid, cnt int) {
	prng := rand.New(rand.NewSource(int64(cnt)))

	for i := 0; i < cnt; i++ {
		c, u := prng.Intn(len(channels)), prng.Intn(len(users))
		g.Insert([2]string{channels[c], users[u]}, nil)
		g.Delete([2]string{channels[c], users[u]})
	}
}

func BenchmarkGridInsDel(b *testing.B) {
	const start = 5000
	const batch = 10000

	g := insertCount(start)

	b.ResetTimer()
	for i := 0; i < b.N; {
		insDelCount(g, batch)
		i += batch
	}
}
