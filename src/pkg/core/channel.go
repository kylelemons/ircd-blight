package core

import (
	"kevlar/ircd/parser"
	"kevlar/ircd/channel"
	"kevlar/ircd/user"
	"strings"
)

var (
	joinhook = Register(parser.CMD_JOIN, User, OptArgs(1, 1), Join)
)

func Join(hook string, msg *parser.Message, ircd *IRCd) {
	nickname, username, _, _, _ := user.GetInfo(msg.SenderID)

	for _, channame := range strings.Split(msg.Args[0], ",", -1) {
		channel, err := channel.Get(channame, true)
		if num, ok := err.(*parser.Numeric); ok {
			ircd.ToClient <- num.Message(msg.SenderID)
			continue
		}

		notify, err := channel.Join(msg.SenderID, nickname+username) // TODO(kevlar): hostname
		if num, ok := err.(*parser.Numeric); ok {
			ircd.ToClient <- num.Message(msg.SenderID)
			continue
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
