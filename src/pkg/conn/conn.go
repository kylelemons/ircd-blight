package conn

import (
	"os"
	"log"
	"net"
	"encoding/line"
	"kevlar/ircd/parser"
	"kevlar/ircd/user"
)

type Conn struct {
	net.Conn
	active      bool
	subscribers map[chan<- *parser.Message]bool
	onclose     map[chan<- string]bool
	Error       os.Error
	id          string
}

func NewConn(nc net.Conn) *Conn {
	c := &Conn{
		Conn:        nc,
		active:      true,
		subscribers: make(map[chan<- *parser.Message]bool),
		onclose:     make(map[chan<- string]bool),
		id:          user.NextUserID(),
	}
	log.Printf("[%s] ** Connected", c.id)
	go c.readthread()
	return c
}

func (c *Conn) Close() os.Error {
	for ch := range c.onclose {
		ch <- c.id
	}
	return c.Conn.Close()
}

func (c *Conn) ID() string {
	return c.id
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
		message.SenderID = c.id
		if message != nil {
			for subscriber := range c.subscribers {
				subscriber <- message
			}
		}
	}
}

func (c *Conn) WriteMessage(message *parser.Message) {
	bytes := message.Bytes()
	bytes = append(bytes, '\r', '\n')
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

func (c *Conn) Subscribe(chn chan<- *parser.Message) {
	c.subscribers[chn] = true
}

func (c *Conn) SubscribeClose(chn chan<- string) {
	c.onclose[chn] = true
}

func (c *Conn) Unsubscribe(chn chan *parser.Message) {
	c.subscribers[chn] = false, false
}

func (c *Conn) UnsubscribeClose(chn chan<- string) {
	c.onclose[chn] = false, false
}
