package util

import (
	"testing"
)

type Comparator struct {
	t *testing.T
	i int
}

func Test(t *testing.T) *Comparator {
	return &Comparator{t,0}
}

func (c *Comparator) EQ(desc string, exp,act interface{}) {
	if exp != act {
		c.t.Errorf("#%d: %s: Expected %#v, but got %#v", c.i, desc, exp, act)
	}
	c.i++
}

func (c *Comparator) NE(desc string, exp,act interface{}) {
	if exp == act {
		c.t.Errorf("#%d: %s: Expected NOT %#v, but got %#v", c.i, desc, exp, act)
	}
	c.i++
}

