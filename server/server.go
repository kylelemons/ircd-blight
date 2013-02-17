package server

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/kylelemons/ircd-blight/server/data"
	"github.com/kylelemons/ircd-blight/server/grid"
)

// Length constants
const (
	ServerIDLen = 3 // [0-9][A-Z0-9]{2}
	UserIDLen   = 6 // [A-Z0-9]{6}
)

// Grid index constants
const (
	gridUser = iota
	gridChan
)

type Server struct {
	rw   sync.RWMutex
	grid grid.Grid

	sid     string
	nextUID uint64

	users map[string]*data.User // users[nick] = u, users[uid] = u
}

func NewServer(sid string) *Server {
	return &Server{
		sid:   sid,
		users: make(map[string]*data.User, 100),
	}
}

func (s *Server) signon(nick, user, name string) (uid string, err error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	// TODO(kevlar): tolower
	// TODO(kevlar): valid nick
	if _, ok := s.users[nick]; ok {
		// TODO(kevlar): log initial user
		return "", fmt.Errorf("NICK %q already in use", nick)
	}

	uid = s.sid + idstr(atomic.AddUint64(&s.nextUID, 1)-1, UserIDLen)
	u := &data.User{
		UID:  uid,
		Nick: nick,
		User: user,
		Name: name,
	}
	// TODO(kevlar): start command routine

	s.users[nick] = u
	s.users[uid] = u

	s.grid.Edges[gridUser].Touch(u.UID, u)

	log.Printf("[%s] Signon: %s!%s :%s", u.UID, u.Nick, u.User, u.Name)
	return uid, nil
}

func idstr(id uint64, length int) string {
	b := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		b[i], id = 'A'+byte(id%26), id/26
	}
	return string(b)
}
