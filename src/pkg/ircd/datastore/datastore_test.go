package datastore

import (
	"testing"
)

func TestNewDataStore(t *testing.T) {
	ds := NewDataStore()
	success := NewReturn()
	ds.Control <- Noop{success}
	if !<-success {
		t.Errorf("No-op should return true")
	}
	close(ds.Control)
}

func TestGetSetUnset(t *testing.T) {
	ds := NewDataStore()
	success := NewReturn()

	// Test get without a Module
	rpc := &Get{"General", "Answer", success, "default"}
	ds.Control <- rpc
	if <-success {
		t.Errorf("Get on nonexistent module should return false")
	} else if "default" != rpc.Value {
		t.Errorf("Nonexistent module get: Expected %q, got %q", "default", rpc.Value)
	}

	// Test set
	ds.Control <- Set{"General", "Answer", 42, success}
	if !<-success {
		t.Errorf("Set on nonexistent module should return true")
	}
	t.Logf("Config after Set: %v", ds.config)

	// Test get with out a value
	rpc = &Get{"General", "Noexist", success, "oops"}
	ds.Control <- rpc
	if <-success {
		t.Errorf("Get on nonexistent value should return false")
	} else if "oops" != rpc.Value {
		t.Errorf("Nonexistent value get: Expected %q, got %q", "oops", rpc.Value)
	}

	// Test get with value
	rpc = &Get{"General", "Answer", success, "default"}
	ds.Control <- rpc
	if !<-success {
		t.Errorf("Get on existing value should return true")
	} else if 42 != rpc.Value {
		t.Errorf("Existing value get: Expected %v, got %v", 42, rpc.Value)
	}

	// Test unset
	ds.Control <- Unset{"General", "Noexist", success}
	if <-success {
		t.Errorf("Unset on nonexistent value should return false")
	}
	if 1 != len(ds.config) {
		t.Errorf("Config should have exactly %d values, has %d", 1, len(ds.config))
	}

	// Test unset
	ds.Control <- Unset{"General", "Answer", success}
	if !<-success {
		t.Errorf("Unset on existing value should return true")
	}
	if 0 != len(ds.config) {
		t.Errorf("Config should have exactly %d values, has %d", 0, len(ds.config))
	}

	// Test unset
	ds.Control <- Unset{"General", "Noexist", success}
	if <-success {
		t.Errorf("Unset on nonexisting value should return false")
	}

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
