package channel

import (
	"kevlar/ircd/parser"
	"kevlar/ircd/user"
	"rand"
	"os"
	"sync"
	"testing"
)

var testJoinPart = []struct {
	ID      string
	Command string
	Channel string
	Error   os.Error
	Chans   int
	Notify  []string
}{
	{
		ID:      "A",
		Command: parser.CMD_JOIN,
		Channel: "#chan",
		Error:   nil,
		Chans:   1,
		Notify:  []string{"A"},
	},
	{
		ID:      "B",
		Command: parser.CMD_PART,
		Channel: "#chan",
		Error:   parser.NewNumeric(parser.ERR_NOTONCHANNEL, ""),
		Chans:   1,
		Notify:  nil,
	},
	{
		ID:      "B",
		Command: parser.CMD_JOIN,
		Channel: "#chan",
		Error:   nil,
		Chans:   1,
		Notify:  []string{"A", "B"},
	},
	{
		ID:      "A",
		Command: parser.CMD_PART,
		Channel: "#chan",
		Error:   nil,
		Chans:   1,
		Notify:  []string{"A", "B"},
	},
	{
		ID:      "B",
		Command: parser.CMD_PART,
		Channel: "#chan",
		Error:   nil,
		Chans:   0,
		Notify:  []string{"B"},
	},
}

func TestJoinPartChannel(t *testing.T) {
	for idx, test := range testJoinPart {
		var err os.Error
		var notify []string
		channel, _ := Get(test.Channel, true)
		switch test.Command {
		case parser.CMD_JOIN:
			notify, err = channel.Join(test.ID)
		case parser.CMD_PART:
			notify, err = channel.Part(test.ID)
		}
		if got, want := err, test.Error; got != want {
			if got == nil || want == nil || got.String() != want.String() {
				t.Errorf("#%d: %s returned %v, want %v", idx, test.Command, got, want)
			}
		}
		if got, want := len(chanMap), test.Chans; got != want {
			t.Errorf("#%d: chans after %s = %d, want %d", idx, test.Command, got, want)
		}
		if got, want := notify, test.Notify; len(got) != len(want) {
			t.Errorf("#%d: %s notify = %v, want %v", idx, test.Command, got, want)
		} else {
			for i := range notify {
				if got, want := notify[i], test.Notify[i]; got != want {
					t.Errorf("#%d: notify[%d] = %s, want %s", idx, i, got, want)
				}
			}
		}
	}
}

func BenchmarkJoin(b *testing.B) {
	b.StopTimer()
	var wg sync.WaitGroup
	users := make([]string, 10000)
	chans := make([]string, 100)
	for i := range users {
		users[i] = user.NextUserID()
		if i < len(chans) {
			chans[i] = "#" + users[i]
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			channame := chans[rand.Intn(len(chans))]
			userid := users[rand.Intn(len(users))]
			channel, _ := Get(channame, true)
			channel.Join(userid, "")
			wg.Done()
		}(i)
	}
	wg.Wait()
}
