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
					rpc.Success <- false
					continue
				}
				ls.locks[rpc.Id] = new(sync.Mutex)
				ls.links[rpc.Id] = rpc.Link
				rpc.Success <- true
			case Noop:
				rpc.Success <- true
			default:
				fmt.Fprintf(os.Stderr, "Unknown LinkStore RPC: %v\n", rpci)
		}
	}
}
