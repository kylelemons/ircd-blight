package core

import (
	"kevlar/ircd/parser"
	"kevlar/ircd/channel"
	//"kevlar/ircd/user"
	"kevlar/ircd/server"
	"kevlar/ircd/log"
	"strings"
)

var (
	joinhooks = []*Hook{
		Register(parser.CMD_JOIN, User, OptArgs(1, 1), Join),
		Register(parser.CMD_JOIN, Server, MinArgs(3), SJoin),
		Register(parser.CMD_SJOIN, Server, MinArgs(4), SJoin),
	}
)

func Join(hook string, msg *parser.Message, ircd *IRCd) {
	for _, channame := range strings.Split(msg.Args[0], ",", -1) {
		channel, err := channel.Get(channame, true)
		if num, ok := err.(*parser.Numeric); ok {
			ircd.ToClient <- num.Message(msg.SenderID)
			continue
		}

		notify, err := channel.Join(msg.SenderID)
		if num, ok := err.(*parser.Numeric); ok {
			ircd.ToClient <- num.Message(msg.SenderID)
			continue
		}

		// Forward to other servers
		for sid := range server.Iter() {
			ircd.ToServer <- &parser.Message{
				Prefix:  msg.SenderID,
				Command: parser.CMD_JOIN,
				Args: []string{
					channel.TS(),
					channel.Name(),
					"+",
				},
				DestIDs: []string{sid},
			}
		}

		ircd.ToClient <- &parser.Message{
			Prefix:  msg.SenderID,
			Command: parser.CMD_JOIN,
			Args: []string{
				channel.Name(),
			},
			DestIDs: notify,
		}

		ircd.ToClient <- channel.NamesMessage(msg.SenderID)
	}
}

func SJoin(hook string, msg *parser.Message, ircd *IRCd) {
	chanTS, channame, mode, uids := msg.Args[0], msg.Args[1], msg.Args[2], msg.Args[3:]
	_ = chanTS
	_ = mode

	// Forward on to other servers
	for fwd := range server.Iter() {
		if fwd != msg.SenderID {
			log.Debug.Printf("Forwarding SJOIN from %s to %s", msg.SenderID, fwd)
			fmsg := msg.Dup()
			fmsg.DestIDs = []string{fwd}
		}
	}

	for i, uid := range uids {
		uids[i] = uid[len(uid)-9:]
	}

	channel, err := channel.Get(channame, true)
	if num, ok := err.(*parser.Numeric); ok {
		ircd.ToServer <- num.Error(msg.SenderID)
		return
	}

	if len(uids) == 0 && len(msg.Prefix) == 9 {
		uids = []string{msg.Prefix}
	}

	chanusers, err := channel.Join(uids...)
	if num, ok := err.(*parser.Numeric); ok {
		ircd.ToServer <- num.Error(msg.SenderID)
		return
	}

	notify := []string{}
	for _, uid := range chanusers {
		if uid[:3] == Config.SID {
			notify = append(notify, uid)
		}
	}

	for _, joiner := range uids {
		ircd.ToClient <- &parser.Message{
			Prefix:  joiner,
			Command: parser.CMD_JOIN,
			Args: []string{
				channel.Name(),
			},
			DestIDs: notify,
		}
	}
}
