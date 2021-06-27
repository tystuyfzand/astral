package router

import (
	"errors"
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"io"
	"strings"
)

var (
	ErrEmptyText = errors.New("text is empty")
)

// Usage builds and shows command usage
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

	if c.Channel.Type == discord.DirectMessage {
		var err error

		c.Channel, err = c.Session.CreatePrivateChannel(c.User.ID)

		if err != nil {
			return nil, err
		}
	}

	return c.Session.SendMessage(c.Channel.ID, text)
}

// Sendf Sends formattable text to the originating channel
func (c *Context) Sendf(format string, a ...interface{}) (*discord.Message, error) {
	return c.Send(fmt.Sprintf(format, a...))
}

// SendFile sends a file by name and the data from r
func (c *Context) SendFile(name string, r io.Reader) (*discord.Message, error) {
	data := api.SendMessageData{
		Files: []sendpart.File{
			{Name: name, Reader: r},
		},
	}

	return c.Session.SendMessageComplex(c.Channel.ID, data)
}

// Reply with a user mention
func (c *Context) Reply(text string) (*discord.Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	if c.Channel.Type == discord.DirectMessage {
		var err error

		c.Channel, err = c.Session.CreatePrivateChannel(c.User.ID)

		if err != nil {
			return nil, err
		}
	}

	return c.Session.SendTextReply(c.Channel.ID, text, c.Message.ID)
}

// Replyf Builds a message and replies with formatted text
func (c *Context) Replyf(format string, a ...interface{}) (*discord.Message, error) {
	return c.Reply(fmt.Sprintf(format, a...))
}

// ReplyTo replies to a specific user
func (c *Context) ReplyTo(to discord.UserID, text string) (*discord.Message, error) {
	return c.Send(fmt.Sprintf("%s %s", to.Mention(), text))
}

// Reply to a user with an embed object
func (c *Context) ReplyEmbed(embed *discord.Embed) (*discord.Message, error) {
	if c.Channel.Type == discord.DirectMessage {
		var err error

		c.Channel, err = c.Session.CreatePrivateChannel(c.User.ID)

		if err != nil {
			return nil, err
		}
	}

	return c.Session.SendEmbedReply(c.Channel.ID, *embed, c.Message.ID)
}

// Reply to a user with a file object
func (c *Context) ReplyFile(name string, r io.Reader) (*discord.Message, error) {
	data := api.SendMessageData{
		Content: "<@" + c.User.ID.String() + ">",
		Files: []sendpart.File{
			{Name: name, Reader: r},
		},
		Reference: &discord.MessageReference{MessageID: c.Message.ID},
	}

	if c.Channel.Type == discord.DirectMessage {
		var err error

		c.Channel, err = c.Session.CreatePrivateChannel(c.User.ID)

		if err != nil {
			return nil, err
		}
	}

	return c.Session.SendMessageComplex(c.Channel.ID, data)
}
