package ircd

type RPC interface {}

type Noop struct {Success chan bool}
type NewLink struct {Id string; Link; Success chan bool}
