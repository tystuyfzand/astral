package valour

import (
	"github.com/auroradevllc/astral/v3"
	valour "github.com/auroradevllc/valourgo"
)

var _ astral.Channel = (*Channel)(nil)

type Channel struct {
	*valour.Channel
	client valour.Client
}

func (c *Channel) ID() astral.ChannelID {
	return astral.ChannelID(c.Channel.ID.String())
}

func (c *Channel) Name() string {
	return c.Channel.Name
}

func (c *Channel) Type() astral.ChannelType {
	switch c.Channel.ChannelType {
	case valour.PlanetChat:
		return astral.ChannelTypeText
	case valour.PlanetVoice:
		return astral.ChannelTypeVoice
	case valour.PlanetCategory:
		return astral.ChannelTypeCategory
	case valour.DirectChat, valour.DirectVoice:
		return astral.ChannelTypeDirect
	case valour.GroupChat, valour.GroupVoice:
		return astral.ChannelTypeGroup
	}

	return astral.ChannelTypeUnknown
}

func (c *Channel) IsNSFW() bool {
	return c.Channel.NSFW
}

func (c *Channel) SendMessage(msg string) (astral.Message, error) {
	m, err := c.client.SendMessage(c.Channel.PlanetID, c.Channel.ID, msg)

	if err != nil {
		return nil, err
	}

	return NewMessage(m), nil
}

func (c *Channel) Send(r astral.MessageContent) (astral.Message, error) {
	return nil, nil
}

func (c *Channel) EditMessage(id astral.MessageID, msg astral.MessageContent) (astral.Message, error) {
	return nil, nil
}

func (c *Channel) DeleteMessage(id astral.MessageID) error {
	return nil
}

func NewChannel(c valour.Client, ch *valour.Channel) *Channel {
	return &Channel{
		Channel: ch,
		client:  c,
	}
}
