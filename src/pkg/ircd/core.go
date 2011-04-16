package ircd

import (
	ds "kevlar/ircd/datastore"
)

type Core struct {
	Data *ds.DataStore
}

func NewCore() *Core {
	core := &Core{
		Data: ds.NewDataStore(),
	}
	return core
}

func (c *Core) Set(module, key string, val ds.Value) {
	r := ds.NewReturn()
	set := ds.Set{
		Module: module,
		Key: key,
		Value: val,
		Return: r,
	}
	c.Data.Control <- set
	<-r
}

func (c *Core) Get(module, key string, def ds.Value) ds.Value {
	r := ds.NewReturn()
	get := &ds.Get{
		Module: module,
		Key: key,
		Return: r,
		Value: def,
	}
	c.Data.Control <- get
	<-r
	return get.Value
}

func (c *Core) Start() {
	//TODO(kevlar)
}
