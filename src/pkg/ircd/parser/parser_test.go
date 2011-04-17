package parser

import (
	"testing"
)

var parse_message_tests = []struct {
	raw             string
	prefix, command string
	args            []string
}{
	{":server.kevlar.net NOTICE user :*** This is a test",
		"server.kevlar.net", "NOTICE", []string{"user", "*** This is a test"}},
	{":A B C", "A", "B", []string{"C"}},
	{"B C", "", "B", []string{"C"}},
}

func TestParseMesage(t *testing.T) {
	for i, test := range parse_message_tests {
		m := ParseMessage([]byte(test.raw))
		if test.prefix != m.Prefix {
			t.Errorf("#d: Expected prefix %q, got %q", i, test.prefix, m.Prefix)
		}
		if test.command != m.Command {
			t.Errorf("#d: Expected command %q, got %q", i, test.command, m.Command)
		}
		if len(test.args) != len(m.Args) {
			t.Errorf("#d: Expected args %v, got %v", i, test.args, m.Args)
		} else {
			for j := 0; j < len(test.args) && j < len(m.Args); j++ {
				if test.args[j] != m.Args[j] {
					t.Errorf("#d: Expected arg[%d] %q, got %q", i, test.args[j], m.Args[j])
				}
			}
		}
	}
}

var build_message_tests = []struct {
	expected        string
	prefix, command string
	args            []string
}{
	{":server.kevlar.net NOTICE user :*** This is a test",
		"server.kevlar.net", "NOTICE", []string{"user", "*** This is a test"}},
	{":A B C", "A", "B", []string{"C"}},
	{"B C", "", "B", []string{"C"}},
	{":A B C D", "A", "B", []string{"C", "D"}},
}

func TestBuildMessage(t *testing.T) {
	for i, test := range build_message_tests {
		m := &Message{
			Prefix:  test.prefix,
			Command: test.command,
			Args:    test.args,
		}
		bytes := m.Bytes()
		str := m.String()
		if test.expected != str {
			t.Errorf("Expected string representation %q, got %q", i, test.expected, str)
		}
		if string(bytes) != str {
			t.Errorf("Expected identical string and byte representation, got %q and %q", i,
				bytes, str)
		}
	}
}

var parse_message_bench = []byte(":server.kevlar.net NOTICE user :*** This is a test")

func BenchmarkParseMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseMessage(parse_message_bench)
	}
}

var build_message_bench = ParseMessage(parse_message_bench)

func BenchmarkBuildMessage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		build_message_bench.String()
	}
}

func BenchmarkBuildMessageBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		build_message_bench.Bytes()
	}
}
