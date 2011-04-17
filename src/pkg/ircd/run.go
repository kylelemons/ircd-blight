package ircd

import (
	//"net"
	"os"
	"fmt"
	//"bufio"
)

import (
//"kevlar/ircd/parser"
)

func errchk(err os.Error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func Run() {
	c := NewCore()
	c.Start()
	select {
	}
}
