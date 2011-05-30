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

// For the purposes of this file, think of the current server as the root of a
// tree where each node has its links hanging below it.  Downstream refers to
// any server directly connected below the node, upstream refers to the single
// server above the node.

// servMap[SID] = &Server{...}
// - Stores the information for any server
// upstream[remote SID] = upstream SID
// - Maps a logical server to the peer to which it's linked
// - Locally linked servers aren't in this map.
// downstream[sid1][sid2] = bool
// - if downstream[sid1][sid2] == true, sid2 is directly downstream of sid1
var (
	servMap    = make(map[string]*Server)
	upstream   = make(map[string]string)
	downstream = make(map[string]map[string]bool)
	servMutex  = new(sync.RWMutex)
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

// Get retrieves and/or creates a server.  If a server is created,
// it is directly linked to this one.
func Get(id string, create bool) *Server {
	servMutex.Lock()
	defer servMutex.Unlock()

	if s, ok := servMap[id]; ok {
		return s
	}

	if !create {
		return nil
	}

	s := &Server{
		mutex: new(sync.RWMutex),
		id:    id,
	}
	servMap[id] = s
	downstream[id] = make(map[string]bool)
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

	servMap[sid] = s
	upstream[sid] = link
	downstream[link][sid] = true
	downstream[sid] = make(map[string]bool)

	up := link
	chain := 1
	for len(up) > 0 {
		up, _ = upstream[up]
		chain++
	}

	log.Info.Printf("Server %s (%s) linked behind %s", name, sid, link)
	log.Info.Printf("%s: %d hops found, %d hops reported", sid, chain, ihops)

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
	links := make([]string, 0, len(downstream))
	for link := range servMap {
		if _, skip := upstream[link]; skip {
			continue
		}
		links = append(links, link)
	}

	go func() {
		defer close(out)
		for _, sid := range links {
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

	for {
		if actual, ok := upstream[skipLink]; ok {
			skipLink = actual
		} else {
			break
		}
	}

	out := make(chan string)
	links := []string{}
	cache := make(map[string]bool)

nextId:
	for _, id := range ids {
		sid := id[:3]
		for {
			if cache[sid] {
				continue nextId
			}
			cache[sid] = true
			if actual, ok := upstream[sid]; ok {
				sid = actual
			} else {
				break
			}
		}
		if sid == skipLink {
			continue nextId
		}
		links = append(links, sid)
	}

	go func() {
		defer close(out)
		for _, sid := range links {
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

// IsLocal returns true if the SID is locally linked
func IsLocal(sid string) bool {
	if _, remote := upstream[sid]; remote {
		return false
	}
	if _, exists := downstream[sid]; !exists {
		return false
	}
	return true
}

// Return the SIDs of all servers behind the given link, starting with the
// server itself.  If the server is unknown, the returned list is empty.
func LinkedTo(link string) []string {
	servMutex.Lock()
	defer servMutex.Unlock()

	return linkedTo(link)
}

// Make sure the server mutex is (r)locked before calling this.
func linkedTo(link string) []string {
	if _, ok := servMap[link]; !ok {
		log.Warn.Printf("Mapping nonexistent link %s", link)
		return nil
	}

	sids := []string{link}
	for downstream := range downstream[link] {
		sids = append(sids, linkedTo(downstream)...)
	}

	return sids
}

// Unlink deletes the given server and all servers behind it.  It returns the list
// of SIDs that were split.
func Unlink(split string) (sids []string) {
	servMutex.Lock()
	defer servMutex.Unlock()

	sids = linkedTo(split)

	for _, sid := range sids {
		log.Info.Printf("Split %s: Unlinking %s", split, sid)

		// Delete the server entry
		servMap[sid] = nil, false

		// Unlink from the upstream server's downstream list
		if up, ok := upstream[sid]; ok {
			// But only if the upstream server is still around
			if _, ok := downstream[up]; ok {
				downstream[up][sid] = false, false
			}
		}

		// Remove the server's downstream list
		downstream[sid] = nil, false
	}

	return
}
