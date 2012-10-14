package mode

import (
	"os"
	"path/filepath"
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
	char   int
	name   string
	typ    modeType
	prefix string
}

func (ms *ModeSpec) Char() int      { return ms.char }
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
	}
	panic("unreachable")
}

func newModeSpec(ch int, typ modeType, name string) *ModeSpec {
	var prefix string
	if typ == StatusMode {
		if idx := strings.IndexRune(statusMode, int(ch)); idx >= 0 {
			prefix = statusPrefix
		}
	}
	return &ModeSpec{
		char:   ch,
		name:   name,
		typ:    typ,
		prefix: prefix,
	}
}

const (
	statusPrefix = "@%+"
	statusMode   = "ohv"
)

type ModeMap map[int]*ModeSpec

var (
	UserModes = ModeMap{
		'D': newModeSpec('D', UserMode, "deaf"),
		'S': newModeSpec('S', UserMode, "network service"),
		'a': newModeSpec('a', UserMode, "server administrator"),
		'i': newModeSpec('i', UserMode, "invisible"),
		'o': newModeSpec('o', UserMode, "IRC operator"),
		'w': newModeSpec('w', UserMode, "wallops recipient"),
		'Z': newModeSpec('Z', UserMode, "SSL user"),
		'r': newModeSpec('r', UserMode, "registered with services"),
	}
	ChannelModes = ModeMap{
		'o': newModeSpec('o', StatusMode, "channel operator"),
		'h': newModeSpec('h', StatusMode, "channel half-operator"),
		'v': newModeSpec('v', StatusMode, "channel voice"),
		'b': newModeSpec('b', ListMode, "banned"),
		'e': newModeSpec('e', ListMode, "exempt from +b"),
		'I': newModeSpec('I', ListMode, "exempt from +i"),
		'k': newModeSpec('k', KeyMode, "key required to join"),
		'l': newModeSpec('l', LimitMode, "user count limit"),
		'm': newModeSpec('m', FlagMode, "moderated"),
		'n': newModeSpec('n', FlagMode, "no external messages"),
		'p': newModeSpec('p', FlagMode, "private"),               // -NAMES -KNOCK
		'r': newModeSpec('r', FlagMode, "registered users only"), // only +r users may join
		's': newModeSpec('s', FlagMode, "secret"),                // -NAMES -WHOIS -LIST -KNOCK
		't': newModeSpec('t', FlagMode, "operators control topic"),
	}
)

type modeOp int

const (
	SetMode modeOp = iota
	UnsetMode
	QueryMode
	ActiveMode
)

// A Mode represents either an active mode or a mode change
type Mode struct {
	// Spec represents the actual mode and can be used in a switch statement.
	Spec *ModeSpec

	// Operation describes whether the mode is active, being set, being unset,
	// or being queried.
	Op modeOp

	// Argument stores the textual representations of the arguments. If this is
	// a status mode or a list mode, it will have at least one argument. If this is a
	// flag mode, it will have no arguments.  If it is a limit mode being unset,
	// it will have no arguments.
	Args []string
}

func ParseModeChange(args []string, modes ModeMap) (changes []Mode, errors []os.Error) {
	op := QueryMode

	if len(args) == 0 {
		return nil, nil
	}

	pmstring := args[0]
	argv := args[1:]

	for _, ch := range pmstring {
		switch ch {
		case '+':
			op = SetMode
			continue
		case '-':
			op = UnsetMode
			continue
		}

		ms, ok := modes[ch]
		if !ok {
			errors = append(errors, &UnknownModeError{ch})
			continue
		}

		modech := Mode{
			Spec: ms,
			Op:   op,
		}

		nargs := 0
		set, unset := ms.Args()

		switch op {
		case QueryMode:
		case SetMode:
			nargs = set
		case UnsetMode:
			nargs = unset
		}

		if nargs > len(argv) {
			errors = append(errors, &MissingArgumentError{ch})
			continue
		}

		if nargs > 0 {
			modech.Args = argv[:nargs]
			argv = argv[nargs:]
		}

		changes = append(changes, modech)
	}
	return
}

type ActiveModes map[int]Mode

func (am ActiveModes) Apply(modes []Mode) (applied []Mode, errors []os.Error) {
	// TODO(kevlar): Itarate over UserModes, ChannelModes to ensure stable order
	for _, m := range modes {
		char, op, typ := m.Spec.char, m.Op, m.Spec.typ
		curr, isset := am[char]

		switch typ {
		case UserMode, FlagMode:
			if op == SetMode {
				if isset {
					continue
				}
				am[char] = m
				applied = append(applied, m)
			} else if op == UnsetMode {
				if !isset {
					continue
				}
				am[char] = Mode{}, false
				applied = append(applied, m)
			}
		case KeyMode, LimitMode:
			if op == SetMode {
				if isset {
					continue
				}
				am[char] = m
				applied = append(applied, m)
			} else if op == UnsetMode {
				if !isset {
					continue
				}
				if curr.Args[0] != m.Args[0] {
					errors = append(errors, &UnsetMatchError{char})
					continue
				}
				am[char] = Mode{}, false
				applied = append(applied, m)
			}
		case StatusMode, ListMode:
			args := make(map[string]bool)
			for _, arg := range curr.Args {
				args[arg] = true
			}
			for _, arg := range m.Args {
				if op == SetMode {
					args[arg] = true
				} else if op == UnsetMode {
					args[arg] = false, false
				}
			}
			applied = append(applied, m)

			if len(args) == 0 {
				am[char] = Mode{}, false
				continue
			}

			curr.Spec = m.Spec
			curr.Args = []string{}
			for arg := range args {
				curr.Args = append(curr.Args, arg)
			}
			am[char] = curr
		}
	}
	return
}

func (am ActiveModes) String() string {
	modes := make([]Mode, 0, len(am))
	for _, mode := range am {
		modes = append(modes, mode)
	}
	return ModeString(modes)
}

func (am ActiveModes) Get(ch int) (isset bool, args []string) {
	var m Mode
	if m, isset = am[ch]; isset {
		args = m.Args
	}
	return
}

func (am ActiveModes) Match(ch int, ref string) bool {
	var m Mode
	var isset bool
	if m, isset = am[ch]; !isset {
		return false
	}
	for _, arg := range m.Args {
		if match, _ := filepath.Match(arg, ref); match {
			return true
		}
	}
	return false
}

func (am ActiveModes) Contains(ch int, ref string) bool {
	var m Mode
	var isset bool
	if m, isset = am[ch]; !isset {
		return false
	}
	for _, arg := range m.Args {
		if arg == ref {
			return true
		}
	}
	return false
}

func ModeString(modes []Mode) string {
	op := QueryMode
	modebytes := []byte{}
	args := []string{""}

	putch := func(ch int) {
		modebytes = append(modebytes, byte(ch))
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
			putch(mode.Spec.char)
			continue
		}
		for _, arg := range mode.Args {
			putch(mode.Spec.char)
			args = append(args, arg)
		}
	}
	args[0] = string(modebytes)
	return strings.Join(args, " ")
}

type UnknownModeError struct{ Char int }
type MissingArgumentError struct{ Char int }
type UnsetMatchError struct{ Char int }

func (e *UnknownModeError) String() string     { return string(e.Char) + " is an unknown mode to me" }
func (e *MissingArgumentError) String() string { return string(e.Char) + " requires argument" }
func (e *UnsetMatchError) String() string      { return "mismatch on unset " + string(e.Char) }
