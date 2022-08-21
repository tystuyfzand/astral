package astral

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"io"
)

// Responder represents an available responder
// This can send messages to channels, direct messages, interactions, etc.
type Responder interface {
	Usage(usage ...string) (*discord.Message, error)
	Send(text string) (*discord.Message, error)
	Sendf(format string, a ...interface{}) (*discord.Message, error)
	SendFile(name string, r io.Reader) (*discord.Message, error)
	Reply(text string) (*discord.Message, error)
	Replyf(format string, a ...interface{}) (*discord.Message, error)
	ReplyTo(to discord.UserID, text string) (*discord.Message, error)
	ReplyEmbed(embed *discord.Embed) (*discord.Message, error)
	ReplyFile(name string, r io.Reader) (*discord.Message, error)
	Respond(r Response) (*discord.Message, error)
}

type Response struct {
	Content string
	Embeds  []discord.Embed
	Files   []sendpart.File
}
