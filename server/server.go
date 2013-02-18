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

	users map[string]*data.User    // users[nick] = u, users[uid] = u
	chans map[string]*data.Channel // chans[channel] = c
}

func NewServer(sid string) *Server {
	return &Server{
		sid:   sid,
		users: make(map[string]*data.User, 100),
		chans: make(map[string]*data.Channel, 100),
	}
}

func (s *Server) signon(nick, user, name string) (*data.User, error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	// TODO(kevlar): tolower
	// TODO(kevlar): valid nick
	if _, ok := s.users[nick]; ok {
		// TODO(kevlar): log initial user
		return nil, fmt.Errorf("NICK %q already in use", nick)
	}

	uid := s.sid + idstr(atomic.AddUint64(&s.nextUID, 1)-1, UserIDLen)
	u := &data.User{
		UID:  uid,
		Nick: nick,
		User: user,
		Name: name,
	}
	// TODO(kevlar): start command routine?

	s.users[nick] = u
	s.users[uid] = u

	s.grid.Edges[gridUser].Touch(u.UID, u)

	log.Printf("[%s] Signon: %s!%s :%s", u.UID, u.Nick, u.User, u.Name)
	return u, nil
}

func (s *Server) join(uid, channel string) (*data.Member, error) {
	s.rw.Lock()
	defer s.rw.Unlock()

	u, ok := s.users[uid]
	if !ok {
		// TODO(kevlar): log stack trace
		return nil, fmt.Errorf("UID %q does not exist", uid)
	}

	// TODO(kevlar): tolower
	// TODO(kevlar): valid chan
	c, chanExist := s.chans[channel]
	if !chanExist {
		c = &data.Channel{
			Name: channel,
		}
		s.grid.Edges[gridChan].Touch(c.Name, c)
		s.chans[channel] = c
	}

	member := &data.Member{
		User:    u,
		Channel: c,
	}
	if !chanExist {
		member.Mode |= data.MemberOp | data.MemberAdmin
	}

	if added := s.grid.Insert([2]string{u.UID, c.Name}, member); !added {
		return nil, fmt.Errorf("UID %q is already on %s", u.UID, c.Name)
	}
	return member, nil
}

func idstr(id uint64, length int) string {
	b := make([]byte, length)
	for i := length - 1; i >= 0; i-- {
		b[i], id = 'A'+byte(id%26), id/26
	}
	return string(b)
}
