package server

import (
	"testing"
)

type tOp int

const (
	tConnect tOp = iota
	tLink
	tUnlink
)

var linkingTests = []struct {
	Op  tOp
	SID string

	// for tLink
	UpSID string
	Hops  string

	Map map[string]map[string]bool
}{
	{ // 0
		Op:  tConnect,
		SID: "0AA",
		Map: map[string]map[string]bool{
			"0AA": {},
		},
	},
	{ // 1
		Op:  tConnect,
		SID: "1AA",
		Map: map[string]map[string]bool{
			"0AA": {},
			"1AA": {},
		},
	},
	{ // 2
		Op:    tLink,
		SID:   "00A",
		UpSID: "0AA",
		Hops:  "2",
		Map: map[string]map[string]bool{
			"0AA": {"00A": true},
			"1AA": {},
			"00A": {},
		},
	},
	{ // 3
		Op:    tLink,
		SID:   "000",
		UpSID: "00A",
		Hops:  "3",
		Map: map[string]map[string]bool{
			"0AA": {"00A": true},
			"00A": {"000": true},
			"1AA": {},
			"000": {},
		},
	},
	{ // 4
		Op:    tLink,
		SID:   "00B",
		UpSID: "0AA",
		Hops:  "2",
		Map: map[string]map[string]bool{
			"0AA": {"00A": true, "00B": true},
			"00A": {"000": true},
			"1AA": {},
			"00B": {},
			"000": {},
		},
	},
	{ // 5
		Op:    tLink,
		SID:   "10A",
		UpSID: "1AA",
		Hops:  "2",
		Map: map[string]map[string]bool{
			"0AA": {"00A": true, "00B": true},
			"00A": {"000": true},
			"1AA": {"10A": true},
			"00B": {},
			"000": {},
			"10A": {},
		},
	},
	{ // 6
		Op:    tLink,
		SID:   "100",
		UpSID: "10A",
		Hops:  "3",
		Map: map[string]map[string]bool{
			"0AA": {"00A": true, "00B": true},
			"00A": {"000": true},
			"1AA": {"10A": true},
			"10A": {"100": true},
			"00B": {},
			"000": {},
			"100": {},
		},
	},
	{ // 7
		Op:  tConnect,
		SID: "2AA",
		Map: map[string]map[string]bool{
			"0AA": {"00A": true, "00B": true},
			"1AA": {"10A": true},
			"2AA": {},
			"00A": {"000": true},
			"00B": {},
			"10A": {"100": true},
			"000": {},
			"100": {},
		},
	},
	{ // 8
		Op:    tLink,
		SID:   "20A",
		UpSID: "2AA",
		Hops:  "2",
		Map: map[string]map[string]bool{
			"0AA": {"00A": true, "00B": true},
			"1AA": {"10A": true},
			"2AA": {"20A": true},
			"00A": {"000": true},
			"00B": {},
			"10A": {"100": true},
			"20A": {},
			"000": {},
			"100": {},
		},
	},
	{ // 9
		Op:    tLink,
		SID:   "101",
		UpSID: "10A",
		Hops:  "2",
		Map: map[string]map[string]bool{
			"0AA": {"00A": true, "00B": true},
			"00A": {"000": true},
			"000": {},
			"00B": {},

			"1AA": {"10A": true},
			"10A": {"100": true, "101": true},
			"100": {},
			"101": {},

			"2AA": {"20A": true},
			"20A": {},
		},
	},
	{ // 10
		Op:    tLink,
		SID:   "102",
		UpSID: "101",
		Hops:  "2",
		Map: map[string]map[string]bool{
			"0AA": {"00A": true, "00B": true},
			"00A": {"000": true},
			"000": {},
			"00B": {},

			"1AA": {"10A": true},
			"10A": {"100": true, "101": true},
			"100": {},
			"101": {"102": true},
			"102": {},

			"2AA": {"20A": true},
			"20A": {},
		},
	},
	{ // 11
		Op:  tUnlink,
		SID: "2AA",
		Map: map[string]map[string]bool{
			"0AA": {"00A": true, "00B": true},
			"00A": {"000": true},
			"000": {},
			"00B": {},

			"1AA": {"10A": true},
			"10A": {"100": true, "101": true},
			"100": {},
			"101": {"102": true},
			"102": {},
		},
	},
	{ // 12
		Op:  tUnlink,
		SID: "10A",
		Map: map[string]map[string]bool{
			"0AA": {"00A": true, "00B": true},
			"00A": {"000": true},
			"000": {},
			"00B": {},

			"1AA": {},
		},
	},
}

func TestLinking(t *testing.T) {
	for idx, test := range linkingTests {
		switch test.Op {
		case tConnect:
			Get(test.SID, true)
		case tLink:
			s := Get(test.UpSID, true)
			if s == nil {
				t.Errorf("%d. Get(%s) returned nil (should be linked)", idx, test.UpSID)
			}
			err := Link(test.UpSID, test.SID, test.SID+"serv", test.Hops, test.SID+"desc")
			if err != nil {
				t.Errorf("%d. Link(%s to %s) returned %q", idx, test.UpSID, test.SID, err)
			}
		case tUnlink:
			Unlink(test.SID)
		}
		if got, want := len(downstream), len(test.Map); got != want {
			t.Errorf("%d. ds len() = %d, want %d", idx, got, want)
		}
		for up, down := range test.Map {
			if got, want := len(downstream[up]), len(down); got != want {
				t.Errorf("%d. ds[%s] len() = %d, want %d", idx, up, got, want)
			}
			for leaf := range down {
				if !downstream[up][leaf] {
					t.Errorf("%d. ds[%s][%s] = false, want true", idx, up, leaf)
				}
			}
		}
	}
}

// MAP
//  |- 0AA
//  |   |- 00A
//  |   |   `- 000
//  |   `- 00B
//  `- 1AA

var linkIterTest = []string{"0AA", "1AA"}
var linkIterForTests = []struct {
	For    []string
	Except string
	SIDs   []string
}{
	{
		For:  []string{"0AA"},
		SIDs: []string{"0AA"},
	},
	{
		For:  []string{"1AA"},
		SIDs: []string{"1AA"},
	},
	{
		For:  []string{"00A"},
		SIDs: []string{"0AA"},
	},
	{
		For:  []string{"000"},
		SIDs: []string{"0AA"},
	},
	{
		For:  []string{"00B"},
		SIDs: []string{"0AA"},
	},
	{
		For:  []string{"000", "00A", "1AA", "00B"},
		SIDs: []string{"0AA", "1AA"},
	},
	{
		For:    []string{"000AAAAA", "00ABLAH", "1AASOOO", "00BTEST"},
		Except: "0AA",
		SIDs:   []string{"1AA"},
	},
	{
		For:    []string{"000AAAAA", "00ABLAH", "1AASOOO", "00BTEST"},
		Except: "00A",
		SIDs:   []string{"1AA"},
	},
	{
		For:    []string{"000AAAAA", "00ABLAH", "1AASOOO", "00BTEST"},
		Except: "1AA",
		SIDs:   []string{"0AA"},
	},
}

func TestIter(t *testing.T) {
	var ch <-chan string

	ch = Iter()
	for idx, test := range linkIterTest {
		if got, want := <-ch, test; got != want {
			t.Errorf("Iter()[%d] = %s, want %s", idx, got, want)
		}
	}
	for extra := range ch {
		t.Errorf("Iter(): %q (extra value)", extra)
	}

	for idx, test := range linkIterForTests {
		ch = IterFor(test.For, test.Except)
		for i, want := range test.SIDs {
			if got := <-ch; got != want {
				t.Errorf("%d. IterFor()[%d] = %s, want %s", idx, i, got, want)
			}
		}
		if _, open := <-ch; open {
			t.Errorf("%d. IterFor(): channel still open", idx)
		}
	}
}
