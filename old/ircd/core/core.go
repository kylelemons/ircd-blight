package core

import (
	"fmt"

	"github.com/kylelemons/ircd-blight/old/ircd/conn"
	"github.com/kylelemons/ircd-blight/old/ircd/parser"
)

type Core struct {
	ports    map[int]bool
	listener *conn.Listener
	messages chan *parser.Message
	params   map[string]map[string]string
}

func NewCore() *Core {
	core := &Core{
		ports:    make(map[int]bool),
		messages: make(chan *parser.Message),
		params:   make(map[string]map[string]string),
	}
	return core
}

func (c *Core) Set(module, key, val string) {
	if _, ok := c.params[module]; !ok {
		c.params[module] = make(map[string]string)
	}
	c.params[module][key] = val
}

func (c *Core) Get(module, key, defval string) string {
	if mod, ok := c.params[module]; ok {
		if val, ok := mod[key]; ok {
			return val
		}
	}
	return defval
}

func (c *Core) Unset(module, key string) {
	if mod, ok := c.params[module]; ok {
		if _, ok := mod[key]; ok {
			delete(mod, key)
		}
		if len(mod) == 0 {
			delete(c.params, module)
		}
	}
}

func (c *Core) Start() {
	//TODO(kevlar)

	// Listen on each port
	c.listener = conn.NewListener()
	for _, port := range []int{6666, 6667} {
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
		delete(c.ports, port)
	}
}

func (c *Core) run() {
	for message := range c.messages {
		fmt.Printf("Message: %s\n", message)
	}
}
