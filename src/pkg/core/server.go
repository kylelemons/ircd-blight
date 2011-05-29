package core

import (
	"kevlar/ircd/conn"
	"kevlar/ircd/log"
	"kevlar/ircd/parser"
	"kevlar/ircd/server"
	"kevlar/ircd/user"
	"sync"
)

func CheckConfig() (okay bool) {
	okay = true

	if Config == nil {
		log.Error.Printf("No configuration loaded")
		return false
	}

	// Check hostname: require at least one .
	if !parser.ValidServerName(Config.Name) {
		log.Error.Printf("invalid server name %q: must match /\\w+(.\\w+)+/", Config.Name)
		okay = false
	}

	// Check prefix; [num][alphanum][alphanum]
	if !parser.ValidServerPrefix(Config.SID) {
		log.Error.Printf("invalid server [refix %q: must match /[0-9][0-9A-Z]{2}/")
		okay = false
	}
	user.UserIDPrefix = Config.SID

	// Check opers
	if len(Config.Operator) == 0 {
		log.Error.Printf("no operators defined: at least one required")
		okay = false
	}

	return
}

func (s *IRCd) manageServers() {
	defer s.running.Done()

	sid2conn := make(map[string]*conn.Conn)
	upstream := ""
	_ = upstream

	var open bool = true
	var msg *parser.Message
	for open {
		// Check if we are connected to our upstream
		// TODO(kevlar): upstream

		select {
		//// Incoming and outgoing messages to and from servers
		// Messages directly from connections
		case msg = <-s.fromServer:
			log.Debug.Printf("{%s} >> %s\n", msg.SenderID, msg)
			DispatchServer(msg, s)

		// Messages from hooks
		case msg, open = <-s.ToServer:
			// Count the number of messages sent
			sentcount := 0
			for _, dest := range msg.DestIDs {
				log.Debug.Printf("{%v} << %s\n", dest, msg)

				if conn, ok := sid2conn[dest]; ok {
					conn.WriteMessage(msg)
					sentcount++
				} else {
					log.Warn.Printf("Unknown SID %s", dest)
				}
			}

			if sentcount == 0 {
				log.Warn.Printf("Dropped outgoing server message: %s", msg)
			}

		//// Connection management
		// Connecting servers
		case conn := <-s.newServer:
			sid := conn.ID()
			sid2conn[sid] = conn
			server.Get(sid)
			log.Debug.Printf("{%s} ** Registered connection", sid)
			conn.Subscribe(s.fromServer)
			conn.SubscribeClose(s.serverClosing)
		// Disconnecting servers
		case closeid := <-s.serverClosing:
			log.Debug.Printf("{%s} ** Connection closed", closeid)
			// TODO(kevlar): Delete server
			user.Delete(closeid)
			sid2conn[closeid] = nil, false
		}
	}
}

func (s *IRCd) manageClients() {
	defer s.running.Done()

	uid2conn := make(map[string]*conn.Conn)

	var open bool = true
	var msg *parser.Message
	for open {
		// anything to do here?

		select {
		//// Incoming and outgoing messages to and from clients
		// Messages directly from connections
		//   TODO(kevlar): Collapse into one case (I don't like that server gets split)
		case msg = <-s.fromClient:
			u := user.Get(msg.SenderID)
			uid := u.ID()

			if msg.Command == parser.CMD_ERROR {
				user.Delete(uid)
				conn := uid2conn[uid]
				if conn != nil {
					log.Debug.Printf("[%s] ** Connection terminated remotely", uid)
					user.Delete(uid)
					uid2conn[uid] = nil, false
					conn.UnsubscribeClose(s.clientClosing)
					conn.Close()
				}
				continue
			}

			log.Debug.Printf("[%s] >> %s", uid, msg)
			DispatchClient(msg, s)

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
					log.Warn.Printf("Nonexistent ID %s as prefix", msg.Prefix)
				} else {
					msg.Prefix = nick + "!" + user + "@host" // TODO(kevlar): hostname
				}
			}
			for i := range msg.Args {
				if isuid(msg.Args[i]) {
					nick, _, _, _, ok := user.GetInfo(msg.Args[i])
					if !ok {
						log.Warn.Printf("Nonexistent ID %s as argument", msg.Args[i])
						continue
					}
					msg.Args[i] = nick
				}
			}

			for _, id := range msg.DestIDs {
				conn, ok := uid2conn[id]
				if !ok {
					log.Warn.Printf("Nonexistent ID %s in send", id)
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
				log.Debug.Printf("[%s] << %s\n", id, msg)
				sentcount++
				if closeafter {
					log.Debug.Printf("[%s] ** Connection terminated", id)
					user.Delete(id)
					uid2conn[id] = nil, false
					conn.UnsubscribeClose(s.clientClosing)
					conn.Close()
				}
			}
			if sentcount == 0 {
				log.Warn.Printf("Dropped outgoing client message: %s", msg)
			}

		// Connecting clients
		case conn := <-s.newClient:
			id := conn.ID()
			uid2conn[id] = conn
			user.Get(id)
			conn.Subscribe(s.fromClient)
			conn.SubscribeClose(s.clientClosing)
		// Disconnecting clients
		case closeid := <-s.clientClosing:
			log.Debug.Printf("[%s] ** Connection closed", closeid)
			user.Delete(closeid)
			uid2conn[closeid] = nil, false
		}
	}
}

func (s *IRCd) manageIncoming() {
	defer s.running.Done()

	quit := false
	defer func() {
		quit = true
	}()

	manage := func(conn *conn.Conn) {
		inc := make(chan *parser.Message)
		stop := make(chan string)
		conn.Subscribe(inc)
		conn.SubscribeClose(stop)

		user, nick := false, false
		pass, server, capab := false, false, false
		sid := ""

		queued := make([]*parser.Message, 0, 3)

		for !quit {
			select {
			case msg := <-inc:
				log.Debug.Printf(" %s  %s", msg.SenderID, msg)
				queued = append(queued, msg)
				switch msg.Command {
				case parser.CMD_PASS:
					if len(msg.Args) == 4 {
						pass = true
						sid = msg.Args[3]
					}
				case parser.CMD_USER:
					user = true
				case parser.CMD_NICK:
					nick = true
				case parser.CMD_CAPAB:
					capab = true
				case parser.CMD_SERVER:
					server = true
				}
			case <-stop:
				return
			}

			if !quit && nick && user {
				conn.Unsubscribe(inc)
				conn.UnsubscribeClose(stop)
				s.newClient <- conn
				for _, msg := range queued {
					s.fromClient <- msg
				}
				return
			}
			if !quit && pass && server && capab {
				conn.SetServer(sid)
				conn.Unsubscribe(inc)
				conn.UnsubscribeClose(stop)
				s.newServer <- conn
				for _, msg := range queued {
					msg.SenderID = sid
					s.fromServer <- msg
				}
				return
			}
		}
	}

	for {
		select {
		// Connecting clients
		case conn := <-s.Incoming:
			go manage(conn)
		}
	}
}

var (
	// TODO(kevlar): Configurable?
	SendQ = 100
	RecvQ = 100
)

func (s *IRCd) Quit() {
	close(s.Incoming)
	close(s.ToClient)
	close(s.ToServer)
	s.running.Wait()
}

func Start() {
	// Make sure the configuration is good before we do anything
	if !CheckConfig() {
		log.Error.Fatalf("Could not start: invalid configuration")
	}

	listener := conn.NewListener()
	defer listener.Close()
	for _, ports := range Config.Ports {
		portlist, err := ports.GetPortList()
		if err != nil {
			log.Warn.Print(err)
		}
		for _, port := range portlist {
			listener.AddPort(port)
		}
	}

	s := &IRCd{
		Incoming:      listener.Incoming,
		newClient:     make(chan *conn.Conn),
		newServer:     make(chan *conn.Conn),
		clientClosing: make(chan string),
		serverClosing: make(chan string),

		ToClient:   make(chan *parser.Message, SendQ),
		ToServer:   make(chan *parser.Message, SendQ),
		fromClient: make(chan *parser.Message, SendQ),
		fromServer: make(chan *parser.Message, SendQ),

		running: new(sync.WaitGroup),
	}

	s.running.Add(1)
	go s.manageClients()

	s.running.Add(1)
	go s.manageServers()

	s.running.Add(1)
	go s.manageIncoming()

	s.running.Wait()
}

func isuid(id string) bool {
	return len(id) == 9 && id[0] >= '0' && id[0] <= '9'
}
