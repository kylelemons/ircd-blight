package ircd

type RPC interface {}

type Return chan bool
func NewReturn() Return { return make(chan bool) }

type LinkFunc func(string,Link) bool
type Noop struct {Return}
type NewLink struct {Id string; Link; Return}
type EditLink struct {Id string; LinkFunc; Return}
type EachLink struct {LinkFunc; Return}

type DataStore struct {
	*LinkStore
	Control chan RPC
}


