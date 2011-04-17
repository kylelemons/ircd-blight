package ircd

import (
	"fmt"
	"kevlar/ircd/conn"
	"kevlar/ircd/parser"
	ds "kevlar/ircd/datastore"
)

type Core struct {
	Data     *ds.DataStore
	ports    map[int]bool
	listener *conn.Listener
	messages chan *parser.Message
}

func NewCore() *Core {
	core := &Core{
		Data:     ds.NewDataStore(),
		ports:    make(map[int]bool),
		messages: make(chan *parser.Message),
	}
	return core
}

func (c *Core) Set(module, key string, val ds.Value) {
	r := ds.NewReturn()
	set := ds.Set{
		Module: module,
		Key:    key,
		Value:  val,
		Return: r,
	}
	c.Data.Control <- set
	<-r
}

func (c *Core) Get(module, key string, def ds.Value) ds.Value {
	r := ds.NewReturn()
	get := &ds.Get{
		Module: module,
		Key:    key,
		Return: r,
		Value:  def,
	}
	c.Data.Control <- get
	<-r
	return get.Value
}

func (c *Core) Start() {
	//TODO(kevlar)

	// Listen on each port
	c.listener = conn.NewListener()
	for _, port := range c.Get("Server", "ports", []int{6666, 6667}).([]int) {
		c.listener.AddPort(port)
		c.ports[port] = true
	}
	go func() {
		for conn := range c.listener.Incoming {
			conn.Subscribe(c.messages)
		}
	}()

	go c.run()
}

func (c *Core) Stop() {
	// Listen on each port
	for port := range c.ports {
		c.listener.ClosePort(port)
		c.ports[port] = false, false
	}
}

func (c *Core) run() {
	for message := range c.messages {
		fmt.Printf("Message: %s", message)
	}
}
