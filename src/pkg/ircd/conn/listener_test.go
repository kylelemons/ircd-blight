package conn

import (
	"testing"
)

import "kevlar/ircd/util"

func TestNewListener(t *testing.T) {
	cmp := util.Test(t)
	l := NewListener()
	cmp.NE("ports", nil, l.ports)
	cmp.NE("Incoming", nil, l.Incoming)
	l.Close()
}

func TestAddPort(t *testing.T) {
	cmp := util.Test(t)
	l := NewListener()
	l.AddPort(56561)
	cmp.EQ("len", 1, len(l.ports))
	listener,ok := l.ports[56561]
	cmp.NE("56561", nil, ok)
	cmp.NE("listener", nil, listener)
	l.Close()
	cmp.EQ("len", 0, len(l.ports))
}

func TestClosePort(t *testing.T) {
	cmp := util.Test(t)
	l := NewListener()
	l.AddPort(56561)
	cmp.EQ("len", 1, len(l.ports))
	l.ClosePort(56561)
	cmp.EQ("len", 0, len(l.ports))
	l.Close()
	cmp.EQ("len", 0, len(l.ports))
}
