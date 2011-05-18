package core

import (
	"testing"
)

func TestGetSetUnset(t *testing.T) {
	c := NewCore()

	c.Set("mod", "key", "val")
	if got, want := len(c.params), 1; got != want {
		t.Fatalf("len(params) = %d, want %d", got, want)
	}
	if _, got := c.params["mod"]; !got {
		t.Fatalf("params[%q] undefined", "mod")
	}
	if got, want := c.params["mod"]["key"], "val"; got != want {
		t.Fatalf("params[%q][%q] = %q, want %q", "mod", "key", got, want)
	}

	if got, want := c.Get("mod", "key", "default"), "val"; got != want {
		t.Fatalf("get(%q) = %q, want %q", "key", got, want)
	}
	if got, want := c.Get("mod", "bad", "default"), "default"; got != want {
		t.Errorf("got(%q) = %q, want %q", "bad", got, want)
	}

	c.Set("mod", "key2", "val2")
	if got, want := len(c.params), 1; got != want {
		t.Fatalf("after get, len(params) = %d, want %d", got, want)
	}
	if got, want := len(c.params["mod"]), 2; got != want {
		t.Fatalf("after get, len(params[%q]) = %d, want %d", "mod", got, want)
	}

	c.Unset("mod", "key2")
	if got, want := len(c.params), 1; got != want {
		t.Fatalf("after get, len(params) = %d, want %d", got, want)
	}
	if got, want := len(c.params["mod"]), 1; got != want {
		t.Fatalf("after get, len(params[%q]) = %d, want %d", "mod", got, want)
	}

	c.Unset("mod", "key")
	if got, want := len(c.params), 0; got != want {
		t.Fatalf("after get, len(params) = %d, want %d", got, want)
	}
}
