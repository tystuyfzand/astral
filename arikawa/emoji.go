package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
)

func NewEmoji(emoji *discord.Emoji) *Emoji {
	return &Emoji{
		Emoji: emoji,
	}
}

type Emoji struct {
	*discord.Emoji
}

func (e *Emoji) ID() astral.EmojiID {
	return astral.EmojiID(e.Emoji.ID.String())
}

func (e *Emoji) Name() string {
	return e.Emoji.Name
}

func (e *Emoji) IsAnimated() bool {
	return e.Emoji.Animated
}
