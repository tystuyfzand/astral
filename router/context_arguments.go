package router

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

// Find the specified argument nand return the information and value
func (c *Context) arg(name string) (*Argument, string) {
	if arg, exists := c.route.Arguments[name]; exists {
		if arg.Index > len(c.Arguments)-1 {
			return arg, ""
		}

		return arg, c.Arguments[arg.Index]
	}

	panic("undefined argument " + name)
}

// Find and return a named argument
func (c *Context) Arg(name string) string {
	_, val := c.arg(name)

	return val
}

// Find and return a named int argument
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

// Find and return a named float argument
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

// Find and return a named bool argument
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

// Find and return a named User argument
func (c *Context) UserArg(name string) *discordgo.User {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeUserMention {
		panic("Trying to use a non-user argument as user")
	}

	m := userMentionRegexp.FindStringSubmatch(val)

	if m == nil {
		return nil
	}

	u, err := c.Session.User(m[1])

	if err != nil {
		return nil
	}

	return u
}

// Find and return a named Channel argument
func (c *Context) ChannelArg(name string) *discordgo.Channel {
	return c.ChannelArgType(name, -1)
}

// Find and return a named Channel argument with a specified type
func (c *Context) ChannelArgType(name string, t discordgo.ChannelType) *discordgo.Channel {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeChannelMention {
		panic("Trying to use a non-channel argument as channel")
	}

	m := channelMentionRegexp.FindStringSubmatch(val)

	if m != nil {
		c, err := c.Session.Channel(m[1])

		if err != nil {
			return nil
		}

		return c
	}

	for _, ch := range c.Guild.Channels {
		if strings.ToLower(ch.Name) == strings.ToLower(val) && (t == -1 || ch.Type == t) {
			c, err := c.Session.Channel(ch.ID)

			if err != nil {
				return nil
			}

			return c
		}
	}

	return nil
}
