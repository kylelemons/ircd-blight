package server

import (
	"log"
	"kevlar/ircd/conn"
	"kevlar/ircd/core"
	"kevlar/ircd/parser"
	"kevlar/ircd/user"
)

func Start() {
	listener := conn.NewListener()
	defer listener.Close()
	for _, ports := range Config.Ports {
		portlist, err := ports.GetPortList()
		if err != nil {
			log.Printf("Warning: %s\n", err)
		}
		for _, port := range portlist {
			listener.AddPort(port)
		}
	}

	// TODO(kevlar): Configurable sendq and recvq

	incoming := make(chan *parser.Message, 100)
	outgoing := make(chan *parser.Message, 100)
	closing := make(chan string)
	connIDs := make(map[string]*conn.Conn)

	for {
		select {
		// TODO(kevlar): Event dispatch channel?
		case closed := <-closing:
			log.Printf("[%s] ** Connection closed", closed)
			user.Delete(closed)
			connIDs[closed] = nil, false
		case msg := <-incoming:
			log.Printf("[%s] >> %s\n", msg.SenderID, msg)
			core.DispatchMessage(msg, outgoing)
		case msg := <-outgoing:
			sentcount := 0
			setnick := len(msg.Args) > 0 && msg.Args[0] == "*"
			setprefix := msg.Prefix == "*"
			closeafter := msg.Command == parser.CMD_ERROR
			if len(msg.Prefix) == 0 && !closeafter {
				msg.Prefix = Config.Name
			}
			for _, id := range msg.DestIDs {
				conn, ok := connIDs[id]
				if !ok {
					log.Printf("Warning: Nonexistent ID %s in send", id)
					continue
				}
				if setnick || setprefix {
					nick, _, _, _, _ := user.GetInfo(id)
					if setnick {
						msg.Args[0] = nick
					}
					if setprefix {
						msg.Prefix = nick
					}
				}
				log.Printf("[%s] << %s\n", id, msg)
				conn.WriteMessage(msg)
				sentcount++
				if closeafter {
					log.Printf("[%s] ** Connection terminated", id)
					user.Delete(id)
					connIDs[id] = nil, false
					conn.UnsubscribeClose(closing)
					conn.Close()
				}
			}
			if sentcount == 0 {
				log.Printf("Dropped outgoing message: %s", msg)
			}
		case conn := <-listener.Incoming:
			id := conn.ID()
			connIDs[id] = conn
			user.Get(id)
			conn.Subscribe(incoming)
			conn.SubscribeClose(closing)
		}
	}
}

func Run() {
	LoadConfigString(DefaultXML)

	Start()
}
