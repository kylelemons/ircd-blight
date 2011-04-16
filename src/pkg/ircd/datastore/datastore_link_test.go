package datastore

import (
	"testing"
)

func TestLinkStoreNewLink(t *testing.T) {
	ls := NewLinkStore()
	success := NewReturn()
	ls.control <- NewLink{"SSSAAAAAA", 42, success}
	if !<-success {
		t.Errorf("NewLink should return true")
	}

	if lock,ok := ls.locks["SSSAAAAAA"]; !ok {
		t.Errorf("Lock should exist for SSSAAAAAA after NewLink")
	} else if lock == nil {
		t.Errorf("Lock should not be nil for SSSAAAAAA after NewLink")
	}

	if link,ok := ls.links["SSSAAAAAA"]; !ok {
		t.Errorf("Link should exist for SSSAAAAAA after NewLink")
	} else if link == nil {
		t.Errorf("Link should not be nil for SSSAAAAAA after NewLink")
	}

	close(ls.control)
}

func TestLinkStoreEditLink(t *testing.T) {
	ls := NewLinkStore()
	success := NewReturn()
	ls.control <- NewLink{"SSSAAAAAA", 42, success}
	<-success

	chk := make(map[int]bool)
	ls.control <- EditLink{"SSSAAAAAA", func(id string, l Link) bool {
		if 42 != l {
			t.Errorf("Link should be %v, got %v", 42, l)
		}
		if "SSSAAAAAA" != id {
			t.Errorf("ID should be %v, got %v", "SSSAAAAAA", id)
		}
		chk[l.(int)] = true
		return true
	}, success}
	<-success
	if val,ok := chk[42]; !ok || !val {
		t.Errorf("Call of func(42) expected, none recorded")
	}

	close(ls.control)
}

func TestLinkStoreEachLink(t *testing.T) {
	ls := NewLinkStore()
	success := NewReturn()
	ls.control <- NewLink{"SSSAAAAAA", 42, success}
	<-success
	ls.control <- NewLink{"SSSAAAAAB", 43, success}
	<-success

	chklink := make(map[int]bool)
	chkid := make(map[string]bool)
	ls.control <- EachLink{func(id string, l Link) bool {
		chklink[l.(int)] = true
		chkid[id] = true
		return true
	}, success}
	<-success
	if val,ok := chklink[42]; !ok || !val {
		t.Errorf("Call of func(42) expected, none recorded")
	}
	if val,ok := chklink[43]; !ok || !val {
		t.Errorf("Call of func(43) expected, none recorded")
	}
	if val,ok := chkid["SSSAAAAAA"]; !ok || !val {
		t.Errorf("Call of func(SSSAAAAAA) expected, none recorded")
	}
	if val,ok := chkid["SSSAAAAAB"]; !ok || !val {
		t.Errorf("Call of func(SSSAAAAAB) expected, none recorded")
	}

	close(ls.control)
}

func BenchmarkLinkStoreControlLoop(b *testing.B) {
	ls := NewLinkStore()
	success := make(chan bool)
	for i := 0; i < b.N; i++ {
		ls.control <- Noop{success}
		<-success
	}
	close(ls.control)
}
