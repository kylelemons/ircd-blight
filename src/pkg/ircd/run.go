package ircd

import (
	"net"
	"os"
	"fmt"
	"bufio"
)

import (
	"kevlar/ircd/parser"
)

func errchk(err os.Error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func Run() {
	proxy, err := net.Listen("tcp", ":6667")
	errchk(err)

	client, err := proxy.Accept()
	errchk(err)

	server, err := net.Dial("tcp", "", "irc.freenode.net:6667")
	errchk(err)

	stopped := make(chan bool)

	go func() {
		in := bufio.NewReader(client)
		for err == nil {
			line, err := in.ReadBytes('\n')
			if err == os.EOF {
				break
			}
			errchk(err)
			server.Write(line)
			m := parser.ParseMessage(line)
			fmt.Printf("<< %s\n", m)
		}

		stopped <- true
	}()

	go func() {
		in := bufio.NewReader(server)
		for err == nil {
			line, err := in.ReadBytes('\n')
			if err == os.EOF {
				break
			}
			errchk(err)
			client.Write(line)
			m := parser.ParseMessage(line)
			fmt.Printf("<< %s\n", m)
		}

		stopped <- true
	}()

	<-stopped
	<-stopped
}
