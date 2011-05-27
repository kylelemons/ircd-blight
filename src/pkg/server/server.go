package server

import (
	"kevlar/ircd/parser"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	servMap   = make(map[string]*Server)
	servMutex = new(sync.RWMutex)
)

type servType int

const (
	Unregistered servType = iota
	RegisteredAsServer
)

type Server struct {
	mutex  *sync.RWMutex
	id     string
	ts     int64
	styp   servType
	pass   string
	sver   int
	server string
	capab  []string
	hops   int
}

func (s *Server) ID() string {
	return s.id
}

func (s *Server) Type() servType {
	return s.styp
}

// Atomically get all of the server's information.
func (s *Server) Info() (sid, server, pass string, capab []string) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.id, s.server, s.pass, s.capab
}

func Get(id string) *Server {
	servMutex.Lock()
	defer servMutex.Unlock()

	if s, ok := servMap[id]; ok {
		return s
	}

	s := &Server{
		mutex: new(sync.RWMutex),
		id:    id,
	}
	servMap[id] = s
	return s
}


// Atomically get all of the server's information.
func GetInfo(id string) (sid, server string, capab []string, typ servType, ok bool) {
	servMutex.Lock()
	defer servMutex.Unlock()

	var s *Server
	if s, ok = servMap[id]; !ok {
		return
	}

	return s.id, s.server, s.capab, s.styp, true
}

func Iter() <-chan string {
	servMutex.RLock()
	defer servMutex.RUnlock()

	out := make(chan string)
	ids := make([]string, 0, len(servMap))
	for _, s := range servMap {
		ids = append(ids, s.id)
	}

	go func() {
		defer close(out)
		for _, sid := range ids {
			out <- sid
		}
	}()
	return out
}

func (s *Server) SetType(typ servType) os.Error {
	if s.styp != Unregistered {
		return os.NewError("Already registered")
	}

	s.styp = typ
	return nil
}

func (s *Server) SetPass(password, ts, prefix string) os.Error {
	if len(password) == 0 {
		return os.NewError("Zero-length password")
	}

	if ts != "6" {
		return os.NewError("TS " + ts + " is unsupported")
	}

	if !parser.ValidServerPrefix(prefix) {
		return os.NewError("SID " + prefix + " is invalid")
	}

	s.pass, s.sver = password, 6
	s.ts = time.Nanoseconds()
	return nil
}

func (s *Server) SetCapab(capab string) os.Error {
	if !strings.Contains(capab, "QS") {
		return os.NewError("QS CAPAB missing")
	}

	if !strings.Contains(capab, "ENCAP") {
		return os.NewError("ENCAP CAPAB missing")
	}

	s.capab = strings.Fields(capab)
	s.ts = time.Nanoseconds()
	return nil
}

func (s *Server) SetServer(serv, hops string) os.Error {
	if len(serv) == 0 {
		return os.NewError("Zero-length server name")
	}

	if hops != "1" {
		return os.NewError("Hops = " + hops + " is unsupported")
	}

	s.server, s.hops = serv, 1
	s.ts = time.Nanoseconds()
	return nil
}
