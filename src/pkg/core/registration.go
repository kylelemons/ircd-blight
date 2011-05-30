package core

import (
	"kevlar/ircd/channel"
	"kevlar/ircd/log"
	"kevlar/ircd/parser"
	"kevlar/ircd/server"
	"kevlar/ircd/user"
	"os"
	"strings"
)

var (
	reghooks = []*Hook{
		Register(parser.CMD_NICK, Registration, AnyArgs, ConnReg),
		Register(parser.CMD_USER, Registration, AnyArgs, ConnReg),
		Register(parser.CMD_SERVER, Registration, AnyArgs, ConnReg),
		Register(parser.CMD_PASS, Registration, AnyArgs, ConnReg),
		Register(parser.CMD_CAPAB, Registration, AnyArgs, ConnReg),
		Register(parser.CMD_UID, Server, NArgs(9), Uid),
		Register(parser.CMD_SID, Server, NArgs(4), Sid),
	}
	quithooks = []*Hook{
		Register(parser.CMD_QUIT, User, AnyArgs, Quit),
		Register(parser.CMD_QUIT, Server, NArgs(2), Quit),
		Register(parser.CMD_SQUIT, Server, NArgs(2), SQuit),
	}
)

// Handle the NICK, USER, SERVER, and PASS messages
func ConnReg(hook string, msg *parser.Message, ircd *IRCd) {
	var err os.Error
	var u *user.User
	var s *server.Server

	switch len(msg.SenderID) {
	case 3:
		s = server.Get(msg.SenderID, true)
	case 9:
		u = user.Get(msg.SenderID)
	}

	switch msg.Command {
	case parser.CMD_NICK:
		// NICK <nick>
		if u != nil {
			nick := msg.Args[0]
			err = u.SetNick(nick)
		}
	case parser.CMD_USER:
		// USER <user> . . :<real name>
		if u != nil {
			username, realname := msg.Args[0], msg.Args[3]
			err = u.SetUser(username, realname)
		}
	case parser.CMD_PASS:
		if s != nil {
			if len(msg.Args) != 4 {
				return
			}
			// PASS <password> TS <ver> <pfx>
			err = s.SetPass(msg.Args[0], msg.Args[2], msg.Args[3])
		}
	case parser.CMD_CAPAB:
		if s != nil {
			err = s.SetCapab(msg.Args[0])
		}
	case parser.CMD_SERVER:
		if s != nil {
			err = s.SetServer(msg.Args[0], msg.Args[1])
		}
	default:
		log.Warn.Printf("Unknown command %q", msg)
	}

	if u != nil {
		if err != nil {
			switch err := err.(type) {
			case *parser.Numeric:
				msg := err.Message()
				msg.DestIDs = append(msg.DestIDs, u.ID())
				ircd.ToClient <- msg
				return
			default:
				msg := &parser.Message{
					Command: parser.CMD_ERROR,
					Args:    []string{err.String()},
					DestIDs: []string{u.ID()},
				}
				ircd.ToClient <- msg
				return
			}
		}

		nickname, username, realname, _ := u.Info()
		if nickname != "*" && username != "" {
			// Notify servers
			for sid := range server.Iter() {
				ircd.ToServer <- &parser.Message{
					Prefix:  Config.SID,
					Command: parser.CMD_UID,
					Args: []string{
						nickname,
						"1",
						u.TS(),
						"+i",
						username,
						"some.host",
						"127.0.0.1",
						u.ID(),
						realname,
					},
					DestIDs: []string{sid},
				}
			}

			// Process signon
			sendSignon(u, ircd)
			return
		}
	}

	if s != nil {
		if err != nil {
			switch err := err.(type) {
			case *parser.Numeric:
				msg := err.Message()
				msg.DestIDs = append(msg.DestIDs, s.ID())
				ircd.ToServer <- msg
				return
			default:
				msg := &parser.Message{
					Command: parser.CMD_ERROR,
					Args:    []string{err.String()},
					DestIDs: []string{s.ID()},
				}
				ircd.ToServer <- msg
				return
			}
		}

		sid, serv, pass, capab := s.Info()
		if sid != "" && serv != "" && pass != "" && len(capab) > 0 {
			// Notify servers
			for sid := range server.Iter() {
				ircd.ToServer <- &parser.Message{
					Prefix:  Config.SID,
					Command: parser.CMD_SID,
					Args: []string{
						serv,
						"2",
						sid,
						"some server",
					},
					DestIDs: []string{sid},
				}
			}

			sendServerSignon(s, ircd)
			Burst(s, ircd)
		}
	}
}

func sendSignon(u *user.User, ircd *IRCd) {
	log.Info.Printf("[%s] ** Registered\n", u.ID())
	u.SetType(user.RegisteredAsUser)

	destIDs := []string{u.ID()}
	// RPL_WELCOME
	msg := parser.NewNumeric(parser.RPL_WELCOME).Message()
	msg.Args[1] = "Welcome to the " + Config.Network.Name + " network, " + u.Nick() + "!"
	msg.DestIDs = destIDs
	ircd.ToClient <- msg

	// RPL_YOURHOST
	msg = parser.NewNumeric(parser.RPL_YOURHOST).Message()
	msg.Args[1] = "Your host is " + Config.Name + ", running IRCD-Blight" // TODO(kevlar): Version
	msg.DestIDs = destIDs
	ircd.ToClient <- msg

	// RPL_CREATED
	// RPL_MYINFO
	// RPL_ISUPPORT

	// RPL_LUSERCLIENT
	// RPL_LUSEROP
	// RPL_LUSERUNKNOWN
	// RPL_LUSERCHANNELS
	// RPL_LUSERME

	// RPL_LOCALUSERS
	// RPL_GLOBALUSERS

	// RPL_NOMOTD
	msg = parser.NewNumeric(parser.ERR_NOMOTD).Message()
	msg.DestIDs = destIDs
	ircd.ToClient <- msg

	msg = &parser.Message{
		Command: parser.CMD_MODE,
		Prefix:  "*",
		Args: []string{
			"*",
			"+i",
		},
		DestIDs: destIDs,
	}
	ircd.ToClient <- msg
}

func sendServerSignon(s *server.Server, ircd *IRCd) {
	log.Info.Printf("{%s} ** Registered As Server\n", s.ID())
	s.SetType(server.RegisteredAsServer)

	destIDs := []string{s.ID()}

	var msg *parser.Message

	msg = &parser.Message{
		Command: parser.CMD_PASS,
		Args: []string{
			"testpass", // TODO
			"TS",
			"6",
			Config.SID,
		},
		DestIDs: destIDs,
	}
	ircd.ToServer <- msg

	msg = &parser.Message{
		Command: parser.CMD_CAPAB,
		Args: []string{
			//"QS EX CHW IE KLN KNOCK TB UNKLN CLUSTER ENCAP SERVICES RSFNC SAVE EUID EOPMOD BAN MLOCK",
			"QS ENCAP", // TODO
		},
		DestIDs: destIDs,
	}
	ircd.ToServer <- msg

	msg = &parser.Message{
		Command: parser.CMD_SERVER,
		Args: []string{
			Config.Name,
			"1",
			"IRCd",
		},
		DestIDs: destIDs,
	}
	ircd.ToServer <- msg
}

func Burst(serv *server.Server, ircd *IRCd) {
	destIDs := []string{serv.ID()}
	sid := Config.SID
	var msg *parser.Message

	// SID/SERVER
	// UID/EUID
	for uid := range user.Iter() {
		u := user.Get(uid)
		nick, username, name, typ := u.Info()
		if typ != user.RegisteredAsUser {
			continue
		}
		msg = &parser.Message{
			Prefix:  sid,
			Command: parser.CMD_UID,
			Args: []string{
				nick,
				// hopcount
				"1",
				u.TS(),
				// umodes
				"+i",
				username,
				// visible hostname
				"some.host",
				// IP addr
				"127.0.0.1",
				uid,
				name,
			},
			DestIDs: destIDs,
		}
		ircd.ToServer <- msg
	}
	// Optional: ENCAP REALHOST, ENCAP LOGIN, AWAY
	// SJOIN
	for channame := range channel.Iter() {
		chanobj, _ := channel.Get(channame, false)
		msg = &parser.Message{
			Prefix:  sid,
			Command: parser.CMD_SJOIN,
			Args: []string{
				chanobj.TS(),
				channame,
				// modes, params...
				"+", // "+nt",
				strings.Join(chanobj.UserIDsWithPrefix(), " "),
			},
			DestIDs: destIDs,
		}
		ircd.ToServer <- msg
	}
	// Optional: BMAST
	// Optional: TB
}

func Uid(hook string, msg *parser.Message, ircd *IRCd) {
	nickname, hopcount, nickTS := msg.Args[0], msg.Args[1], msg.Args[2]
	umode, username, hostname := msg.Args[3], msg.Args[4], msg.Args[5]
	ip, uid, name := msg.Args[6], msg.Args[7], msg.Args[8]
	_ = umode

	err := user.Import(uid, nickname, username, hostname, ip, hopcount, nickTS, name)
	if err != nil {
		// TODO: TS check - Kill remote or local? For now, we kill remote.
		ircd.ToServer <- &parser.Message{
			Prefix:  Config.SID,
			Command: parser.CMD_SQUIT,
			Args: []string{
				uid,
				err.String(),
			},
			DestIDs: []string{msg.SenderID},
		}
	}

	for fwd := range server.Iter() {
		if fwd != msg.SenderID {
			log.Debug.Printf("Forwarding UID from %s to %s", msg.SenderID, fwd)
			fmsg := msg.Dup()
			fmsg.DestIDs = []string{fwd}
		}
	}
}

func Sid(hook string, msg *parser.Message, ircd *IRCd) {
	servname, hopcount, sid, desc := msg.Args[0], msg.Args[1], msg.Args[2], msg.Args[3]

	err := server.Link(msg.Prefix, sid, servname, hopcount, desc)
	if err != nil {
		ircd.ToServer <- &parser.Message{
			Prefix:  Config.SID,
			Command: parser.CMD_SQUIT,
			Args: []string{
				sid,
				err.String(),
			},
			DestIDs: []string{msg.SenderID},
		}
	}

	for fwd := range server.Iter() {
		if fwd != msg.SenderID {
			log.Debug.Printf("Forwarding SID from %s to %s", msg.SenderID, fwd)
			fmsg := msg.Dup()
			fmsg.DestIDs = []string{fwd}
		}
	}
}

func Quit(hook string, msg *parser.Message, ircd *IRCd) {
	quitter := msg.SenderID
	reason := "Client Quit"

	if len(msg.Args) > 0 {
		reason = msg.Args[0]
	}

	if len(msg.SenderID) == 3 {
		quitter = msg.Prefix
	}

	for sid := range server.Iter() {
		log.Debug.Printf("Forwarding QUIT from %s to %s", quitter, sid)
		if sid != msg.SenderID {
			ircd.ToServer <- &parser.Message{
				Prefix:  quitter,
				Command: parser.CMD_QUIT,
				Args: []string{
					reason,
				},
				DestIDs: []string{sid},
			}
		}
	}

	members := channel.PartAll(quitter)
	log.Debug.Printf("QUIT recipients: %#v", members)
	peers := make(map[string]bool)
	for _, users := range members {
		for _, uid := range users {
			if uid[:3] == Config.SID && uid != quitter {
				peers[uid] = true
			}
		}
	}
	if len(peers) > 0 {
		notify := []string{}
		for peer := range peers {
			notify = append(notify, peer)
		}
		ircd.ToClient <- &parser.Message{
			Prefix:  quitter,
			Command: parser.CMD_QUIT,
			Args: []string{
				"Quit: " + reason,
			},
			DestIDs: notify,
		}
	}

	// Will be dropped if it's a remote client
	error := &parser.Message{
		Command: parser.CMD_ERROR,
		Args: []string{
			"Closing Link (" + reason + ")",
		},
		DestIDs: []string{
			quitter,
		},
	}
	ircd.ToClient <- error
}

func SQuit(hook string, msg *parser.Message, ircd *IRCd) {
	split, reason := msg.Args[0], msg.Args[1]

	if split == Config.SID {
		split = msg.SenderID
	}

	// Forward
	for sid := range server.Iter() {
		if sid != msg.SenderID {
			msg := msg.Dup()
			msg.DestIDs = []string{sid}
			ircd.ToServer <- msg
		}
	}
	if server.IsLocal(split) {
		ircd.ToServer <- &parser.Message{
			Command: parser.CMD_ERROR,
			Args: []string{
				"SQUIT: " + reason,
			},
		}
	}

	sids := server.Unlink(split)
	peers := user.Netsplit(sids)
	notify := channel.Netsplit(Config.SID, peers)

	log.Debug.Printf("NET SPLIT: %s", split)
	log.Debug.Printf(" -   SIDs: %v", sids)
	log.Debug.Printf(" -  Peers: %v", peers)
	log.Debug.Printf(" - Notify: %v", notify)

	for uid, peers := range notify {
		ircd.ToClient <- &parser.Message{
			Prefix:  uid,
			Command: parser.CMD_QUIT,
			Args: []string{
				"*.net *.split",
			},
			DestIDs: peers,
		}
	}
}
