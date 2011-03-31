package ircd

import (
	"testing"
)

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
	for i,test := range parse_message_tests {
		m := ParseMessage([]byte(test.raw));
		eq(t,i, "prefix", test.prefix, m.Prefix)
		eq(t,i, "command", test.command, m.Command)
		eq(t,i, "arglen", len(test.args), len(m.Args))
		for j := 0; j < len(test.args) && j < len(m.Args); j++ {
			eq(t,i, "arg", test.args[j], m.Args[j])
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
