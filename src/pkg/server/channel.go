package server

import (
	"kevlar/ircd/parser"
	"kevlar/ircd/core"
	"kevlar/ircd/channel"
	"kevlar/ircd/user"
	"strings"
)

var (
	joinhook = core.Register(parser.CMD_JOIN, core.User, core.OptArgs(1, 1), Join)
)

func Join(hook string, msg *parser.Message, out chan<- *parser.Message) {
	nickname, username, _, _, _ := user.GetInfo(msg.SenderID)

	for _, channame := range strings.Split(msg.Args[0], ",", -1) {
		channel, err := channel.Get(channame, true)
		if num, ok := err.(*parser.Numeric); ok {
			out <- num.Message(msg.SenderID)
			continue
		}

		notify, err := channel.Join(msg.SenderID, nickname+username) // TODO(kevlar): hostname
		if num, ok := err.(*parser.Numeric); ok {
			out <- num.Message(msg.SenderID)
			continue
		}

		out <- &parser.Message{
			Prefix:  msg.SenderID,
			Command: parser.CMD_JOIN,
			Args: []string{
				channel.Name(),
			},
			DestIDs: notify,
		}

		out <- channel.NamesMessage(msg.SenderID)
	}
}
