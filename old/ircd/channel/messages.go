package channel

import (
	"bytes"
	"log"

	"github.com/kylelemons/ircd-blight/old/ircd/parser"
	"github.com/kylelemons/ircd-blight/old/ircd/user"
)

// Construct a names message for the channel.
func (c *Channel) NamesMessage(destIDs ...string) *parser.Message {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	buf := bytes.NewBuffer(nil)
	for id := range c.users {
		nick, _, _, _, ok := user.GetInfo(id)
		if !ok {
			log.Printf("Warning: Unknown id %q in %s", id, c.name)
			continue
		}
		buf.WriteByte(' ')
		buf.WriteString(nick)
	}
	buf.ReadByte()
	return &parser.Message{
		Command: parser.RPL_NAMREPLY,
		Args: []string{
			// =public *private @secret
			"*",
			"@",
			c.name,
			buf.String(),
		},
		DestIDs: destIDs,
	}
}
