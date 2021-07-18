package parser

import (
	"bytes"
	"strings"
)

type Message struct {
	Prefix  string
	Command string
	Args    []string

	SenderID string
	DestIDs  []string
}

func NewMessage(pfx, cmd string, args []string) *Message {
	m := new(Message)
	m.Prefix = pfx
	m.Command = cmd
	m.Args = args
	return m
}

// Duplicate a message.  This is a deep copy.
func (m *Message) Dup() *Message {
	n := new(Message)
	n.Prefix, n.Command, n.SenderID = m.Prefix, m.Command, m.SenderID
	n.Args = make([]string, len(m.Args))
	copy(n.Args, m.Args)
	n.DestIDs = make([]string, len(m.DestIDs))
	copy(n.DestIDs, m.DestIDs)
	return n
}

func ParseMessage(line []byte) *Message {
	line = bytes.TrimSpace(line)
	if len(line) <= 0 {
		return nil
	}
	m := new(Message)
	if line[0] == ':' {
		split := bytes.SplitN(line, []byte{' '}, 2)
		if len(split) <= 1 {
			return nil
		}
		m.Prefix = string(split[0][1:])
		line = split[1]
	}
	split := bytes.SplitN(line, []byte{':'}, 2)
	args := bytes.Split(bytes.TrimSpace(split[0]), []byte{' '})
	m.Command = string(bytes.ToUpper(args[0]))
	m.Args = make([]string, 0, len(args))
	for _, arg := range args[1:] {
		m.Args = append(m.Args, string(arg))
	}
	if len(split) > 1 {
		m.Args = append(m.Args, string(split[1]))
	}
	return m
}

func (m Message) Bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	if len(m.Prefix) > 0 {
		buf.WriteByte(':')
		buf.WriteString(m.Prefix)
		buf.WriteByte(' ')
	}
	buf.WriteString(m.Command)
	for i, arg := range m.Args {
		buf.WriteByte(' ')
		if i == len(m.Args)-1 {
			if strings.IndexAny(arg, " :") >= 0 {
				buf.WriteByte(':')
			}
		}
		buf.WriteString(m.Args[i])
	}
	return buf.Bytes()
}

func (m Message) String() string {
	return string(m.Bytes())
}
