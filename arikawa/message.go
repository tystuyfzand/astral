package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
)

func NewMessage(m *discord.Message) *Message {
	return &Message{
		Message: m,
	}
}

type Message struct {
	*discord.Message
}

func (m *Message) ID() astral.MessageID {
	return astral.MessageID(m.Message.ID.String())
}

func (m *Message) Content() string {
	return m.Message.Content
}
