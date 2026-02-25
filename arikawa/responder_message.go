package arikawa

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

type MessageResponder struct {
	ctx     *astral.Context
	event   *gateway.MessageCreateEvent
	state   *state.State
	channel *discord.Channel
	user    discord.User
	message discord.Message
}

var (
	ErrEmptyText = errors.New("text is empty")
)

// Usage builds and shows command usage
func (m *MessageResponder) Usage(usage ...string) (astral.Message, error) {
	if len(usage) == 0 {
		usage = []string{m.ctx.Route.Usage}
	}

	usage[0] = strings.Replace(usage[0], "{command}", strings.Join(m.ctx.Route.Path(), " "), -1)

	return m.Reply(usage[0])
}

// Send text to the originating channel
func (m *MessageResponder) Send(text string) (astral.Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	if err := m.checkMessageChannel(); err != nil {
		return nil, err
	}

	msg, err := m.state.SendMessage(m.channel.ID, text)

	if err != nil {
		return nil, err
	}

	return NewMessage(msg), nil
}

// Sendf Sends formattable text to the originating channel
func (m *MessageResponder) Sendf(format string, a ...interface{}) (astral.Message, error) {
	return m.Send(fmt.Sprintf(format, a...))
}

// SendFile sends a file by name and the data from r
func (m *MessageResponder) SendFile(name string, r io.Reader) (astral.Message, error) {
	data := api.SendMessageData{
		Files: []sendpart.File{
			{Name: name, Reader: r},
		},
	}

	msg, err := m.state.SendMessageComplex(m.channel.ID, data)

	if err != nil {
		return nil, err
	}

	return NewMessage(msg), nil
}

// Replyf Builds a message and replies with formatted text
func (m *MessageResponder) Replyf(format string, a ...interface{}) (astral.Message, error) {
	return m.Reply(fmt.Sprintf(format, a...))
}

// ReplyTo replies to a specific user
func (m *MessageResponder) ReplyTo(to discord.UserID, text string) (astral.Message, error) {
	return m.Send(fmt.Sprintf("%s %s", to.Mention(), text))
}

func (m *MessageResponder) checkMessageChannel() error {
	if m.channel.Type == discord.DirectMessage {
		var err error

		m.channel, err = m.state.CreatePrivateChannel(m.user.ID)

		if err != nil {
			return err
		}
	}

	return nil
}

// Reply to a message
func (m *MessageResponder) Reply(text string) (astral.Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	if err := m.checkMessageChannel(); err != nil {
		return nil, err
	}

	msg, err := m.state.SendTextReply(m.channel.ID, text, m.message.ID)

	if err != nil {
		return nil, err
	}

	return NewMessage(msg), nil
}

// ReplyEmbed replies to a user with an embed object
func (m *MessageResponder) ReplyEmbed(embed astral.Embed) (astral.Message, error) {
	if err := m.checkMessageChannel(); err != nil {
		return nil, err
	}

	return nil, nil
	//return m.ctx.Session.SendEmbedReply(m.ctx.Channel.ID, m.ctx.Message.ID, *embed)
}

// ReplyFile replies to a user with a file object
func (m *MessageResponder) ReplyFile(name string, r io.Reader) (astral.Message, error) {
	data := api.SendMessageData{
		Content: m.user.Mention(),
		Files: []sendpart.File{
			{Name: name, Reader: r},
		},
		Reference: &discord.MessageReference{MessageID: m.message.ID},
	}

	if err := m.checkMessageChannel(); err != nil {
		return nil, err
	}

	msg, err := m.state.SendMessageComplex(m.channel.ID, data)

	if err != nil {
		return nil, err
	}

	return NewMessage(msg), nil
}

// Respond replies to a user by serializing Response
func (m *MessageResponder) Respond(r astral.Response) (astral.Message, error) {
	var files []sendpart.File

	if len(r.Files) > 0 {
		files = make([]sendpart.File, len(r.Files))

		for i, f := range r.Files {
			files[i] = sendpart.File{
				Name:   f.Name,
				Reader: f.Reader,
			}
		}
	}

	var embeds []discord.Embed

	data := api.SendMessageData{
		Content:   r.Content,
		Files:     files,
		Embeds:    embeds,
		Reference: &discord.MessageReference{MessageID: m.message.ID},
	}

	if err := m.checkMessageChannel(); err != nil {
		return nil, err
	}

	msg, err := m.state.SendMessageComplex(m.channel.ID, data)

	if err != nil {
		return nil, err
	}

	return NewMessage(msg), nil
}

func (m *MessageResponder) Acknowledge() error {
	return nil
}

func (m *MessageResponder) Error(message string) error {
	return nil
}
