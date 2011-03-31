package ircd

import (
	"testing"
)

func TestLinkStoreNewLink(t *testing.T) {
	ls := NewLinkStore()
	success := make(chan bool)
	ls.control <- NewLink{"SSSAAAAAA", 42, success}
	eq(t,0, "success", true, <-success)

	lock,ok := ls.locks["SSSAAAAAA"]
	eq(t,1, "lock exists", true, ok)
	ne(t,2, "lock", nil, lock)

	link,ok := ls.links["SSSAAAAAA"]
	eq(t,3, "link exists", true, ok)
	eq(t,4, "link", 42, link)
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
