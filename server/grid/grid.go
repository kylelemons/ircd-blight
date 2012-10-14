package grid

import (
	"sync"
	"log"
)

// TODO(kevlar): remove
var _ = log.Printf

var (
	// DefaultCapacity stores the initial capacity of each edge of the grid.
	DefaultCapacity = 100
)

// A Grid represents a mutual membership data structure.  Its zero value is
// ready to use.
//
// To use IRC as an example: a user is a member of channels, but channels also
// need to store the users who are on the channel.
type Grid struct {
	Edges [2]Edge
}

// Insert inserts the specified pair into the data structure and returns
// whether an actual insertion was performed.
//
// This function is goroutine safe.
func (g *Grid) Insert(keys [2]string) bool {
	// Get the List on each Edge
	lists := [2]*List{g.Edges[0].Touch(keys[0]), g.Edges[1].Touch(keys[1])}

	// Find the insertion point along each axis
	ptrs, exacts := [2]**membership{}, [2]bool{}
	for i := range keys {
		lists[i].lock.Lock()
		defer lists[i].lock.Unlock()

		ptrs[i], exacts[i] = lists[i].find(i, keys[1-i])
		if exacts[i] {
			return false
		}
	}

	// Construct and insert new membership node
	mem := membership{
		Edges: lists,
		next:  [2]*membership{*ptrs[0], *ptrs[1]},
	}
	for i := range ptrs {
		*ptrs[i] = &mem
	}

	return true
}

// Delete removes an association between the given keys and returns whether or
// not a deletion was performed.
func (g *Grid) Delete(keys [2]string) bool {
	// Find the two edges
	lists, oks := [2]*List{}, [2]bool{}
	for i, key := range keys {
		lists[i], oks[i] = g.Edges[i].Get(key)
		if !oks[i] {
			return false
		}
	}

	// Find the membership links to remove
	mems := [2]**membership{}
	for i, lst := range lists {
		lst.lock.Lock()
		defer lst.lock.Unlock()

		mems[i], oks[i] = lst.find(i, keys[1-i])
		if !oks[i] {
			return false
		}
	}

	// Skip the membership and let the GC deal with the link.
	for i, mem := range mems {
		*mem = (*mem).next[i]
	}

	return true
}

// dump returns a snapshot of the grid.
//
// This function is not goroutine safe and should only be used in unit tests.
func (g *Grid) dump() (out [2]map[string][]string) {
	for e, edg := range g.Edges {
		out[e] = make(map[string][]string)
		for name, lst := range edg.lists {
			for m := lst.members; m != nil; m = m.next[e] {
				out[e][name] = append(out[e][name], m.Edges[1-e].Name)
			}
		}
	}
	return out
}

// An Edge represents one dimension of a Grid.
type Edge struct {
	lock  sync.RWMutex
	lists map[string]*List
}

// Get returns a List from the Edge and a boolean indicating whether the
// edge was found.
//
// This function is goroutine safe.
func (e *Edge) Get(name string) (lst *List, ok bool) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	lst, ok = e.lists[name]
	return
}

// Touch returns a List, crating and initializing it if necessary.
//
// This function is goroutine safe.
func (e *Edge) Touch(name string) *List {
	e.lock.Lock()
	defer e.lock.Unlock()

	// Check if the List already exists
	lst, ok := e.lists[name]
	if ok {
		return lst
	}

	// Make sure lists exists
	if e.lists == nil {
		e.lists = make(map[string]*List, DefaultCapacity)
	}

	// Store and return a new one
	lst = &List{
		Name: name,
	}
	e.lists[name] = lst
	return lst
}

// A List acts as the head node of one axis of membership and contains its
// name.  Mutations of the linked list should be interlocked.
type List struct {
	Name    string
	lock    sync.RWMutex
	members *membership
}

// A membership is a member of two linked lists along two axes.  Once added to
// the lists, it should only be mutated while both lists are locked.
type membership struct {
	Edges [2]*List
	next  [2]*membership
}

// find returns a pointer to the link into which the given name would be
// inserted along the given axis.  If name already exists in the List, the link
// is the link pointing to the memnership with the matching name.  If the
// returned link is going to be modified, this should be called after the lock
// is acquired.
func (l *List) find(along int, name string) (**membership, bool) {
	// pm will never be nil
	pm := &l.members

	// Find a link into which to insert a membership
	for *pm != nil {
		if next := (*pm).Edges[1-along].Name; name == next {
			return pm, true
		} else if name < next {
			return pm, false
		}
		pm = &(*pm).next[along]
	}

	return pm, false
}
