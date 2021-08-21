package router

import (
	"github.com/diamondburned/arikawa/v3/discord"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
	"strconv"
)

// Find the specified argument nand return the information and value
func (c *Context) arg(name string) (*Argument, string) {
	if arg, exists := c.route.Arguments[name]; exists {
		if arg.Index > c.ArgumentCount-1 {
			return arg, ""
		}

		return arg, c.Arguments[arg.Index]
	}

	panic("undefined argument " + name)
}

// Arg finds and returns a named argument as a string
func (c *Context) Arg(name string) string {
	_, val := c.arg(name)

	return val
}

// IntArg finds and returns a named int argument
func (c *Context) IntArg(name string) int64 {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeInt {
		panic("Trying to use a non-int argument as int")
	}

	v, err := strconv.ParseInt(val, 10, 64)

	if err != nil {
		return -1
	}

	return v
}

// FloatArg finds and returns a named float argument
func (c *Context) FloatArg(name string) float64 {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeFloat {
		panic("Trying to use a non-float argument as float")
	}

	v, err := strconv.ParseFloat(val, 64)

	if err != nil {
		return -1
	}

	return v
}

// BoolArg finds and returns a named bool argument
func (c *Context) BoolArg(name string) bool {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeBool {
		panic("Trying to use a non-bool argument as bool")
	}

	v, err := strconv.ParseBool(val)

	if err != nil {
		return false
	}

	return v
}

// UserArg finds and returns a named User argument
func (c *Context) UserArg(name string) *discord.User {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeUserMention {
		panic("Trying to use a non-user argument as user")
	}

	m := userMentionRegexp.FindStringSubmatch(val)

	if m == nil {
		return nil
	}

	sf, err := discord.ParseSnowflake(m[1])

	if err != nil {
		return nil
	}

	u, err := c.Session.User(discord.UserID(sf))

	if err != nil {
		return nil
	}

	return u
}

// ChannelArg finds and returns a named Channel argument
func (c *Context) ChannelArg(name string) *discord.Channel {
	return c.ChannelArgType(name, 255)
}

// ChannelArgType finds and returns Channel argument with a specified type
func (c *Context) ChannelArgType(name string, t discord.ChannelType) *discord.Channel {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeChannelMention {
		panic("Trying to use a non-channel argument as channel")
	}

	m := channelMentionRegexp.FindStringSubmatch(val)

	if m != nil {
		sf, err := discord.ParseSnowflake(m[1])

		if err != nil {
			return nil
		}

		c, err := c.Session.Channel(discord.ChannelID(sf))

		if err != nil {
			return nil
		}

		if t != 255 && c.Type != t {
			return nil
		}

		return c
	}

	return nil
}

// EmojiArg finds and returns an argument as an emoji
func (c *Context) EmojiArg(name string) *discord.Emoji {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeEmoji {
		panic("Trying to use a non-emoji argument as emoji")
	}

	m := emojiRegexp.FindStringSubmatch(val)

	if m != nil {
		sf, err := discord.ParseSnowflake(m[3])

		if err != nil {
			return nil
		}

		return &discord.Emoji{
			ID:       discord.EmojiID(sf),
			Name:     m[2],
			Animated: m[1] == "a",
		}
	}

	result, err := emoji.LookupEmoji(val)

	if err == nil {
		return &discord.Emoji{
			Name: result.Value,
		}
	}

	return nil
}
