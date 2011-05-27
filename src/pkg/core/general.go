package core

import (
	"kevlar/ircd/parser"
)

var (
	pinghook = Register(parser.CMD_PING, Any, NArgs(1), Ping)
)

func Ping(hook string, msg *parser.Message, ircd *IRCd) {
	pongmsg := msg.Args[0]
	ircd.ToClient <- &parser.Message{
		Command: parser.CMD_PONG,
		Args: []string{
			Config.Name,
			pongmsg,
		},
		DestIDs: []string{
			msg.SenderID,
		},
	}
}
