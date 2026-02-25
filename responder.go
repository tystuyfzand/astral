package astral

import (
	"io"

	"github.com/diamondburned/arikawa/v3/discord"
)

// Responder represents an available responder
// This can send messages to channels, direct messages, interactions, etc.
type Responder interface {
	Usage(usage ...string) (Message, error)
	Send(text string) (Message, error)
	Sendf(format string, a ...interface{}) (Message, error)
	SendFile(name string, r io.Reader) (Message, error)
	Reply(text string) (Message, error)
	Replyf(format string, a ...interface{}) (Message, error)
	ReplyTo(to discord.UserID, text string) (Message, error)
	ReplyEmbed(embed Embed) (Message, error)
	ReplyFile(name string, r io.Reader) (Message, error)
	Respond(r Response) (Message, error)
	Acknowledge() error
	Error(message string) error
}
