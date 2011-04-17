package datastore

import (
	"fmt"
	"os"
	"sync"
)

type Link interface{}

type LinkFunc func(string, Link) bool
type NewLink struct {
	Id string
	Link
	Return
}
type EditLink struct {
	Id string
	LinkFunc
	Return
}
type EachLink struct {
	LinkFunc
	Return
}

type LinkStore struct {
	locks   map[string]*sync.Mutex
	links   map[string]Link
	Control chan RPC
}

func newLinkStore() *LinkStore {
	ls := new(LinkStore)
	ls.locks = make(map[string]*sync.Mutex)
	ls.links = make(map[string]Link)
	ls.Control = make(chan RPC)
	go ls.ControlLoop()
	return ls
}

func (ls *LinkStore) ControlLoop() {
	for rpci := range ls.Control {
		switch rpc := rpci.(type) {
		case NewLink:
			if _, ok := ls.links[rpc.Id]; ok {
				rpc.Return <- false
				continue
			}
			ls.locks[rpc.Id] = new(sync.Mutex)
			ls.links[rpc.Id] = rpc.Link
			rpc.Return <- true
		case EditLink:
			go func() {
				link, ok := ls.links[rpc.Id]
				if !ok {
					rpc.Return <- false
					return
				}
				ls.locks[rpc.Id].Lock()
				defer ls.locks[rpc.Id].Unlock()
				rpc.Return <- rpc.LinkFunc(rpc.Id, link)
			}()
		case EachLink:
			go func() {
				success := true
				for id, link := range ls.links {
					ls.locks[id].Lock()
					if !rpc.LinkFunc(id, link) {
						success = false
					}
					ls.locks[id].Unlock()
				}
				rpc.Return <- success
			}()
		case Noop:
			rpc.Return <- true
		default:
			fmt.Fprintf(os.Stderr, "Unknown LinkStore RPC: %v\n", rpci)
		}
	}
}
