package mode

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

type modeType int

const (
	// User modes can only be applied to users
	UserMode modeType = iota

	// Status modes apply to specific users on a channel
	StatusMode

	// Mask modes are lists (typically of nick!user@host masks)
	ListMode

	// Key modes require the same argument to remove as to set
	KeyMode

	// Limit modes require no argument to remove them
	LimitMode

	// Flag modes require no arguments to set or remove
	FlagMode
)

type ModeSpec struct {
	char   rune
	name   string
	typ    modeType
	prefix string
}

func (ms *ModeSpec) Char() rune     { return ms.char }
func (ms *ModeSpec) Name() string   { return ms.name }
func (ms *ModeSpec) Type() modeType { return ms.typ }
func (ms *ModeSpec) Prefix() string { return ms.prefix }

// Args returns the number of arguments for this mode change
// when it is being set and when it is being unset.
func (ms *ModeSpec) Args() (set, unset int) {
	switch ms.typ {
	case UserMode, FlagMode:
		return 0, 0
	case StatusMode, ListMode, KeyMode:
		return 1, 1
	case LimitMode:
		return 1, 0
	default:
		panic("unknown mode type")
	}
}

func newModeSpec(r rune, typ modeType, name string) ModeSpec {
	const (
		StatusPrefixes = "@%+"
		StatusModes    = "ohv"
	)
	var prefix string
	if typ == StatusMode {
		if idx := strings.IndexRune(StatusModes, r); idx >= 0 {
			prefix = StatusPrefixes[idx : idx+1]
		}
	}

	return ModeSpec{
		char:   r,
		name:   name,
		typ:    typ,
		prefix: prefix,
	}
}

type ModeMap struct {
	Modes []ModeSpec // The set of modes that are valid in this mode set
	Index [256]byte  // The index into the Modes map for a given character
}

func MakeModeMap(specs ...ModeSpec) *ModeMap {
	mm := new(ModeMap)
	mm.Modes = make([]ModeSpec, len(specs)+1)
	copy(mm.Modes[1:], specs)
	if len(mm.Modes) > len(mm.Index) {
		log.Fatalf("Too many modes for index array")
	}
	for i, s := range mm.Modes {
		if i == 0 {
			continue
		}
		if int(s.char) >= len(mm.Index) {
			log.Fatalf("Mode character %q (%d) out of bounds", s.char, s.char)
		}
		mm.Index[int(s.char)] = byte(i)
	}
	return mm
}

func (mm *ModeMap) Mode(r rune) (spec *ModeSpec, index byte, ok bool) {
	if int(r) >= len(mm.Index) {
		return nil, 0, false
	}
	i := mm.Index[int(r)]
	if i == 0 {
		return nil, 0, false
	}
	return &mm.Modes[int(i)], i, true
}

func (mm *ModeMap) For(m Mode) *ModeSpec {
	return &mm.Modes[m.Index]
}

var (
	UserModes = MakeModeMap(
		newModeSpec('D', UserMode, "deaf"),
		newModeSpec('S', UserMode, "network service"),
		newModeSpec('a', UserMode, "server administrator"),
		newModeSpec('i', UserMode, "invisible"),
		newModeSpec('o', UserMode, "IRC operator"),
		newModeSpec('w', UserMode, "wallops recipient"),
		newModeSpec('Z', UserMode, "SSL user"),
		newModeSpec('r', UserMode, "registered with services"),
	)
	ChannelModes = MakeModeMap(
		newModeSpec('o', StatusMode, "channel operator"),
		newModeSpec('h', StatusMode, "channel half-operator"),
		newModeSpec('v', StatusMode, "channel voice"),
		newModeSpec('b', ListMode, "banned"),
		newModeSpec('e', ListMode, "exempt from +b"),
		newModeSpec('I', ListMode, "exempt from +i"),
		newModeSpec('k', KeyMode, "key required to join"),
		newModeSpec('l', LimitMode, "user count limit"),
		newModeSpec('m', FlagMode, "moderated"),
		newModeSpec('n', FlagMode, "no external messages"),
		newModeSpec('p', FlagMode, "private"),               // -NAMES -KNOCK
		newModeSpec('r', FlagMode, "registered users only"), // only +r users may join
		newModeSpec('s', FlagMode, "secret"),                // -NAMES -WHOIS -LIST -KNOCK
		newModeSpec('t', FlagMode, "operators control topic"),
	)
)

type modeOp int

const (
	SetMode   modeOp = iota // Set a mode, e.g. "v"
	UnsetMode               // Remove a mode, e.g. "-v"
	QueryMode               // Check a mode, e.g. "+v"
)

// A Mode represents either an active mode or a mode change
type Mode struct {
	// Index represents the index of the mode in its accompanying mode spec
	Index byte

	// Operation describes whether the mode is active, being set, being unset,
	// or being queried.
	Op modeOp

	// Argument stores the textual representations of the arguments. If this is
	// a status mode or a list mode, it will have at least one argument. If this is a
	// flag mode, it will have no arguments.  If it is a limit mode being unset,
	// it will have no arguments.
	Args []string
}

func (mm *ModeMap) ParseModeChange(args []string) (changes []Mode, errors []error) {
	op := QueryMode

	if len(args) == 0 {
		return nil, nil
	}

	pmstring := args[0]
	argv := args[1:]

nextMode:
	for _, r := range pmstring {
		switch r {
		case '+':
			op = SetMode
			continue
		case '-':
			op = UnsetMode
			continue
		}

		spec, index, ok := mm.Mode(r)
		if !ok {
			errors = append(errors, &UnknownModeError{r})
			continue
		}

		modech := Mode{
			Index: index,
			Op:    op,
		}

		nargs := 0
		set, unset := spec.Args()

		switch op {
		case QueryMode:
		case SetMode:
			nargs = set
		case UnsetMode:
			nargs = unset
		}

		if nargs > len(argv) {
			errors = append(errors, &MissingArgumentError{r})
			continue
		}

		if nargs > 0 {
			modech.Args = argv[:nargs]
			argv = argv[nargs:]
		}

		for _, arg := range modech.Args {
			if len(arg) == 0 {
				errors = append(errors, &MissingArgumentError{r})
				continue nextMode
			}
			for _, ar := range arg {
				// TODO: Use something unicode safe like:
				//   unicode.Letter, unicode.Number, unicode.Punct, unicode.Sm
				if ar < '!' || ar > '~' {
					errors = append(errors, &InvalidArgumentError{r, arg, string(ar)})
					continue nextMode
				}
			}
		}

		changes = append(changes, modech)
	}
	if len(argv) > 0 {
		errors = append(errors, &TooManyArgumentsError{argv})
	}
	if len(changes) == 0 {
		errors = append(errors, ErrNoModeChange)
	}
	return changes, errors
}

type ActiveModes struct {
	Type  *ModeMap
	Modes ByModeOrder // must be kept sorted
}

func NewActiveModes(set *ModeMap) *ActiveModes {
	return &ActiveModes{
		Type: set,
	}
}

func (am *ActiveModes) Apply(modes []Mode) (applied []Mode, errors []error) {
	for _, m := range modes {
		if err := am.apply(m); err != nil {
			errors = append(errors, err)
		} else {
			applied = append(applied, m)
		}
	}
	return modes, errors
}

func (am *ActiveModes) Lookup(r rune) (args []string, found bool) {
	index := am.Type.Index[int(r)]
	i := sort.Search(len(am.Modes), func(i int) bool {
		return am.Modes[i].Index >= index
	})
	if i < len(am.Modes) && am.Modes[i].Index == index {
		return am.Modes[i].Args, true
	}
	return nil, false
}

func (am *ActiveModes) apply(mode Mode) error {
	i := sort.Search(len(am.Modes), func(i int) bool {
		return am.Modes[i].Index >= mode.Index
	})
	found := i < len(am.Modes) && am.Modes[i].Index == mode.Index

	switch am.Type.For(mode).typ {
	case FlagMode, UserMode, KeyMode, LimitMode:
		switch {
		case found && mode.Op == UnsetMode:
			// Mode is being unset:
			am.deleteAt(i)
		case !found && mode.Op == SetMode:
			// Mode is being set, insert it:
			am.insertAt(i, mode)
		}
	case StatusMode, ListMode:
		switch {
		case !found && mode.Op == SetMode:
			// First value in a list mode
			am.insertAt(i, mode)
		case found && mode.Op == SetMode:
			want := mode.Args[0]
			for _, v := range am.Modes[i].Args {
				if v == want {
					// Value is already present
					return nil
				}
			}
			am.Modes[i].Args = append(am.Modes[i].Args, want)
		case found && mode.Op == UnsetMode:
			want := mode.Args[0]
			for j, v := range am.Modes[i].Args {
				if v == want {
					copy(am.Modes[i].Args[j:], am.Modes[i].Args[j+1:])
					am.Modes[i].Args = am.Modes[i].Args[:len(am.Modes[i].Args)-1]
				}
			}
			if len(am.Modes[i].Args) == 0 {
				am.deleteAt(i)
			}
		}
	}
	return nil
}

func (am *ActiveModes) insertAt(i int, m Mode) {
	am.Modes = append(am.Modes, Mode{})
	copy(am.Modes[i+1:], am.Modes[i:])
	am.Modes[i] = m
}

func (am *ActiveModes) deleteAt(i int) {
	copy(am.Modes[i:], am.Modes[i+1:])
	am.Modes = am.Modes[:len(am.Modes)-1]
}

func (am ActiveModes) String() string {
	return am.Type.ModeString(am.Modes)
}

func (am ActiveModes) Match(r rune, ref string) bool {
	args, _ := am.Lookup(r)
	for _, arg := range args {
		if match, _ := filepath.Match(arg, ref); match {
			return true
		}
	}
	return false
}

func (am ActiveModes) Contains(r rune, ref string) bool {
	args, _ := am.Lookup(r)
	for _, arg := range args {
		if arg == ref {
			return true
		}
	}
	return false
}

func (mm *ModeMap) ModeString(modes []Mode) string {
	op := QueryMode
	modebytes := []byte{}
	args := []string{""}

	putch := func(ch byte) {
		modebytes = append(modebytes, ch)
	}
	putm := func(i int) {
		modebytes = append(modebytes, byte(mm.Modes[i].char))
	}
	puts := func(s string) {
		args = append(args, s)
	}

	for _, mode := range modes {
		if mode.Op != op {
			switch mode.Op {
			case SetMode:
				putch('+')
			case UnsetMode:
				putch('-')
			case QueryMode:
				// can't do a query after the first + or -
				continue
			}
		}
		op = mode.Op
		if len(mode.Args) == 0 {
			putm(int(mode.Index))
			continue
		}
		for _, arg := range mode.Args {
			putm(int(mode.Index))
			puts(arg)
		}
	}
	args[0] = string(modebytes)
	return strings.Join(args, " ")
}

type ByModeOrder []Mode

func (v ByModeOrder) Len() int           { return len(v) }
func (v ByModeOrder) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByModeOrder) Less(i, j int) bool { return CompareModes(v[i], v[j]) < 0 }

func CompareModes(x, y Mode) int {
	if a, b := x.Index, y.Index; a != b {
		return int(b - a)
	}
	if a, b := len(x.Args), len(y.Args); a != b {
		return b - a // shouldn't happen, but makes the loop below safe
	}
	for arg := range x.Args {
		a, b := x.Args[arg], y.Args[arg]
		if d := strings.Compare(a, b); d != 0 {
			return d
		}
	}
	return 0
}

type UnknownModeError struct{ Rune rune }
type MissingArgumentError struct{ Rune rune }
type UnsetMatchError struct{ Rune rune }
type TooManyArgumentsError struct{ Extra []string }

func (e *UnknownModeError) Error() string     { return string(e.Rune) + " is an unknown mode to me" }
func (e *MissingArgumentError) Error() string { return string(e.Rune) + " requires argument" }
func (e *UnsetMatchError) Error() string      { return "mismatch on unset " + string(e.Rune) }

func (e *TooManyArgumentsError) Error() string {
	return "extra arguments: " + strings.Join(e.Extra, " ")
}

type InvalidArgumentError struct {
	Rune     rune
	Argument string
	What     string
}

func (e *InvalidArgumentError) Error() string {
	return fmt.Sprintf("argument %q invalid for mode %c: %s", e.Argument, e.Rune, e.What)
}

var ErrNoModeChange = errors.New("no mode changes requested")
