package server

import (
	"log"
	"os"
	"kevlar/ircd/parser"
	"kevlar/ircd/user"
	"kevlar/ircd/core"
)

var (
	reghooks = []*core.Hook{
		core.Register(parser.CMD_NICK, core.Registration, core.AnyArgs, Registration),
		core.Register(parser.CMD_USER, core.Registration, core.AnyArgs, Registration),
		core.Register(parser.CMD_SERVER, core.Registration, core.AnyArgs, Registration),
		core.Register(parser.CMD_PASS, core.Registration, core.AnyArgs, Registration),
	}
	quithook = core.Register(parser.CMD_QUIT, core.Any, core.AnyArgs, Quit)
)

// Handle the NICK, USER, SERVER, and PASS messages
func Registration(hook string, msg *parser.Message, ircd *core.IRCd) {
	u := user.Get(msg.SenderID)

	var err os.Error

	switch msg.Command {
	case parser.CMD_NICK:
		// NICK <nick>
		nick := msg.Args[0]
		err = u.SetNick(nick)
	case parser.CMD_USER:
		// USER <user> . . :<real name>
		username, realname := msg.Args[0], msg.Args[3]
		err = u.SetUser(username, realname)
	}

	if num, ok := err.(*parser.Numeric); ok {
		msg := num.Message()
		msg.DestIDs = append(msg.DestIDs, u.ID())
		ircd.ToClient <- msg
		return
	}

	nickname, username, _, _ := u.Info()
	if nickname == "*" || username == "" {
		return
	}

	log.Printf("[%s] ** Registered\n", u.ID())
	u.SetType(user.RegisteredAsUser)

	destIDs := []string{u.ID()}
	// RPL_WELCOME
	msg = parser.NewNumeric(parser.RPL_WELCOME).Message()
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

func Quit(hook string, msg *parser.Message, ircd *core.IRCd) {
	reason := "Client Quit"
	if len(msg.Args) > 0 {
		reason = msg.Args[0]
	}
	error := &parser.Message{
		Command: parser.CMD_ERROR,
		Args: []string{
			"Closing Link (" + reason + ")",
		},
		DestIDs: []string{
			msg.SenderID,
		},
	}
	ircd.ToClient <- error
}
