package discord

import "github.com/diamondburned/arikawa/v3/discord"

type Channel struct {
	*discord.Channel
}

func (c *Channel) ID() string {
	return c.Channel.ID.String()
}

func (c *Channel) Name() string {
	return c.Channel.Name
}
