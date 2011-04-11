package conn

import (
	"net"
)

type Conn struct {
	net.Conn
}

func NewConn(nc net.Conn) *Conn {
	c := new(Conn)
	c.Conn = c
	return c
}
