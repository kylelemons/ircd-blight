package server

import (
	"kevlar/ircd/parser"
	"kevlar/ircd/core"
)

var (
	pinghook = core.Register(parser.CMD_PING, core.Any, core.NArgs(1), Ping)
)

func Ping(hook string, msg *parser.Message, out chan<- *parser.Message) {
	pongmsg := msg.Args[0]
	out <- &parser.Message{
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
