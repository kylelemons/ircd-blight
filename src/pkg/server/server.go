package server

import (
	"kevlar/ircd/conn"
	"kevlar/ircd/core"
	"kevlar/ircd/parser"
	"kevlar/ircd/user"
	"log"
	"sync"
)

func checkConfig() {
	okay := true

	// Check hostname: require at least one .
	if !parser.ValidServerName(Config.Name) {
		log.Printf("Error: invalid server name %q: must match /\\w+(.\\w+)+/", Config.Name)
		okay = false
	}

	// Check prefix; [num][alphanum][alphanum]
	if !parser.ValidServerPrefix(Config.SID) {
		log.Printf("Error: invalid server [refix %q: must match /[0-9][0-9A-Z]{2}/")
		okay = false
	}
	user.UserIDPrefix = Config.SID

	// Check opers
	if len(Config.Operator) == 0 {
		log.Printf("Error: no operators defined: at least one required")
		okay = false
	}

	if !okay {
		log.Fatalf("Invalid configuration; Exiting.")
	}
}

type IRCd struct {
	*core.IRCd
	clientClosing chan string
	newClient     chan *conn.Conn
	serverClosing chan string
	newServer     chan *conn.Conn
	shutdown      chan bool
	running       *sync.WaitGroup
}

func (s *IRCd) manageServers() {
	defer s.running.Done()

	sid2conn := make(map[string]*conn.Conn)
	upstream := ""

	incoming := make(chan *parser.Message)

	var open bool = true
	var msg *parser.Message
	for open {
		// Check if we are connected to our upstream
		// TODO(kevlar): upstream

		select {
		//// Incoming and outgoing messages to and from servers
		// Messages directly from connections
		case msg = <-incoming:
			log.Printf("{%s} >> %s\n", msg.SenderID, msg)
			core.DispatchMessage(msg, s.IRCd)

		// Messages from hooks
		case msg, open = <-s.ToServer:
			// Count the number of messages sent
			//sentcount := 0

			_ = upstream

			log.Printf("{%v} << %s\n", msg.DestIDs, msg)

		//// Connection management
		// Connecting servers
		case conn := <-s.newServer:
			sid := conn.ID()
			sid2conn[sid] = conn
			conn.Subscribe(incoming)
			conn.SubscribeClose(s.serverClosing)
		// Disconnecting servers
		case closeid := <-s.serverClosing:
			log.Printf("{%s} ** Connection closed", closeid)
			user.Delete(closeid)
			sid2conn[closeid] = nil, false
		}
	}
}

func (s *IRCd) manageClients() {
	defer s.running.Done()

	uid2conn := make(map[string]*conn.Conn)
	incoming := make(chan *parser.Message)

	var open bool = true
	var msg *parser.Message
	for open {
		// anything to do here?

		select {
		//// Incoming and outgoing messages to and from clients
		// Messages directly from connections
		case msg = <-incoming:
			log.Printf("[%s] >> %s\n", msg.SenderID, msg)
			core.DispatchMessage(msg, s.IRCd)
		// Messages from hooks
		case msg, open = <-s.ToClient:
			// Count the number of messages sent
			sentcount := 0

			// For simplicity, * in the prefix or as the first argument
			// is replaced by the nick of the user the message is sent to
			setnick := len(msg.Args) > 0 && msg.Args[0] == "*"
			setprefix := msg.Prefix == "*"

			closeafter := false
			if msg.Command == parser.CMD_ERROR {
				// Close the connection and remove the prefix if we are sending an ERROR
				closeafter = true
				msg.Prefix = ""
			} else if len(msg.Prefix) == 0 {
				// Make sure a prefix is specified (use the server name)
				msg.Prefix = Config.Name
			}

			local := make([]string, 0, len(msg.DestIDs))
			remote := make([]string, 0, len(msg.DestIDs))

			for _, id := range msg.DestIDs {
				if sid := id[0:3]; sid != Config.SID {
					remote = append(remote, id)
					continue
				}
				local = append(local, id)
			}

			// Pass the message to the server goroutine
			if len(remote) > 0 {
				msg.DestIDs = local
				msg := msg.Dup()
				msg.DestIDs = remote
				s.ToServer <- msg

				// Short circuit if there are no local recipients
				if len(local) == 0 {
					continue
				}
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
				conn, ok := uid2conn[id]
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
				conn.WriteMessage(msg)
				log.Printf("[%s] << %s\n", id, msg)
				sentcount++
				if closeafter {
					log.Printf("[%s] ** Connection terminated", id)
					user.Delete(id)
					uid2conn[id] = nil, false
					conn.UnsubscribeClose(s.clientClosing)
					conn.Close()
				}
			}
			if sentcount == 0 {
				log.Printf("Warning: Dropped outgoing message: %s", msg)
			}

		//// Connection management
		// Connecting clients
		case conn := <-s.newClient:
			id := conn.ID()
			uid2conn[id] = conn
			user.Get(id)
			conn.Subscribe(incoming)
			conn.SubscribeClose(s.clientClosing)
		// Disconnecting clients
		case closeid := <-s.clientClosing:
			log.Printf("[%s] ** Connection closed", closeid)
			user.Delete(closeid)
			uid2conn[closeid] = nil, false
		}
	}
}

func (s *IRCd) manageIncoming() {
	defer s.running.Done()

	var open bool = true
	var c *conn.Conn
	for open {
		// Do anything?
		select {
		case c, open = <-s.Incoming:
			id := c.ID()
			switch len(id) {
			case 3:
				s.newServer <- c
			case 9:
				s.newClient <- c
			default:
				log.Printf("Warning: Unknown connection ID type: %s", id)
			}
		}
	}
}

var (
	// TODO(kevlar): Configurable?
	SendQ = 100
	RecvQ = 100
)

func (s *IRCd) Quit() {
	for i := 0; i < cap(s.shutdown); i++ {
		s.shutdown <- true
	}
	s.running.Wait()
}

func Start() {
	// Make sure the configuration is good before we do anything
	checkConfig()

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

	s := &IRCd{
		IRCd: &core.IRCd{
			ToClient: make(chan *parser.Message, SendQ),
			ToServer: make(chan *parser.Message, SendQ),
			Incoming: listener.Incoming,
		},

		clientClosing: make(chan string),
		newClient:     make(chan *conn.Conn),

		serverClosing: make(chan string),
		newServer:     make(chan *conn.Conn),

		shutdown: make(chan bool, 3),
		running:  new(sync.WaitGroup),
	}

	s.running.Add(1)
	go s.manageClients()

	s.running.Add(1)
	go s.manageServers()

	s.running.Add(1)
	go s.manageIncoming()

	s.running.Wait()
}

func Run() {
	LoadConfigString(DefaultXML)

	Start()
}

func isuid(id string) bool {
	return len(id) == 9 && id[0] >= '0' && id[0] <= '9'
}
