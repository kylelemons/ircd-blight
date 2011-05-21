package server

import (
	"fmt"
	"kevlar/ircd/conn"
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

	for {
		select {
		case conn := <-listener.Incoming:
			go handleRegistration(conn)
		}
	}
}

func Run() {
	LoadConfigString(DefaultXML)

	Start()
}
