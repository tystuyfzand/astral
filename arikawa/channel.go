package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Channel struct {
	*discord.Channel
}

func NewChannel(c *discord.Channel) *Channel {
	return &Channel{c}
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
