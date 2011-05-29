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
		Register(parser.CMD_PART, User, OptArgs(1, 1), Part),
		Register(parser.CMD_PART, Server, NArgs(2), SPart),
	}
)

// Local joins only
func Join(hook string, msg *parser.Message, ircd *IRCd) {
	// todo keys
	for _, channame := range strings.Split(msg.Args[0], ",", -1) {
		channel, err := channel.Get(channame, true)
		if num, ok := err.(*parser.Numeric); ok {
			ircd.ToClient <- num.Message(msg.SenderID)
			continue
		}

		members, err := channel.Join(msg.SenderID)
		if num, ok := err.(*parser.Numeric); ok {
			ircd.ToClient <- num.Message(msg.SenderID)
			continue
		}

		notify := []string{}
		for _, uid := range members {
			if uid[:3] == Config.SID {
				notify = append(notify, uid)
			}
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

		if len(notify) > 0 {
			ircd.ToClient <- &parser.Message{
				Prefix:  msg.SenderID,
				Command: parser.CMD_JOIN,
				Args: []string{
					channel.Name(),
				},
				DestIDs: notify,
			}
		}

		ircd.ToClient <- channel.NamesMessage(msg.SenderID)
	}
}

// Local PARTs only
func Part(hook string, msg *parser.Message, ircd *IRCd) {
	reason := "Leaving"
	if len(msg.Args) > 1 {
		reason = msg.Args[1]
	}

	leftchans := []string{}

	for _, channame := range strings.Split(msg.Args[0], ",", -1) {
		channel, err := channel.Get(channame, false)
		if num, ok := err.(*parser.Numeric); ok {
			ircd.ToClient <- num.Message(msg.SenderID)
			continue
		}

		members, err := channel.Part(msg.SenderID)
		if num, ok := err.(*parser.Numeric); ok {
			ircd.ToClient <- num.Message(msg.SenderID)
			continue
		}

		notify := []string{}
		for _, uid := range members {
			if uid[:3] == Config.SID {
				notify = append(notify, uid)
			}
		}

		if len(notify) > 0 {
			ircd.ToClient <- &parser.Message{
				Prefix:  msg.SenderID,
				Command: parser.CMD_PART,
				Args: []string{
					channel.Name(),
				},
				DestIDs: notify,
			}
		}

		leftchans = append(leftchans, channame)
	}

	// Forward to other servers
	if len(leftchans) > 0 {
		leftstr := strings.Join(leftchans, ",")
		for sid := range server.Iter() {
			ircd.ToServer <- &parser.Message{
				Prefix:  msg.SenderID,
				Command: parser.CMD_PART,
				Args: []string{
					leftstr,
					reason,
				},
				DestIDs: []string{sid},
			}
		}
	}
}

// Server JOIN and SJOIN
func SJoin(hook string, msg *parser.Message, ircd *IRCd) {
	chanTS, channame, mode := msg.Args[0], msg.Args[1], msg.Args[2]

	uids := []string{msg.Prefix}
	if len(msg.Prefix) == 3 {
		if len(msg.Args) == 3 {
			return
		}
		uids = strings.Split(msg.Args[3], " ", -1)
	}

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

	if len(notify) > 0 {
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
}

// Server PART
func SPart(hook string, msg *parser.Message, ircd *IRCd) {
	chanlist, reason := strings.Split(msg.Args[0], ",", -1), msg.Args[1]

	// Forward on to other servers
	for fwd := range server.Iter() {
		if fwd != msg.SenderID {
			log.Debug.Printf("Forwarding PART from %s to %s", msg.SenderID, fwd)
			fmsg := msg.Dup()
			fmsg.DestIDs = []string{fwd}
		}
	}

	for _, channame := range chanlist {
		channel, err := channel.Get(channame, false)
		if num, ok := err.(*parser.Numeric); ok {
			ircd.ToServer <- num.Error(msg.SenderID)
			return
		}

		chanusers, err := channel.Part(msg.Prefix)
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

		if len(notify) > 0 {
			ircd.ToClient <- &parser.Message{
				Prefix:  msg.Prefix,
				Command: parser.CMD_PART,
				Args: []string{
					channel.Name(),
					reason,
				},
				DestIDs: notify,
			}
		}
	}
}
