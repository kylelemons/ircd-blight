package mode

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

var parseModeChangeTests = []struct {
	ModeSet    *ModeMap
	ModeChange string
	Parsed     []Mode
	Errors     []error
}{
	{
		ModeSet:    ChannelModes,
		ModeChange: "+o User",
		Parsed: []Mode{
			{
				Index: ChannelModes.Index['o'],
				Op:    SetMode,
				Args:  []string{"User"},
			},
		},
	},
	{
		ModeSet:    ChannelModes,
		ModeChange: "b+mkolsv key Op 10 Voice",
		Parsed: []Mode{
			{
				Index: ChannelModes.Index['b'],
				Op:    QueryMode,
			},
			{
				Index: ChannelModes.Index['m'],
				Op:    SetMode,
			},
			{
				Index: ChannelModes.Index['k'],
				Op:    SetMode,
				Args:  []string{"key"},
			},
			{
				Index: ChannelModes.Index['o'],
				Op:    SetMode,
				Args:  []string{"Op"},
			},
			{
				Index: ChannelModes.Index['l'],
				Op:    SetMode,
				Args:  []string{"10"},
			},
			{
				Index: ChannelModes.Index['s'],
				Op:    SetMode,
			},
			{
				Index: ChannelModes.Index['v'],
				Op:    SetMode,
				Args:  []string{"Voice"},
			},
		},
	},
	{
		ModeSet:    ChannelModes,
		ModeChange: "+o-h+v NewOp OldHop NewVoice",
		Parsed: []Mode{
			{
				Index: ChannelModes.Index['o'],
				Op:    SetMode,
				Args:  []string{"NewOp"},
			},
			{
				Index: ChannelModes.Index['h'],
				Op:    UnsetMode,
				Args:  []string{"OldHop"},
			},
			{
				Index: ChannelModes.Index['v'],
				Op:    SetMode,
				Args:  []string{"NewVoice"},
			},
		},
	},
	{
		ModeSet:    UserModes,
		ModeChange: "-i+Zr",
		Parsed: []Mode{
			{
				Index: UserModes.Index['i'],
				Op:    UnsetMode,
			},
			{
				Index: UserModes.Index['Z'],
				Op:    SetMode,
			},
			{
				Index: UserModes.Index['r'],
				Op:    SetMode,
			},
		},
	},
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestParseModeChange(t *testing.T) {
	for idx, test := range parseModeChangeTests {
		parsed, _ := test.ModeSet.ParseModeChange(strings.Fields(test.ModeChange))
		if got, want := len(parsed), len(test.Parsed); got != want {
			t.Errorf("%d. len(parsed) = %d, want %d", idx, got, want)
		}
		for i := 0; i < min(len(parsed), len(test.Parsed)); i++ {
			got, want := parsed[i], test.Parsed[i]
			if got, want := got.Index, want.Index; got != want {
				t.Errorf("%d. mode[%d] is %q, want %q",
					idx, i, test.ModeSet.Modes[int(got)], test.ModeSet.Modes[int(want)])
			}
			if got, want := got.Op, want.Op; got != want {
				t.Errorf("%d. mode[%d].Op = %d, want %d", idx, i, got, want)
			}
			if got, want := got.Args, want.Args; !reflect.DeepEqual(got, want) {
				t.Errorf("%d. mode[%d].Args = %v, want %v", idx, i, got, want)
			}
		}
	}
}

func TestModeMap_ParseModeChange_goodModes(t *testing.T) {
	tests := []struct {
		name  string
		modes string
	}{
		{"ban", "+b some*!bad@user123.com"},
		{"letters", "+o abcdABCD"},
		{"numbers", "+o abcd1234"},
		{"punctuation", "+o abcd~`!@#$%^&*()_-+={}[]|\\:;\"'<,>.?/"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ChannelModes.ParseModeChange(strings.Split(test.modes, " ")); err != nil {
				t.Errorf("ParseModeChange(%q) failed: %s", test.modes, err)
			}
		})
	}
}

func TestModeMap_ParseModeChange_badModes(t *testing.T) {
	tests := []struct {
		name  string
		modes string
	}{
		{"bad", "\\rp-h\\t@OpHovse"},
		{"empty", " "},
		{"empty_arg", "+ooo a  b"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ChannelModes.ParseModeChange(strings.Split(test.modes, " ")); err != nil {
				return
			}
			t.Errorf("ParseModeChange(%q) succeeded unexpectedly", test.modes)
		})
	}
}

func TestModeString(t *testing.T) {
	for idx, test := range parseModeChangeTests {
		if got, want := test.ModeSet.ModeString(test.Parsed), test.ModeChange; got != want {
			t.Errorf("%d. ModeString = %q, want %q", idx, got, want)
		}
	}
}

var applyModeTests = []struct {
	Apply  string
	Result string
}{
	// Channel modes
	{
		Apply:  "+o O",
		Result: "+o O",
	},
	{
		Apply:  "+hv H V",
		Result: "+ohv O H V",
	},
	{
		Apply:  "-ohv O H H",
		Result: "+v V",
	},
	{
		Apply:  "+vv 1 2",
		Result: "+vvv V 1 2",
	},
	{
		Apply:  "+kl test 10",
		Result: "+vvvkl V 1 2 test 10",
	},
	{
		Apply:  "+b *!testuser@*",
		Result: "+vvvbkl V 1 2 *!testuser@* test 10",
	},
	{
		Apply:  "-b *!anotheruser@*",
		Result: "+vvvbkl V 1 2 *!testuser@* test 10",
	},
	{
		Apply:  "-lk foo",
		Result: "+vvvb V 1 2 *!testuser@*",
	},
	{
		Apply:  "-b *!testuser@*",
		Result: "+vvv V 1 2",
	},
	{
		Apply:  "-vvv 2 1 V",
		Result: "",
	},
	// end result should be "" for benchmark to work
}

func TestApplyModes(t *testing.T) {
	var prev []string
	running := NewActiveModes(ChannelModes)
	for _, test := range applyModeTests {
		t.Run("apply_"+test.Apply, func(t *testing.T) {
			parsed, errs := ChannelModes.ParseModeChange(strings.Fields(test.Apply))
			if len(errs) > 0 {
				t.Errorf("ParseModeChange(%q):", test.Apply)
				for _, err := range errs {
					t.Errorf(" - %s", err)
				}
				t.FailNow()
			}

			applied, errs := running.Apply(parsed)
			if len(errs) > 0 {
				t.Errorf("Apply(%#v):", parsed)
				for _, err := range errs {
					t.Errorf(" - %s", err)
				}
				t.FailNow()
			}

			if got, want := ChannelModes.ModeString(applied), test.Apply; got != want {
				t.Errorf("After %q, applying %q, ModeString = %q, want %q", prev, test.Apply, got, want)
			}
			if !sort.IsSorted(running.Modes) {
				t.Errorf("Apply de-sorted modes:")
				for i, m := range running.Modes {
					t.Errorf("modes[%d] = %#v", i, m)
				}
			}
			if got, want := running.String(), test.Result; got != want {
				t.Errorf("After %q, applying %q, result = %q, want %q", prev, test.Apply, got, want)
				for i, m := range running.Modes {
					t.Errorf("modes[%d] = %#v", i, m)
				}
			}
			prev = append(prev, test.Apply)
		})
	}
}

func BenchmarkModeMap_ParseModeChange(b *testing.B) {
	var argsets [][]string
	for _, test := range parseModeChangeTests {
		argsets = append(argsets, strings.Fields(test.ModeChange))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % len(parseModeChangeTests)
		test := parseModeChangeTests[index]
		args := argsets[index]
		test.ModeSet.ParseModeChange(args)
	}
}

func BenchmarkActiveModes_Apply(b *testing.B) {
	var changesets [][]Mode
	for _, test := range applyModeTests {
		parsed, errs := ChannelModes.ParseModeChange(strings.Fields(test.Apply))
		if len(errs) > 0 {
			b.Errorf("ParseModeChange(%q):", test.Apply)
			for _, err := range errs {
				b.Errorf(" - %s", err)
			}
			b.FailNow()
		}
		changesets = append(changesets, parsed)
	}

	running := NewActiveModes(ChannelModes)

	b.ReportAllocs()
	b.ResetTimer()
	for n := b.N; n > 0; n -= len(changesets) {
		for _, changes := range changesets {
			running.Apply(changes)
		}
	}
	if after := running.String(); len(after) != 0 {
		b.Errorf("Applies did not end up at empty: %q", after)
	}
}
