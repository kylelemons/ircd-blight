package ircd

import (
	"fmt"
	"os"
	"sync"
)

type Link interface {}

type LinkStore struct {
	locks map[string]*sync.Mutex
	links map[string]Link
	control chan RPC
}

func NewLinkStore() *LinkStore {
	ls := new(LinkStore)
	ls.locks = make(map[string]*sync.Mutex)
	ls.links = make(map[string]Link)
	ls.control = make(chan RPC)
	go ls.ControlLoop()
	return ls
}

func (ls *LinkStore) ControlLoop() {
	for rpci := range ls.control {
		switch rpc := rpci.(type) {
			case NewLink:
				if _,ok := ls.links[rpc.Id]; ok {
					rpc.Return <- false
					continue
				}
				ls.locks[rpc.Id] = new(sync.Mutex)
				ls.links[rpc.Id] = rpc.Link
				rpc.Return <- true
			case EditLink:
				go func() {
					link,ok := ls.links[rpc.Id]
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
					for id,link := range ls.links {
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
