package server

import (
	"kevlar/ircd/parser"
	"kevlar/ircd/log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// servMap[SID] = &Server{...}
// - Stores the information for any server
// linkMap[server SID] = linked SID
// - Maps a logical server to the peer to which it's linked
var (
	servMap   = make(map[string]*Server)
	linkMap   = make(map[string]string)
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
	desc   string
	sver   int
	server string
	link   string
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
	linkMap[id] = id
	return s
}

// Link registers a new server linked behind link.
func Link(link, sid, name, hops, desc string) os.Error {
	servMutex.Lock()
	defer servMutex.Unlock()

	if _, ok := servMap[sid]; ok {
		return os.NewError("Server already linked: " + sid)
	}

	ihops, _ := strconv.Atoi(hops)

	s := &Server{
		mutex:  new(sync.RWMutex),
		id:     sid,
		server: name,
		desc:   desc,
		link:   link,
		hops:   ihops,
	}

	if downstream, ok := linkMap[link]; ok {
		link = downstream
	}

	log.Info.Printf("Server %s (%s) linked behind %s", name, sid, link)

	servMap[sid] = s
	linkMap[sid] = link
	return nil
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

// Iter iterates over all server links
func Iter() <-chan string {
	servMutex.RLock()
	defer servMutex.RUnlock()

	out := make(chan string)
	ids := make([]string, 0, len(servMap))
	for sid := range linkMap {
		ids = append(ids, sid)
	}

	go func() {
		defer close(out)
		for _, sid := range ids {
			out <- sid
		}
	}()
	return out
}

// IterFor iterates over the link IDs for all of the ID in the given list.
// The list may contain SIDs, UIDs, or both.  If the skipLink is given,
// any servers behind that link will be skipped.
func IterFor(ids []string, skipLink string) <-chan string {
	servMutex.RLock()
	defer servMutex.RUnlock()

	if actual, ok := linkMap[skipLink]; ok {
		skipLink = actual
	}

	out := make(chan string)
	links := make(map[string]bool)
	for _, id := range ids {
		sid := id[:3]
		if link, ok := linkMap[sid]; ok {
			links[link] = true
		}
	}
	links[skipLink] = false, false

	go func() {
		defer close(out)
		for sid := range links {
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
