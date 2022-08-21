package astral

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
	"io"
	"strings"
)

type InteractionResponder struct {
	ctx *Context
}

// Usage builds and shows command usage
func (m *InteractionResponder) Usage(usage ...string) (*discord.Message, error) {
	if len(usage) == 0 {
		usage = []string{m.ctx.route.Usage}
	}

	usage[0] = strings.Replace(usage[0], "{command}", m.ctx.route.Name, -1)

	return m.Reply(usage[0])
}

// Send text to the originating channel
func (m *InteractionResponder) Send(text string) (*discord.Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	if err := checkMessageChannel(m.ctx); err != nil {
		return nil, err
	}

	return m.Reply(text)
}

// Sendf Sends formattable text to the originating channel
func (m *InteractionResponder) Sendf(format string, a ...interface{}) (*discord.Message, error) {
	return m.Reply(fmt.Sprintf(format, a...))
}

// SendFile sends a file by name and the data from r
func (m *InteractionResponder) SendFile(name string, r io.Reader) (*discord.Message, error) {
	data := api.SendMessageData{
		Files: []sendpart.File{
			{Name: name, Reader: r},
		},
	}

	return m.ctx.Session.SendMessageComplex(m.ctx.Channel.ID, data)
}

// Replyf Builds a message and replies with formatted text
func (m *InteractionResponder) Replyf(format string, a ...interface{}) (*discord.Message, error) {
	return m.Reply(fmt.Sprintf(format, a...))
}

// ReplyTo replies to a specific user
func (m *InteractionResponder) ReplyTo(to discord.UserID, text string) (*discord.Message, error) {
	return m.Reply(fmt.Sprintf("%s %s", to.Mention(), text))
}

// Reply with a user mention
func (m *InteractionResponder) Reply(text string) (*discord.Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	err := m.ctx.Session.RespondInteraction(m.ctx.Interaction.ID, m.ctx.Interaction.Token, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{Content: option.NewNullableString(text)},
	})

	return nil, err
}

// ReplyEmbed replies to a user with an embed object
func (m *InteractionResponder) ReplyEmbed(embed *discord.Embed) (*discord.Message, error) {
	err := m.ctx.Session.RespondInteraction(m.ctx.Interaction.ID, m.ctx.Interaction.Token, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Embeds: &[]discord.Embed{*embed},
		},
	})

	return nil, err
}

// ReplyFile replies to a user with a file object
func (m *InteractionResponder) ReplyFile(name string, r io.Reader) (*discord.Message, error) {
	err := m.ctx.Session.RespondInteraction(m.ctx.Interaction.ID, m.ctx.Interaction.Token, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Files: []sendpart.File{
				{Name: name, Reader: r},
			},
		},
	})

	return nil, err
}

// Respond replies to a user by serializing Response
func (m *InteractionResponder) Respond(r Response) (*discord.Message, error) {
	var embeds *[]discord.Embed = nil

	if r.Embeds != nil {
		embeds = &r.Embeds
	}

	var content option.NullableString = nil

	if r.Content != "" {
		content = option.NewNullableString(r.Content)
	}

	err := m.ctx.Session.RespondInteraction(m.ctx.Interaction.ID, m.ctx.Interaction.Token, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: content,
			Embeds:  embeds,
			Files:   r.Files,
		},
	})

	return nil, err
}
