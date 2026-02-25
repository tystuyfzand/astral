package discord

import (
	"github.com/auroradevllc/astral/v3/adapter"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Channel struct {
	*discord.Channel
}

func (c *Channel) ID() adapter.ID {
	return adapter.ID(c.Channel.ID.String())
}

func (c *Channel) Name() string {
	return c.Channel.Name
}

func (c *Channel) Type() adapter.ChannelType {
	switch c.Channel.Type {
	case discord.GuildText:
		return adapter.ChannelTypeText
	case discord.GuildVoice:
		return adapter.ChannelTypeVoice
	case discord.GuildCategory:
		return adapter.ChannelTypeCategory
	}

	return adapter.ChannelTypeUnknown
}
