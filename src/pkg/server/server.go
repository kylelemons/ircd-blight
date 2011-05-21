package server

import (
	"fmt"
	"kevlar/ircd/conn"
	"kevlar/ircd/core"
	"kevlar/ircd/parser"
	"kevlar/ircd/user"
	"os"
)

func Start() {
	listener := conn.NewListener()
	defer listener.Close()
	for _, ports := range Config.Ports {
		portlist, err := ports.GetPortList()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %s\n", err)
		}
		for _, port := range portlist {
			listener.AddPort(port)
		}
	}

	messages := make(chan *parser.Message)

	for {
		select {
		// TODO(kevlar): Event dispatch channel?
		case msg := <-messages:
			core.DispatchMessage(msg)
		case conn := <-listener.Incoming:
			user.Get(conn.ID())
			conn.Subscribe(messages)
		}
	}
}

func Run() {
	LoadConfigString(DefaultXML)

	Start()
}
