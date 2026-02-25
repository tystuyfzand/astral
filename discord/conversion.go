package discord

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
)

func Snowflake(id astral.ID) (discord.Snowflake, error) {
	return discord.ParseSnowflake(string(id))
}

func GuildID(id astral.ID) discord.GuildID {
	sf, err := Snowflake(id)

	if err != nil {
		return discord.NullGuildID
	}

	return discord.GuildID(sf)
}

func ChannelID(id astral.ID) discord.ChannelID {
	sf, err := Snowflake(id)

	if err != nil {
		return discord.NullChannelID
	}

	return discord.ChannelID(sf)
}

func UserID(id astral.ID) discord.UserID {
	sf, err := Snowflake(id)

	if err != nil {
		return discord.NullUserID
	}

	return discord.UserID(sf)
}
