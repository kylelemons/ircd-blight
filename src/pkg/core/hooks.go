package core

import (
	"kevlar/ircd/log"
	"kevlar/ircd/parser"
	"kevlar/ircd/server"
	"kevlar/ircd/user"
)

// Choose in what contexts a hook is called
type ExecutionMask int

const (
	Registration ExecutionMask = 1 << iota
	User
	Server
	Any ExecutionMask = Registration | User | Server
)

// Choose how many arguments a hook needs to be called
type CallConstraints struct {
	MinArgs int
	MaxArgs int
}

func NArgs(count int) CallConstraints {
	return CallConstraints{
		MinArgs: count,
		MaxArgs: count,
	}
}

func MinArgs(min int) CallConstraints {
	return CallConstraints{
		MinArgs: min,
		MaxArgs: -1,
	}
}

func OptArgs(required, optional int) CallConstraints {
	return CallConstraints{
		MinArgs: required,
		MaxArgs: required + optional,
	}
}

var (
	AnyArgs = CallConstraints{
		MinArgs: 0,
		MaxArgs: -1,
	}
)

// Allow registration of hooks in any module
type Hook struct {
	When        ExecutionMask
	Constraints CallConstraints
	Calls       int
	Func        func(hook string, message *parser.Message, ircd *IRCd)
}

var (
	registeredHooks = map[string][]*Hook{}
)

func Register(hook string, when ExecutionMask, args CallConstraints,
fn func(string, *parser.Message, *IRCd)) *Hook {
	if _, ok := registeredHooks[hook]; !ok {
		registeredHooks[hook] = make([]*Hook, 0, 1)
	}
	h := &Hook{
		When:        when,
		Constraints: args,
		Func:        fn,
	}
	registeredHooks[hook] = append(registeredHooks[hook], h)
	return h
}

// TODO(kevlar): Add channel to send messages back on
func DispatchClient(message *parser.Message, ircd *IRCd) {
	hookName := message.Command
	_, _, _, reg, ok := user.GetInfo(message.SenderID)
	if !ok {
		panic("Unknown user: " + message.SenderID)
	}
	var mask ExecutionMask
	switch reg {
	case user.Unregistered:
		mask |= Registration
	case user.RegisteredAsUser:
		mask |= User
	}
	for _, hook := range registeredHooks[hookName] {
		if hook.When&mask == mask {
			// TODO(kevlar): Check callconstraints
			go hook.Func(hookName, message, ircd)
			hook.Calls++
		}
	}
}

func DispatchServer(message *parser.Message, ircd *IRCd) {
	hookName := message.Command
	_, _, _, reg, ok := server.GetInfo(message.SenderID)
	if !ok {
		log.Warn.Printf("Unknown source server: %s", message.SenderID)
		return
	}
	var mask ExecutionMask
	switch reg {
	case server.Unregistered:
		mask |= Registration
	case server.RegisteredAsServer:
		mask |= Server
	}
	for _, hook := range registeredHooks[hookName] {
		if hook.When&mask == mask {
			// TODO(kevlar): Check callconstraints
			go hook.Func(hookName, message, ircd)
			hook.Calls++
		}
	}
}
