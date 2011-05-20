package core

import (
	"kevlar/ircd/parser"
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
	AnyArgs = &CallConstraints{
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
		registeredHooks[hook] = make([]*Hook, 1)
	}
	h := &Hook{
		When:        when,
		Constraints: args,
		Func:        fn,
	}
	registeredHooks[hook] = append(registeredHooks[hook], h)
	return h
}

// TODO(kevlar): Add source? Something to send numerics back on
func Dispatch(hookName string, mask ExecutionMask, message *parser.Message) {
	for _, hook := range registeredHooks[hookName] {
		if hook.When|mask == mask {
			// TODO(kevlar): Check callconstraints
			go hook.Func(hookName, mask, message)
			hook.Calls++
		}
	}
}
