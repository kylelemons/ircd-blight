package mode

import (
	"reflect"
	"strings"
	"testing"
)

var parseModeChangeTests = []struct {
	ModeSet    ModeMap
	ModeChange string
	Parsed     []Mode
	Errors     []error
}{
	{
		ModeSet:    ChannelModes,
		ModeChange: "+o User",
		Parsed: []Mode{
			{
				Spec: ChannelModes['o'],
				Op:   SetMode,
				Args: []string{"User"},
			},
		},
	},
	{
		ModeSet:    ChannelModes,
		ModeChange: "b+mkolsv key Op 10 Voice",
		Parsed: []Mode{
			{
				Spec: ChannelModes['b'],
				Op:   QueryMode,
			},
			{
				Spec: ChannelModes['m'],
				Op:   SetMode,
			},
			{
				Spec: ChannelModes['k'],
				Op:   SetMode,
				Args: []string{"key"},
			},
			{
				Spec: ChannelModes['o'],
				Op:   SetMode,
				Args: []string{"Op"},
			},
			{
				Spec: ChannelModes['l'],
				Op:   SetMode,
				Args: []string{"10"},
			},
			{
				Spec: ChannelModes['s'],
				Op:   SetMode,
			},
			{
				Spec: ChannelModes['v'],
				Op:   SetMode,
				Args: []string{"Voice"},
			},
		},
	},
	{
		ModeSet:    ChannelModes,
		ModeChange: "+o-h+v NewOp OldHop NewVoice",
		Parsed: []Mode{
			{
				Spec: ChannelModes['o'],
				Op:   SetMode,
				Args: []string{"NewOp"},
			},
			{
				Spec: ChannelModes['h'],
				Op:   UnsetMode,
				Args: []string{"OldHop"},
			},
			{
				Spec: ChannelModes['v'],
				Op:   SetMode,
				Args: []string{"NewVoice"},
			},
		},
	},
	{
		ModeSet:    UserModes,
		ModeChange: "-i+Zr",
		Parsed: []Mode{
			{
				Spec: UserModes['i'],
				Op:   UnsetMode,
			},
			{
				Spec: UserModes['Z'],
				Op:   SetMode,
			},
			{
				Spec: UserModes['r'],
				Op:   SetMode,
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
		parsed, _ := ParseModeChange(strings.Fields(test.ModeChange), test.ModeSet)
		if got, want := len(parsed), len(test.Parsed); got != want {
			t.Errorf("%d. len(parsed) = %d, want %d", idx, got, want)
		}
		for i := 0; i < min(len(parsed), len(test.Parsed)); i++ {
			got, want := parsed[i], test.Parsed[i]
			if got, want := got.Spec, want.Spec; got != want {
				t.Errorf("%d. mode[%d] is %q, want %q", idx, i, got.Name(), want.Name())
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

func TestModeString(t *testing.T) {
	for idx, test := range parseModeChangeTests {
		if got, want := ModeString(test.Parsed), test.ModeChange; got != want {
			t.Errorf("%d. ModeString = %q, want %q", idx, got, want)
		}
	}
}

var applyModeTests = []struct {
	Apply   string
	Applied string
	Result  string
}{
	// Channel modes
	{
		Apply:   "+o O",
		Applied: "+o O",
		Result:  "+o O",
	},
	{
		Apply:   "+hv H V",
		Applied: "+hv H V",
		Result:  "+voh V O H",
	},
	{
		Apply:   "-ohv O H H",
		Applied: "-ohv O H H",
		Result:  "+v V",
	},
	{
		Apply:   "+vv 1 2",
		Applied: "+vv 1 2",
		Result:  "+vvv V 1 2",
	},
}

func TestApplyModes(t *testing.T) {
	running := make(ActiveModes)
	for idx, test := range applyModeTests {
		parsed, _ := ParseModeChange(strings.Fields(test.Apply), ChannelModes)
		applied, _ := running.Apply(parsed)
		if got, want := ModeString(applied), test.Applied; got != want {
			t.Errorf("%d. applied = %q, want %q", idx, got, want)
		}
		if got, want := running.String(), test.Result; got != want {
			t.Errorf("%d. result = %q, want %q", idx, got, want)
		}
	}
}
