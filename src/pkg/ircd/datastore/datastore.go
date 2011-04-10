package datastore

import (
	"fmt"
	"os"
)

// General types
type RPC interface {}
type Value interface {}

type Return chan bool
func NewReturn() Return { return make(chan bool) }

// Universal RPCs
type Noop struct {Return}

// DataStore RPCs
type Get struct {Module, Key string; Return; Value} // Use &Get{}
type Set struct {Module, Key string; Value; Return}
type Unset struct {Module, Key string; Return}

type DataStore struct {
	*LinkStore
	Control chan RPC
	config map[string]map[string]Value
}

func NewDataStore() *DataStore {
	ds := new(DataStore)
	ds.LinkStore = NewLinkStore()
	ds.Control = make(chan RPC)
	ds.config = make(map[string]map[string]Value)
	go ds.ControlLoop()
	return ds
}

func (ds *DataStore) ControlLoop() {
	for rpci := range ds.Control {
		switch rpc := rpci.(type) {
			case *Get:
				// Value will only be modified if the Module and Key both exist.
				// Thus, setting it in the RPC is equivalent to setting a default.
				if mod,ok := ds.config[rpc.Module]; ok {
					if val,ok := mod[rpc.Key]; ok {
						rpc.Value = val
						rpc.Return <- true
					} else {
						rpc.Return <- false
					}
				} else {
					rpc.Return <- false
				}

			case Set:
				if _,ok := ds.config[rpc.Module]; !ok {
					ds.config[rpc.Module] = make(map[string]Value)
				}
				ds.config[rpc.Module][rpc.Key] = rpc.Value
				rpc.Return <- true

			case Unset:
				if mod,ok := ds.config[rpc.Module]; ok {
					if _,ok := mod[rpc.Key]; ok {
						mod[rpc.Key] = nil,false
						if len(mod) == 0 {
							ds.config[rpc.Module] = nil,false
						}
						rpc.Return <- true
					} else {
						rpc.Return <- false
					}
				} else {
					rpc.Return <- false
				}

			case Noop:
				rpc.Return <- true

			default:
				fmt.Fprintf(os.Stderr, "Unknown DataStore RPC: %v\n", rpci)
		}
	}
}
