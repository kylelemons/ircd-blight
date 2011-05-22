package parser

import (
	"strings"
	"log"
)

type Numeric struct {
	num  string
	args []string
}

func NewNumeric(num string, args ...string) *Numeric {
	text, ok1 := NumericText[num]
	name, ok2 := NumericName[num]
	if !ok1 || !ok2 {
		return nil
	}

	// Strip out the format (before the : in NumericText rfc*.go)
	pieces := strings.Split(text, " :", 2)
	numargs := append([]string{"*"}) // reserve one for the nick

	argcnt := 0
	if len(pieces) > 1 {
		argcnt = strings.Count(pieces[0], "<")
	}

	if got, want := len(args), argcnt; got != want {
		log.Printf("Warning: %d arguments to %s, want %d", got, name, want)
	}
	numargs = append(numargs, args...)
	numargs = append(numargs, pieces[len(pieces)-1])

	return &Numeric{
		num:  num,
		args: numargs,
	}
}

func (n *Numeric) String() string {
	name, ok := NumericName[n.num]
	if !ok {
		return n.num
	}
	return name
}

func (n *Numeric) Message(destIDs ...string) *Message {
	return &Message{
		Command: n.num,
		Args:    n.args,
		DestIDs: destIDs,
	}
}
