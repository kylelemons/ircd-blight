package conn

import (
	"fmt"
	"net"
	"sync"
)

type Listener struct {
	ports    map[int]net.Listener
	Incoming chan *Conn
	wg       sync.WaitGroup
}

func NewListener() *Listener {
	l := new(Listener)
	l.ports = make(map[int]net.Listener)
	l.Incoming = make(chan *Conn)
	return l
}

// AddPort starts a new goroutine listening on the given port number.
// If the port number is already being listened to, nothing happens.
func (l *Listener) AddPort(portno int) {
	if _, ok := l.ports[portno]; ok {
		return
	}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", portno))
	if err != nil {
		fmt.Sprintf("Error[%d]: %s\n", portno, err)
		return
	}

	l.ports[portno] = listener
	l.wg.Add(1)
	go func() {
		defer listener.Close()
		defer l.wg.Done()
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Sprintf("Error[%d]: %s\n", portno, err)
				break
			}
			if listener.Addr() == nil {
				conn.Close()
				break
			}
			go func(c net.Conn) {
				l.Incoming <- NewConn(c)
			}(conn)
		}
		l.ports[portno] = nil, false
	}()
}

// ClosePort stops listening on the given port.  If this listener
// is not listening on the port, nothing happens.
func (l *Listener) ClosePort(portno int) {
	listener, ok := l.ports[portno]
	if !ok {
		return
	}
	listener.Close()
	c, _ := net.Dial("tcp", fmt.Sprintf(":%d", portno))
	if c != nil {
		c.Close()
	}
}

// Close signals all of the listening ports to stop listening.
func (l *Listener) Close() {
	for port, listener := range l.ports {
		listener.Close()
		c, _ := net.Dial("tcp", fmt.Sprintf(":%d", port))
		if c != nil {
			c.Close()
		}
	}
	l.wg.Wait()
}
