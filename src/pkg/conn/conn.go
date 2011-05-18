package conn

import (
	"os"
	"net"
	"encoding/line"
	"kevlar/ircd/parser"
)

type Conn struct {
	net.Conn
	active      bool
	subscribers map[chan *parser.Message]bool
	Error       os.Error
}

func NewConn(nc net.Conn) *Conn {
	c := new(Conn)
	c.Conn = nc
	c.active = true
	c.subscribers = make(map[chan *parser.Message]bool)
	go c.readthread()
	return c
}

func (c *Conn) readthread() {
	// Always close the connection
	defer c.Close()

	// Read lines by \r\n or \n
	linereader := line.NewReader(c, 512)
	for c.active {
		line, _, err := linereader.ReadLine()
		if err != nil {
			c.active = false
			c.Error = err
			return
		}
		message := parser.ParseMessage(line)
		if message != nil {
			for subscriber := range c.subscribers {
				subscriber <- message
			}
		}
	}
}

func (c *Conn) WriteMessage(message *parser.Message) {
	bytes := message.Bytes()
	n, err := c.Write(bytes)
	if err != nil || n != len(bytes) {
		c.Error = err
		c.active = false
		c.Close()
	}
}

func (c *Conn) Active() bool {
	return c.active
}

func (c *Conn) Subscribe(chn chan *parser.Message) {
	c.subscribers[chn] = true
}

func (c *Conn) Unsubscribe(chn chan *parser.Message) {
	c.subscribers[chn] = false, false
}
