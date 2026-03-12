package valour

import (
	"github.com/auroradevllc/astral/v3"
	valour "github.com/auroradevllc/valourgo"
)

var _ astral.Message = (*Message)(nil)

func NewMessage(m *valour.Message) *Message {
	return &Message{
		Message: m,
	}
}

type Message struct {
	*valour.Message
}

func (m *Message) ID() astral.MessageID {
	return astral.MessageID(m.Message.ID.String())
}

func (m *Message) Content() string {
	return m.Message.Content
}
