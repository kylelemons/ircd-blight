// +build gofuzzbeta

package mode

import (
	"bufio"
	"sort"
	"strings"
	"testing"
)

func FuzzModeMap_ParseModeChange(f *testing.F) {
	for _, test := range parseModeChangeTests {
		f.Add(test.ModeChange)
	}
	for _, test := range applyModeTests {
		f.Add(test.Apply)
		f.Add(test.Result)
	}
	f.Fuzz(func(t *testing.T, modeChange string) {
		parsed, err := ChannelModes.ParseModeChange(strings.Split(modeChange, " "))
		if err != nil {
			t.SkipNow()
		}
		_ = parsed // are there any sanity checks we can do?
	})
}

func FuzzActiveModes_Apply(f *testing.F) {
	for _, test := range parseModeChangeTests {
		f.Add(test.ModeChange)
	}
	var accum string
	for _, test := range applyModeTests {
		f.Add(test.Apply)
		accum += test.Apply + "\n"
		f.Add(accum)
	}
	f.Fuzz(func(t *testing.T, modeChanges string) {
		lines := bufio.NewScanner(strings.NewReader(modeChanges))
		state := NewActiveModes(ChannelModes)
		for lines.Scan() {
			parsed, err := ChannelModes.ParseModeChange(strings.Split(lines.Text(), " "))
			if err != nil {
				t.SkipNow()
			}
			if _, err := state.Apply(parsed); err != nil {
				t.Errorf("Failed to apply valid mode change %q: %s", lines.Text(), err)
				t.Errorf("... prior state: %q", state)
			}
		}
		if err := lines.Err(); err != nil {
			t.Fatalf("failed reading lines: %s", err)
		}
		if !sort.IsSorted(state.Modes) {
			t.Errorf("Modes unsorted after operations:")
			t.Errorf("OPERATIONS:\n%s", modeChanges)
			t.Errorf("MODES:\n%s", state.String())
		}
	})
}
