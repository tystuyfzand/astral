package arikawa

import (
	"fmt"
	"io"
	"strings"

	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

func NewInteractionResponder(ctx *astral.Context, state *state.State, interaction discord.InteractionEvent) *InteractionResponder {
	return &InteractionResponder{
		ctx:         ctx,
		state:       state,
		interaction: interaction,
	}
}

type InteractionResponder struct {
	ctx          *astral.Context
	state        *state.State
	interaction  discord.InteractionEvent
	acknowledged bool
}

// Usage builds and shows command usage
func (m *InteractionResponder) Usage(usage ...string) (astral.Message, error) {
	if len(usage) == 0 {
		usage = []string{m.ctx.Route.Usage}
	}

	usage[0] = strings.Replace(usage[0], "{command}", m.ctx.Route.Name, -1)

	return m.Reply(usage[0])
}

// Send text to the originating channel
func (m *InteractionResponder) Send(text string) (astral.Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	return m.Reply(text)
}

// Sendf Sends formattable text to the originating channel
func (m *InteractionResponder) Sendf(format string, a ...interface{}) (astral.Message, error) {
	return m.Reply(fmt.Sprintf(format, a...))
}

// SendFile sends a file by name and the data from r
func (m *InteractionResponder) SendFile(name string, r io.Reader) (astral.Message, error) {
	data := api.SendMessageData{
		Files: []sendpart.File{
			{Name: name, Reader: r},
		},
	}

	msg, err := m.state.SendMessageComplex(ChannelID(m.ctx.Channel.ID()), data)

	if err != nil {
		return nil, err
	}

	return NewMessage(msg), nil
}

// Replyf Builds a message and replies with formatted text
func (m *InteractionResponder) Replyf(format string, a ...interface{}) (astral.Message, error) {
	return m.Reply(fmt.Sprintf(format, a...))
}

// ReplyTo replies to a specific user
func (m *InteractionResponder) ReplyTo(to astral.ID, text string) (astral.Message, error) {
	return m.Reply(fmt.Sprintf("%s %s", UserID(to).Mention(), text))
}

// Reply with a user mention
func (m *InteractionResponder) Reply(text string) (astral.Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	err := m.state.RespondInteraction(m.interaction.ID, m.interaction.Token, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{Content: option.NewNullableString(text)},
	})

	return nil, err
}

// ReplyEmbed replies to a user with an embed object
func (m *InteractionResponder) ReplyEmbed(embed astral.Embed) (astral.Message, error) {
	err := m.state.RespondInteraction(m.interaction.ID, m.interaction.Token, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			//Embeds: &[]discord.Embed{embed},
		},
	})

	return nil, err
}

// ReplyFile replies to a user with a file object
func (m *InteractionResponder) ReplyFile(name string, r io.Reader) (astral.Message, error) {
	err := m.state.RespondInteraction(m.interaction.ID, m.interaction.Token, api.InteractionResponse{
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
func (m *InteractionResponder) Respond(r astral.Response) (astral.Message, error) {
	var embeds *[]discord.Embed = nil

	if r.Embeds != nil {
		//embeds = &r.Embeds
	}

	var content option.NullableString = nil

	if r.Content != "" {
		content = option.NewNullableString(r.Content)
	}

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

	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: content,
			Embeds:  embeds,
			Files:   files,
		},
	}

	var err error

	if m.acknowledged {
		_, err = m.state.FollowUpInteraction(m.interaction.AppID, m.interaction.Token, *data.Data)
	} else {
		err = m.state.RespondInteraction(m.interaction.ID, m.interaction.Token, data)
	}

	return nil, err
}

func (m *InteractionResponder) Error(message string) error {
	var err error

	data := api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: option.NewNullableString(message),
			Flags:   discord.EphemeralMessage,
		},
	}

	if m.acknowledged {
		_, err = m.state.FollowUpInteraction(m.interaction.AppID, m.interaction.Token, *data.Data)
	} else {
		err = m.state.RespondInteraction(m.interaction.ID, m.interaction.Token, data)
	}

	return err
}

func (m *InteractionResponder) Acknowledge() error {
	err := m.state.RespondInteraction(m.interaction.ID, m.interaction.Token, api.InteractionResponse{
		Type: api.DeferredMessageInteractionWithSource,
	})

	if err == nil {
		m.acknowledged = true
	}

	return err
}
