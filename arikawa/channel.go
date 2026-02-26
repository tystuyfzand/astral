package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

type Channel struct {
	*discord.Channel
	state *state.State
}

func (c *Channel) SendMessage(msg string) (astral.Message, error) {
	m, err := c.state.SendMessage(c.Channel.ID, msg)

	if err != nil {
		return nil, err
	}

	return NewMessage(m), nil
}

func (c *Channel) Send(r astral.Response) (astral.Message, error) {
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
		Content: r.Content,
		Files:   files,
		Embeds:  embeds,
	}

	msg, err := c.state.SendMessageComplex(c.Channel.ID, data)

	if err != nil {
		return nil, err
	}

	return NewMessage(msg), nil
}

func NewChannel(state *state.State, c *discord.Channel) *Channel {
	return &Channel{
		Channel: c,
		state:   state,
	}
}

func (c *Channel) ID() astral.ID {
	return astral.ID(c.Channel.ID.String())
}

func (c *Channel) Name() string {
	return c.Channel.Name
}

func (c *Channel) Type() astral.ChannelType {
	switch c.Channel.Type {
	case discord.GuildText:
		return astral.ChannelTypeText
	case discord.GuildVoice:
		return astral.ChannelTypeVoice
	case discord.GuildCategory:
		return astral.ChannelTypeCategory
	}

	return astral.ChannelTypeUnknown
}
