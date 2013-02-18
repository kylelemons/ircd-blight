// Package grid implements a data structure which represents mutual membership.
//
// This package should generally not be interacted with directly; the IRC server
// etc should wrap its functionality with a cleaner interface.
package grid

import (
	"log"
	"sync"
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
// In the case of IRC: a user is a member of channels, but channels also need
// to store the users who are on the channel.  This structure could be used
// largely unchanged if the IRC-specific data is removed.
type Grid struct {
	Edges [2]Edge
}

// Insert inserts the specified pair into the data structure and returns
// whether an actual insertion was performed.
//
// If either of the keys do not exist along the appropriate edge, they will
// be created (as with Touch) but with nil data.
//
// This function is goroutine safe.
func (g *Grid) Insert(keys [2]string, data interface{}) bool {
	// Get the List on each Edge
	lists := [2]*List{g.Edges[0].Touch(keys[0], nil), g.Edges[1].Touch(keys[1], nil)}

	// Find the insertion point along each axis
	ptrs, exacts := [2]**Membership{}, [2]bool{}
	for i := range keys {
		lists[i].lock.Lock()
		defer lists[i].lock.Unlock()

		ptrs[i], exacts[i] = lists[i].find(i, keys[1-i])
		if exacts[i] {
			return false
		}
	}

	// Construct and insert new membership node
	mem := Membership{
		Edges: lists,
		next:  [2]*Membership{*ptrs[0], *ptrs[1]},
		Data:  data,
	}
	for i := range ptrs {
		*ptrs[i] = &mem
	}

	return true
}

// Get gets the membership between the two keys if it exists and also
// returns true if the memberhsip was found.
func (g *Grid) Get(keys [2]string) (*Membership, bool) {
	// Find the two edges
	lists, ok := [2]*List{}, false
	for i, key := range keys {
		if lists[i], ok = g.Edges[i].Get(key); !ok {
			return nil, false
		}
	}

	// Find the membership links
	mems := [2]**Membership{}
	for i, lst := range lists {
		lst.lock.Lock()
		defer lst.lock.Unlock()

		mems[i], ok = lst.find(i, keys[1-i])
		if ok {
			return nil, false
		}
	}

	// Sanity check before returning externally
	if *mems[0] != *mems[1] {
		log.Panicf("grid[%q][%q] returned mismatched memberships", keys[0], keys[1])
	}

	// TODO(kevlar): Someday we can optimize this by not looking in both
	// lists; I'm not sure how to choose which one, though.

	return *mems[0], true
}

// Delete removes an association between the given keys and returns whether or
// not a deletion was performed.
//
// This function is goroutine-safe.
func (g *Grid) Delete(keys [2]string) bool {
	// Find the two edges
	lists, ok := [2]*List{}, false
	for i, key := range keys {
		if lists[i], ok = g.Edges[i].Get(key); !ok {
			return false
		}
	}

	// Find the membership links to remove
	mems := [2]**Membership{}
	for i, lst := range lists {
		lst.lock.Lock()
		defer lst.lock.Unlock()

		mems[i], ok = lst.find(i, keys[1-i])
		if !ok {
			return false
		}
	}

	// Skip the membership and let the GC deal with the link.
	for i, mem := range mems {
		*mem, (*mem).next[i] = (*mem).next[i], nil
	}

	return true
}

// DeleteAll removes the given key along the given edge and all memberships
// it has and returns:
//   - the key of any list from which a membership was deleted
//   - the keys of other members in the above lists
//   - a boolean indicating whether the key was found
//
// This function is goroutine-safe, but very expensive.
func (g *Grid) DeleteAll(edge int, key string) ([]string, [][]string, bool) {
	// Always acquire the locks in order (to prevent deadlock)
	for _, e := range g.Edges {
		e := e
		e.lock.Lock()
		defer e.lock.Unlock()
	}

	// Get the primary and secondary axes
	pri, pidx := g.Edges[edge], edge
	sec, sidx := g.Edges[1-edge], 1-edge
	_ = sec

	// Find the list we're deleting and delete it or return
	lst, ok := pri.lists[key]
	delete(pri.lists, key)
	if !ok {
		return nil, nil, false
	}

	var keys []string
	var affected [][]string

	// Perform the deletions
	lst.lock.Lock()
	defer lst.lock.Unlock()
	for m := lst.members; m != nil; m = m.next[pidx] {
		other := m.Edges[sidx]

		var notify []string
		for n := other.members; n != nil; n = n.next[sidx] {
			notify = append(notify, n.Edges[pidx].Name)
		}

		keys = append(keys, other.Name)
		affected = append(affected, notify)

		mem, ok := other.find(sidx, key)
		if !ok {
			log.Printf("failed to find %q along %q axis", key, other.Name)
			continue
		}
		*mem, (*mem).next[sidx] = (*mem).next[sidx], nil
	}

	return keys, affected, true
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

	Data interface{}
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
func (e *Edge) Touch(name string, data interface{}) *List {
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
		Data: data,
	}
	e.lists[name] = lst
	return lst
}

// A List acts as the head node of one axis of membership and contains its
// name.  Mutations of the linked list should be interlocked.
type List struct {
	Name string
	lock sync.RWMutex

	// If this List is on Grid.Edges[idx], follow the list by iterating over
	// member.next[idx].
	members *Membership

	Data interface{}
}

// A Membership is a member of two linked lists along two axes.  Once added to
// the lists, it should only be mutated while both lists are locked.
type Membership struct {
	Edges [2]*List
	next  [2]*Membership

	Data interface{}
}

// find returns a pointer to the link into which the given name would be
// inserted along the given axis.  If name already exists in the List, the link
// is the link pointing to the memnership with the matching name.  This should
// be called only after the lock is acquired.
func (l *List) find(along int, name string) (**Membership, bool) {
	// pm will never be nil
	pm := &l.members

	// Find a link into which to insert a Membership
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
