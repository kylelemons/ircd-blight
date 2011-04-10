package ircd

import (
	"testing"
)

func eq(t *testing.T, i int, desc string, exp,act interface{}) {
	if exp != act {
		t.Errorf("#%d: %s: Expected %#v, got %#v", i, desc, exp, act)
	}
}

func ne(t *testing.T, i int, desc string, exp,act interface{}) {
	if exp == act {
		t.Errorf("#%d: %s: Expected NOT %#v, got %#v", i, desc, exp, act)
	}
}

