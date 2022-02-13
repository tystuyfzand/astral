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

type MessageResponder struct {
	ctx *Context
}

var (
	ErrEmptyText = errors.New("text is empty")
)

// Usage builds and shows command usage
func (m *MessageResponder) Usage(usage ...string) (*discord.Message, error) {
	if len(usage) == 0 {
		usage = []string{m.ctx.route.Usage}
	}

	usage[0] = strings.Replace(usage[0], "{command}", strings.Join(m.ctx.route.Path(), " "), -1)

	return m.Reply(usage[0])
}

// Send text to the originating channel
func (m *MessageResponder) Send(text string) (*discord.Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	if err := checkMessageChannel(m.ctx); err != nil {
		return nil, err
	}

	return m.ctx.Session.SendMessage(m.ctx.Channel.ID, text)
}

// Sendf Sends formattable text to the originating channel
func (m *MessageResponder) Sendf(format string, a ...interface{}) (*discord.Message, error) {
	return m.Send(fmt.Sprintf(format, a...))
}

// SendFile sends a file by name and the data from r
func (m *MessageResponder) SendFile(name string, r io.Reader) (*discord.Message, error) {
	data := api.SendMessageData{
		Files: []sendpart.File{
			{Name: name, Reader: r},
		},
	}

	return m.ctx.Session.SendMessageComplex(m.ctx.Channel.ID, data)
}

// Replyf Builds a message and replies with formatted text
func (m *MessageResponder) Replyf(format string, a ...interface{}) (*discord.Message, error) {
	return m.Reply(fmt.Sprintf(format, a...))
}

// ReplyTo replies to a specific user
func (m *MessageResponder) ReplyTo(to discord.UserID, text string) (*discord.Message, error) {
	return m.Send(fmt.Sprintf("%s %s", to.Mention(), text))
}

func checkMessageChannel(ctx *Context) error {
	if ctx.Channel.Type == discord.DirectMessage {
		var err error

		ctx.Channel, err = ctx.Session.CreatePrivateChannel(ctx.User.ID)

		if err != nil {
			return err
		}
	}

	return nil
}

// Reply to a message
func (m *MessageResponder) Reply(text string) (*discord.Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	if err := checkMessageChannel(m.ctx); err != nil {
		return nil, err
	}

	return m.ctx.Session.SendTextReply(m.ctx.Channel.ID, text, m.ctx.Message.ID)
}

// ReplyEmbed replies to a user with an embed object
func (m *MessageResponder) ReplyEmbed(embed *discord.Embed) (*discord.Message, error) {
	if err := checkMessageChannel(m.ctx); err != nil {
		return nil, err
	}

	return m.ctx.Session.SendEmbedReply(m.ctx.Channel.ID, m.ctx.Message.ID, *embed)
}

// ReplyFile replies to a user with a file object
func (m *MessageResponder) ReplyFile(name string, r io.Reader) (*discord.Message, error) {
	data := api.SendMessageData{
		Content: "<@" + m.ctx.User.ID.String() + ">",
		Files: []sendpart.File{
			{Name: name, Reader: r},
		},
		Reference: &discord.MessageReference{MessageID: m.ctx.Message.ID},
	}

	if err := checkMessageChannel(m.ctx); err != nil {
		return nil, err
	}

	return m.ctx.Session.SendMessageComplex(m.ctx.Channel.ID, data)
}
