package core

import (
	"kevlar/ircd/parser"
	"kevlar/ircd/user"
)

// Choose in what contexts a hook is called
type ExecutionMask int

const (
	Registration ExecutionMask = 1 << iota
	User
	Server
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
	Func        func(hook string, when ExecutionMask, message *parser.Message)
}

var (
	registeredHooks = map[string][]*Hook{}
)

func Register(hook string, when ExecutionMask, args CallConstraints,
fn func(string, ExecutionMask, *parser.Message)) *Hook {
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
func DispatchMessage(message *parser.Message) {
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
	case user.RegisteredAsServer:
		mask |= Server
	}
	for _, hook := range registeredHooks[hookName] {
		if hook.When&mask == mask {
			// TODO(kevlar): Check callconstraints
			go hook.Func(hookName, mask, message)
			hook.Calls++
		}
	}
}
