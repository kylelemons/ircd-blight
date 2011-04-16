package parser

import (
	"testing"
)

import u "kevlar/ircd/util"

var parse_message_tests = []struct{
	raw string
	prefix, command string
	args []string
}{
{":server.kevlar.net NOTICE user :*** This is a test",
 "server.kevlar.net", "NOTICE", []string{"user", "*** This is a test"}},
{":A B C", "A", "B", []string{"C"}},
{"B C", "", "B", []string{"C"}},
}

func TestParseMesage(t *testing.T) {
	cmp := u.Test(t)
	for _,test := range parse_message_tests {
		m := ParseMessage([]byte(test.raw));
		cmp.EQ("prefix", test.prefix, m.Prefix)
		cmp.EQ("command", test.command, m.Command)
		cmp.EQ("arglen", len(test.args), len(m.Args))
		for j := 0; j < len(test.args) && j < len(m.Args); j++ {
			cmp.EQ("arg", test.args[j], m.Args[j])
		}
	}
}

var build_message_tests = []struct{
	expected string
	prefix, command string
	args []string
}{
{":server.kevlar.net NOTICE user :*** This is a test",
 "server.kevlar.net", "NOTICE", []string{"user", "*** This is a test"}},
{":A B C", "A", "B", []string{"C"}},
{"B C", "", "B", []string{"C"}},
{":A B C D", "A", "B", []string{"C", "D"}},
}
func TestBuildMessage(t *testing.T) {
	cmp := u.Test(t)
	for _,test := range build_message_tests {
		m := &Message{
			Prefix: test.prefix,
			Command: test.command,
			Args: test.args,
		}
		bytes := m.Bytes()
		str := m.String()
		cmp.EQ("bytes", test.expected, string(bytes))
		cmp.EQ("string", test.expected, str)
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
