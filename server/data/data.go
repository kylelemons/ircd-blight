// Package data contains common data types used in the IRC server.
package data

import (
	"sync"
)

// A Command is a command sent to a component of the IRC server.
type Command interface {
	String() string
	Done(error)
}

// A Server represents an IRC server to which this server is connected.
type Server struct {
	sync.Mutex

	SID  string
	Name string

	Control chan Command
}

// A ChanMode is a bitmask indicating what channel modes are set.
type ChanMode uint64

// A Channel represents a channel on this server.
type Channel struct {
	sync.Mutex

	Name string
	Mode ChanMode

	Control chan Command
}

// A UserMode is a bitmask indicating what user modes are set.
type UserMode uint64

// A User represents a user on this server.
type User struct {
	sync.Mutex

	UID  string
	Nick string
	User string
	Name string
	Mode UserMode

	Control chan Command
}

// A MemberMode is a bitmask indicating what modes a user has on a channel.
type MemberMode uint64

// Member mode constants
const (
	MemberVoice MemberMode = 1 << iota
	MemberHalfOp
	MemberOp
	MemberAdmin
)

// A Member stores data about a user's membership on a channel.
type Member struct {
	sync.Mutex

	User    *User
	Channel *Channel

	Mode MemberMode
}
