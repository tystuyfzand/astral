package astral

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"io"
)

func (c *Context) Usage(usage ...string) (*discord.Message, error) {
	return c.responder.Usage(usage...)
}

func (c *Context) Send(text string) (*discord.Message, error) {
	return c.responder.Send(text)
}

func (c *Context) Sendf(format string, a ...interface{}) (*discord.Message, error) {
	return c.responder.Sendf(format, a...)
}

func (c *Context) SendFile(name string, r io.Reader) (*discord.Message, error) {
	return c.responder.SendFile(name, r)
}

func (c *Context) Reply(text string) (*discord.Message, error) {
	return c.responder.Reply(text)
}

func (c *Context) Replyf(format string, a ...interface{}) (*discord.Message, error) {
	return c.responder.Replyf(format, a...)
}

func (c *Context) ReplyTo(to discord.UserID, text string) (*discord.Message, error) {
	return c.responder.ReplyTo(to, text)
}

func (c *Context) ReplyEmbed(embed *discord.Embed) (*discord.Message, error) {
	return c.responder.ReplyEmbed(embed)
}

func (c *Context) ReplyFile(name string, r io.Reader) (*discord.Message, error) {
	return c.responder.ReplyFile(name, r)
}

func (c *Context) Respond(r Response) (*discord.Message, error) {
	return c.responder.Respond(r)
}
