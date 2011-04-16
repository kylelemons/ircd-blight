package conn

import (
	"fmt"
	"net"
	"os"
	"testing"
	"kevlar/ircd/parser"
	u "kevlar/ircd/util"
)

type MockConn struct {
	data []string
	lastwrite []byte
}

func (mc *MockConn) LocalAddr() net.Addr { return nil }
func (mc *MockConn) RemoteAddr() net.Addr { return nil }
func (mc *MockConn) SetTimeout(nsec int64) os.Error { return nil }
func (mc *MockConn) SetReadTimeout(nsec int64) os.Error { return nil }
func (mc *MockConn) SetWriteTimeout(nsec int64) os.Error { return nil }
func (mc *MockConn) Close() os.Error {
	mc.data = nil
	return nil
}
func (mc *MockConn) Write(b []byte) (n int, err os.Error) {
	mc.lastwrite = b
	fmt.Printf("Write: %v\n", b)
	return len(b), nil
}
func (mc *MockConn) Read(b []byte) (n int, err os.Error) {
	if len(mc.data) <= 0 {
		err = os.EOF
		return
	}
	next := mc.data[0]
	n = len(next)
	if n > len(b) {
		n = len(b)
	}
	copy(b, next[:n])
	if len(next) > n {
		mc.data[0] = next[n+1:]
	} else {
		mc.data = mc.data[1:]
	}
	return
}
func (mc *MockConn) Add(s string) {
	mc.data = append(mc.data, s)
}

func TestConn(t *testing.T) {
	cmp := u.Test(t)
	mc := new(MockConn)
	mc.Add(":source command arg :longarg\r\n")
	mc.Add(":source command arg :longarg\n")
	conn := NewConn(mc)
	cmp.NE("conn", nil, conn)
	// Test with \r\n
	message := conn.ReadMessage()
	cmp.EQ("prefix", "source", message.Prefix)
	cmp.EQ("command", "command", message.Command)
	cmp.EQ("args", []string{"arg", "longarg"}, message.Args)
	fmt.Println("Done")
	message = conn.ReadMessage()
	cmp.EQ("prefix", "source", message.Prefix)
	cmp.EQ("command", "command", message.Command)
	cmp.EQ("args", []string{"arg", "longarg"}, message.Args)
}

func TestWriteMessage(t *testing.T) {
	cmp := u.Test(t)
	msg := &parser.Message{
		Prefix: "server",
		Command: "COMMAND",
		Args: []string{"arg1", "arg2", "arg3 arg3"},
	}
	mc := new(MockConn)
	conn := NewConn(mc)
	conn.WriteMessage(msg)
	cmp.EQ("Write", ":server COMMAND arg1 arg2 :arg3 arg3", string(mc.lastwrite))
}

func BenchmarkWriteMessage(b *testing.B) {
	msg := &parser.Message{
		Prefix: "server",
		Command: "COMMAND",
		Args: []string{"arg1", "arg2", "arg3 arg3"},
	}
	mc := new(MockConn)
	conn := NewConn(mc)
	for i := 0; i < b.N; i++ {
		conn.WriteMessage(msg)
	}
}
