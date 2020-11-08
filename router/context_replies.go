package router

import (
	"errors"
	"fmt"
	"github.com/diamondburned/arikawa/v2/api"
	"github.com/diamondburned/arikawa/v2/discord"
	"io"
	"strings"
)

var (
	ErrEmptyText = errors.New("text is empty")
)

// Show context usage
func (c *Context) Usage(usage ...string) (*discord.Message, error) {
	if len(usage) == 0 {
		usage = []string{c.route.Usage}
	}

	usage[0] = strings.Replace(usage[0], "{command}", c.Command, -1)

	return c.Reply(usage[0])
}

// Send text to the originating channel
func (c *Context) Send(text string) (*discord.Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	return c.Session.SendMessage(c.Channel.ID, text, nil)
}

// Send formattable text to the originating channel
func (c *Context) Sendf(format string, a ...interface{}) (*discord.Message, error) {
	return c.Send(fmt.Sprintf(format, a...))
}

// Send a file by name and read from r
func (c *Context) SendFile(name string, r io.Reader) (*discord.Message, error) {
	data := api.SendMessageData{
		Files: []api.SendMessageFile{
			{Name: name, Reader: r},
		},
	}

	return c.Session.SendMessageComplex(c.Channel.ID, data)
}

// Reply with a user mention
func (c *Context) Reply(text string) (*discord.Message, error) {
	return c.Send(fmt.Sprintf("<@%s> %s", c.User.ID, text))
}

// Reply with formatted text
func (c *Context) Replyf(format string, a ...interface{}) (*discord.Message, error) {
	return c.Reply(fmt.Sprintf(format, a...))
}

// Reply to a specific user
func (c *Context) ReplyTo(to, text string) (*discord.Message, error) {
	return c.Send(fmt.Sprintf("<@%s> %s", to, text))
}

// Reply to a user with an embed object
func (c *Context) ReplyEmbed(embed *discord.Embed) (*discord.Message, error) {
	return c.Session.SendMessage(c.Channel.ID, "<@"+c.User.ID.String()+">", embed)
}

// Reply to a user with a file object
func (c *Context) ReplyFile(name string, r io.Reader) (*discord.Message, error) {
	data := api.SendMessageData{
		Content: "<@" + c.User.ID.String() + ">",
		Files: []api.SendMessageFile{
			{Name: name, Reader: r},
		},
	}

	return c.Session.SendMessageComplex(c.Channel.ID, data)
}
