package core

import (
	"kevlar/ircd/parser"
	"kevlar/ircd/channel"
	"kevlar/ircd/user"
	"strings"
)

var (
	msghooks = []*Hook{
		Register(parser.CMD_PRIVMSG, User, NArgs(2), Privmsg),
		Register(parser.CMD_NOTICE, User, NArgs(2), Privmsg),
	}
)

func Privmsg(hook string, msg *parser.Message, ircd *IRCd) {
	quiet := hook == parser.CMD_NOTICE
	recipients, text := strings.Split(msg.Args[0], ",", -1), msg.Args[1]
	destIDs := make([]string, 0, len(recipients))
	for _, name := range recipients {
		if parser.ValidChannel(name) {
			channel, err := channel.Get(name, false)
			if num, ok := err.(*parser.Numeric); ok {
				if !quiet {
					ircd.ToClient <- num.Message(msg.SenderID)
				}
				continue
			}
			userids := channel.UserIDs()
			for i := 0; i < len(userids); i++ {
				if userids[i] == msg.SenderID {
					userids[i] = userids[len(userids)-1]
					userids = userids[:len(userids)-1]
				}
			}
			ircd.ToClient <- &parser.Message{
				Prefix:  msg.SenderID,
				Command: hook,
				Args: []string{
					channel.Name(),
					text,
				},
				DestIDs: userids,
			}
			continue
		}

		id, err := user.GetID(name)
		if num, ok := err.(*parser.Numeric); ok {
			if !quiet {
				ircd.ToClient <- num.Message(msg.SenderID)
			}
			continue
		}
		destIDs = append(destIDs, id)
	}
	if len(destIDs) > 0 {
		ircd.ToClient <- &parser.Message{
			Prefix:  msg.SenderID,
			Command: hook,
			Args: []string{
				"*",
				text,
			},
			DestIDs: destIDs,
		}
	}
}
