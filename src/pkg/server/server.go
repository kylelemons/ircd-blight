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
			// Count the number of messages sent
			sentcount := 0

			// For simplicity, * in the prefix or as the first argument
			// is replaced by the user the message is sent to
			setnick := len(msg.Args) > 0 && msg.Args[0] == "*"
			setprefix := msg.Prefix == "*"

			// Close the connection and remove the prefix if we are sending an ERROR
			closeafter := false
			if msg.Command == parser.CMD_ERROR {
				closeafter = true
				msg.Prefix = ""
				// Make sure a prefix is specified (use the server name)
			} else if len(msg.Prefix) == 0 {
				msg.Prefix = Config.Name
			}

			// Examine all arguments for UIDs and replace them
			if isuid(msg.Prefix) {
				nick, user, _, _, ok := user.GetInfo(msg.Prefix)
				if !ok {
					log.Printf("Warning: Nonexistent ID %s as prefix", msg.Prefix)
				} else {
					msg.Prefix = nick + "!" + user + "@host" // TODO(kevlar): hostname
				}
			}
			for i := range msg.Args {
				if isuid(msg.Args[i]) {
					nick, _, _, _, ok := user.GetInfo(msg.Args[i])
					if !ok {
						log.Printf("Warning: Nonexistent ID %s as argument", msg.Args[i])
						continue
					}
					msg.Args[i] = nick
				}
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
				log.Printf("Warning: Dropped outgoing message: %s", msg)
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

func isuid(id string) bool {
	return len(id) == 9 && id[0] >= '0' && id[0] <= '9'
}
