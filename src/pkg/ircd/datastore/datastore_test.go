package datastore

import (
	"testing"
)

import u "kevlar/ircd/util"

func TestNewDataStore(t *testing.T) {
	ds := NewDataStore()
	success := NewReturn()
	ds.Control <- Noop{success}
	u.EQ(t,0, "success", true, <-success)
	close(ds.Control)
}

func TestGetSetUnset(t *testing.T) {
	ds := NewDataStore()
	success := NewReturn()

	// Test get without a Module
	rpc := &Get{"General", "Answer", success, "default"}
	ds.Control <- rpc
	u.EQ(t,0, "No module - Return", false, <-success)
	u.EQ(t,1, "No module - Value", "default", rpc.Value)

	// Test set
	ds.Control <- Set{"General", "Answer", 42, success}
	u.EQ(t,2, "Set Answer", true, <-success)

	// Test get with out a value
	rpc = &Get{"General", "Noexist", success, "oops"}
	ds.Control <- rpc
	u.EQ(t,3, "No module - Return", false, <-success)
	u.EQ(t,4, "No module - Value", "oops", rpc.Value)

	// Test get with value
	rpc = &Get{"General", "Answer", success, "default"}
	ds.Control <- rpc
	u.EQ(t,5, "Get Answer - Return", true, <-success)
	t.Logf("Config: %v", ds.config)
	u.EQ(t,6, "Get Answer - Value", 42, rpc.Value)

	// Test unset
	ds.Control <- Unset{"General", "Noexist", success}
	u.EQ(t,7, "Unsea Noexist - Return", false, <-success)
	u.EQ(t,8, "Unset Noexist - Length", 1, len(ds.config))

	// Test unset
	ds.Control <- Unset{"General", "Answer", success}
	u.EQ(t,9, "Unset Answer - Return", true, <-success)
	u.EQ(t,10, "Unset Answer - Length", 0, len(ds.config))

	// Test unset
	ds.Control <- Unset{"General", "Noexist", success}
	u.EQ(t,11, "Unset Nomodule - Return", false, <-success)

	close(ds.Control)
}

func BenchmarkDataStoreControlLoop(b *testing.B) {
	ds := NewDataStore()
	success := make(chan bool)
	for i := 0; i < b.N; i++ {
		ds.Control <- Noop{success}
		<-success
	}
	close(ds.Control)
}
