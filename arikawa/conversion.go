package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
)

func GuildID(id astral.ServerID) discord.GuildID {
	sf, err := discord.ParseSnowflake(string(id))

	if err != nil {
		return discord.NullGuildID
	}

	return discord.GuildID(sf)
}

func ChannelID(id astral.ChannelID) discord.ChannelID {
	sf, err := discord.ParseSnowflake(string(id))

	if err != nil {
		return discord.NullChannelID
	}

	return discord.ChannelID(sf)
}

func UserID(id astral.UserID) discord.UserID {
	sf, err := discord.ParseSnowflake(string(id))

	if err != nil {
		return discord.NullUserID
	}

	return discord.UserID(sf)
}

func RoleID(id astral.RoleID) discord.RoleID {
	sf, err := discord.ParseSnowflake(string(id))

	if err != nil {
		return discord.NullRoleID
	}

	return discord.RoleID(sf)
}

func MessageID(id astral.MessageID) discord.MessageID {
	sf, err := discord.ParseSnowflake(string(id))

	if err != nil {
		return discord.NullMessageID
	}

	return discord.MessageID(sf)
}

func EmojiID(id astral.EmojiID) discord.EmojiID {
	sf, err := discord.ParseSnowflake(string(id))

	if err != nil {
		return discord.NullEmojiID
	}

	return discord.EmojiID(sf)
}
