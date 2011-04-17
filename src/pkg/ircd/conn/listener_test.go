package conn

import (
	"runtime"
	"testing"
	"time"
)

func TestNewListener(t *testing.T) {
	l := NewListener()
	if l.ports == nil {
		t.Errorf("Listener ports should not be nil")
	}
	if l.Incoming == nil {
		t.Errorf("Listener Incoming should not be nil")
	}
	l.Close()
}

func TestAddPort(t *testing.T) {
	l := NewListener()
	gcnt := runtime.Goroutines()
	l.AddPort(56561)
	if 1 != len(l.ports) {
		t.Errorf("Length of ports array should be 1, got %d", len(l.ports))
	}
	if runtime.Gosched(); gcnt >= runtime.Goroutines() {
		t.Errorf("Expected more than %d goroutines after AddPort, %d running", gcnt,
			runtime.Goroutines())
	}
	if listener, ok := l.ports[56561]; ok {
		if listener == nil {
			t.Errorf("Port listener should not be nil")
		}
	} else {
		t.Errorf("Listener should have entry for port 56561, got %v", l.ports)
	}
	gcnt = runtime.Goroutines()
	l.Close()
	if 0 != len(l.ports) {
		t.Errorf("After Close(), ports should have 0 entries, got %d", len(l.ports))
	}
	if runtime.Gosched(); gcnt <= runtime.Goroutines() {
		t.Errorf("Expected fewer than %d goroutines after Close(), %d running", gcnt,
			runtime.Goroutines())
	}
}

func TestClosePort(t *testing.T) {
	l := NewListener()
	l.AddPort(56561)
	gcnt := runtime.Goroutines()
	l.ClosePort(56561)
	// ClosePort is not synchronized, so give it some time (on mac, dialog pops up)
	for i := 0; i < 100 && 0 != len(l.ports); i++ {
		time.Sleep(1e6)
	}
	if runtime.Gosched(); 0 != len(l.ports) {
		t.Errorf("After ClosePort(), ports should have 0 entries, got %d", len(l.ports))
	}
	if runtime.Gosched(); gcnt <= runtime.Goroutines() {
		t.Errorf("Expected fewer than %d goroutines after ClosePort(), %d running", gcnt,
			runtime.Goroutines())
	}
	l.Close()
	if 0 != len(l.ports) {
		t.Errorf("After Close(), ports should have 0 entries, got %d", len(l.ports))
	}
	if runtime.Gosched(); gcnt <= runtime.Goroutines() {
		t.Errorf("Expected fewer than %d goroutines after Close(), %d running", gcnt,
			runtime.Goroutines())
	}
}
