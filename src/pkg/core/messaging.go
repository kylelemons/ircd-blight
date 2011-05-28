package core

import (
	"kevlar/ircd/parser"
	"kevlar/ircd/channel"
	"kevlar/ircd/user"
	"strings"
)

var (
	msghooks = []*Hook{
		Register(parser.CMD_PRIVMSG, User|Server, NArgs(2), Privmsg),
		Register(parser.CMD_NOTICE, User|Server, NArgs(2), Privmsg),
	}
)

func Privmsg(hook string, msg *parser.Message, ircd *IRCd) {
	quiet := hook == parser.CMD_NOTICE
	recipients, text := strings.Split(msg.Args[0], ",", -1), msg.Args[1]
	destIDs := make([]string, 0, len(recipients))
	sender := msg.SenderID
	if len(msg.Prefix) == 9 {
		sender = msg.Prefix
	}
	for _, name := range recipients {
		if parser.ValidChannel(name) {
			channel, err := channel.Get(name, false)
			if num, ok := err.(*parser.Numeric); ok {
				if !quiet {
					ircd.ToClient <- num.Message(msg.SenderID)
				}
				continue
			}
			userids := []string{}
			sids := make(map[string]bool)
			for _, uid := range channel.UserIDs() {
				if uid != sender && uid[:3] != msg.SenderID {
					if uid[:3] == Config.SID {
						userids = append(userids, uid)
					} else {
						sids[uid[:3]] = true
					}
				}
			}
			for sid := range sids {
				ircd.ToServer <- &parser.Message{
					Prefix:  sender,
					Command: hook,
					Args: []string{
						channel.Name(),
						text,
					},
					DestIDs: []string{sid},
				}
			}
			ircd.ToClient <- &parser.Message{
				Prefix:  sender,
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
			Prefix:  sender,
			Command: hook,
			Args: []string{
				"*",
				text,
			},
			DestIDs: destIDs,
		}
	}
}
