package astral

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/sendpart"
)

type Response struct {
	Content string
	Embeds  []discord.Embed
	Files   []sendpart.File
}
