package conn

import (
	"net"
	"os"
	"testing"
	"kevlar/ircd/parser"
)

type MockConn struct {
	data      []string
	lastwrite []byte
}

func (mc *MockConn) LocalAddr() net.Addr                 { return nil }
func (mc *MockConn) RemoteAddr() net.Addr                { return nil }
func (mc *MockConn) SetTimeout(nsec int64) os.Error      { return nil }
func (mc *MockConn) SetReadTimeout(nsec int64) os.Error  { return nil }
func (mc *MockConn) SetWriteTimeout(nsec int64) os.Error { return nil }
func (mc *MockConn) Close() os.Error {
	mc.data = nil
	return nil
}
func (mc *MockConn) Write(b []byte) (n int, err os.Error) {
	mc.lastwrite = b
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
	mc := new(MockConn)
	mc.Add(":source command arg :longarg\r\n")
	mc.Add(":SOURCE COMMAND ARG :LONGARG\n")
	conn := NewConn(mc)
	messages := make(chan *parser.Message)
	conn.Subscribe(messages)
	if conn == nil {
		t.Fatalf("NewConn should not return a nil connection")
	}
	// Test with \r\n
	message := <-messages
	if "source" != message.Prefix {
		t.Errorf("Message prefix %q expected, got %q", "source", message.Prefix)
	}
	if "COMMAND" != message.Command {
		t.Errorf("Message command %q expected, got %q", "command", message.Command)
	}
	if 2 != len(message.Args) || "arg" != message.Args[0] || "longarg" != message.Args[1] {
		t.Errorf("Message args %v expected, got %v", []string{"arg", "longarg"}, message.Args)
	}
	message = <-messages
	if "SOURCE" != message.Prefix {
		t.Errorf("Message prefix %q expected, got %q", "SOURCE", message.Prefix)
	}
	if "COMMAND" != message.Command {
		t.Errorf("Message command %q expected, got %q", "COMMAND", message.Command)
	}
	if 2 != len(message.Args) || "ARG" != message.Args[0] || "LONGARG" != message.Args[1] {
		t.Errorf("Message args %v expected, got %v", []string{"ARG", "LONGARG"}, message.Args)
	}
}

func TestWriteMessage(t *testing.T) {
	msg := &parser.Message{
		Prefix:  "server",
		Command: "COMMAND",
		Args:    []string{"arg1", "arg2", "arg3 arg3"},
	}
	mc := new(MockConn)
	conn := NewConn(mc)
	conn.WriteMessage(msg)
	if ":server COMMAND arg1 arg2 :arg3 arg3" != string(mc.lastwrite) {
		t.Errorf("Expected write of %q, got %q", ":server COMMAND arg1 arg2 :arg3 arg3",
			string(mc.lastwrite))
	}
}

func BenchmarkWriteMessage(b *testing.B) {
	msg := &parser.Message{
		Prefix:  "server",
		Command: "COMMAND",
		Args:    []string{"arg1", "arg2", "arg3 arg3"},
	}
	mc := new(MockConn)
	conn := NewConn(mc)
	for i := 0; i < b.N; i++ {
		conn.WriteMessage(msg)
	}
}
