package core

import (
	"kevlar/ircd/log"
	"kevlar/ircd/parser"
	"kevlar/ircd/server"
)

var (
	pinghooks = []*Hook{
		Register(parser.CMD_PING, User, NArgs(1), Ping),
		Register(parser.CMD_PING, Server, OptArgs(1, 1), SPing),
		Register(parser.CMD_PONG, Server, OptArgs(1, 1), SPing),
	}
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

func SPing(hook string, msg *parser.Message, ircd *IRCd) {
	source := msg.Args[0]
	dest := Config.SID
	if len(msg.Args) > 1 {
		dest = msg.Args[1]
	}

	if dest == Config.SID {
		switch hook {
		case parser.CMD_PING:
			ircd.ToServer <- &parser.Message{
				Prefix:  Config.SID,
				Command: parser.CMD_PONG,
				Args: []string{
					Config.Name,
					source,
				},
				DestIDs: []string{
					msg.SenderID,
				},
			}
		case parser.CMD_PONG:
			log.Info.Printf("End of BURST from %s", source)
		}
	} else {
		for sid := range server.IterFor([]string{dest}, "") {
			log.Debug.Printf("Forwarding %s to %s", hook, sid)
			msg := msg.Dup()
			msg.DestIDs = []string{sid}
			ircd.ToServer <- msg
		}
	}
}
