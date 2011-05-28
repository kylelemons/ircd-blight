package core

import (
	"kevlar/ircd/parser"
	"kevlar/ircd/channel"
	"kevlar/ircd/user"
	"kevlar/ircd/server"
	"kevlar/ircd/log"
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
	sender := msg.SenderID
	if len(msg.Prefix) == 9 {
		sender = msg.Prefix
	}
	local := []string{}
	remote := []string{}
	for _, name := range recipients {
		if parser.ValidChannel(name) {
			channel, err := channel.Get(name, false)
			if num, ok := err.(*parser.Numeric); ok {
				if !quiet {
					ircd.ToClient <- num.Message(msg.SenderID)
				}
				continue
			}
			local := []string{}
			remote := []string{}
			for _, uid := range channel.UserIDs() {
				if uid != sender {
					if uid[:3] == Config.SID {
						local = append(local, uid)
					} else {
						remote = append(remote, uid)
					}
				}
			}
			if len(remote) > 0 {
				for sid := range server.IterFor(remote, msg.SenderID) {
					log.Debug.Printf("Forwarding PRIVMSG from %s to %s", msg.SenderID, sid)
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
			}
			if len(local) > 0 {
				ircd.ToClient <- &parser.Message{
					Prefix:  sender,
					Command: hook,
					Args: []string{
						channel.Name(),
						text,
					},
					DestIDs: local,
				}
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
		if id[:3] == Config.SID {
			local = append(local, id)
		} else {
			remote = append(remote, id)
		}
	}
	if len(remote) > 0 {
		for _, remoteid := range remote {
			for sid := range server.IterFor([]string{remoteid}, "") {
				ircd.ToServer <- &parser.Message{
					Prefix:  sender,
					Command: hook,
					Args: []string{
						remoteid,
						text,
					},
					DestIDs: []string{sid},
				}
			}
		}
	}
	if len(local) > 0 {
		ircd.ToClient <- &parser.Message{
			Prefix:  sender,
			Command: hook,
			Args: []string{
				"*",
				text,
			},
			DestIDs: local,
		}
	}
}
