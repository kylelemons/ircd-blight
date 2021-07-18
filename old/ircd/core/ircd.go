package core

import (
	"sync"

	"github.com/kylelemons/ircd-blight/old/ircd/conn"
	"github.com/kylelemons/ircd-blight/old/ircd/parser"
)

type IRCd struct {
	Incoming      chan *conn.Conn
	newClient     chan *conn.Conn
	newServer     chan *conn.Conn
	clientClosing chan string
	serverClosing chan string

	ToClient   chan *parser.Message
	ToServer   chan *parser.Message
	fromClient chan *parser.Message
	fromServer chan *parser.Message

	running *sync.WaitGroup
}
