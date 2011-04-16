package conn

import (
	"os"
	"net"
	"encoding/line"
	"kevlar/ircd/parser"
)

type Conn struct {
	net.Conn
	active bool
	messages chan *parser.Message
	Error os.Error
}

func NewConn(nc net.Conn) *Conn {
	c := new(Conn)
	c.Conn = nc
	c.active = true
	c.messages = make(chan *parser.Message)
	go c.readthread()
	return c
}

func (c *Conn) readthread() {
	// Always close the connection
	defer c.Close()

	// Always close the message channel
	defer close(c.messages)

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
			c.messages <- message
		}
	}
}

func (c *Conn) Iter() chan<- *parser.Message {
	return c.messages
}

func (c *Conn) ReadMessage() *parser.Message {
	return <-c.messages
}

func (c *Conn) WriteMessage(message *parser.Message) {
	bytes := message.Bytes()
	n,err := c.Write(bytes)
	if err != nil || n != len(bytes) {
		c.Error = err
		c.active = false
		c.Close()
	}
}

func (c *Conn) Active() bool {
	return c.active
}
