package router

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"strings"
)

type Context struct {
	route *Route
	Session *discordgo.Session
	Event *discordgo.MessageCreate
	Guild *discordgo.Guild
	Channel *discordgo.Channel
	User *discordgo.User
	Prefix string
	Command string
	ArgumentString string
	Arguments []string
	ArgumentCount int
	Vars map[string]interface{}
}

// Set sets a variable on the context
func (c *Context) Set(key string, d interface{}) {
	c.Vars[key] = d
}

// Get retrieves a variable from the context
func (c *Context) Get(key string) interface{} {
	if c, ok := c.Vars[key]; ok {
		return c
	}
	return nil
}

// Send text to the originating channel
func (c *Context) Send(text string) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSend(c.Channel.ID, text)
}

// Send formattable text to the originating channel
func (c *Context) Sendf(format string, a ...interface{}) (*discordgo.Message, error) {
	return c.Send(fmt.Sprintf(format, a...))
}

// Send a file by name and read from r
func (c *Context) SendFile(name string, r io.Reader) (*discordgo.Message, error) {
	return c.Session.ChannelFileSend(c.Channel.ID, name, r)
}

// Show context usage
func (c *Context) Usage(usage ...string) (*discordgo.Message, error) {
	if len(usage) == 0 {
		usage = []string{c.route.Usage}
	}

	usage[0] = strings.Replace(usage[0], "{command}", c.Command, -1)

	return c.Reply(usage[0])
}

// Reply with a user mention
func (c *Context) Reply(text string) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSend(c.Channel.ID, fmt.Sprintf("<@%s> %s", c.User.ID, text))
}

// Reply with formatted text
func (c *Context) Replyf(format string, a ...interface{}) (*discordgo.Message, error) {
	return c.Reply(fmt.Sprintf(format, a...))
}

// Reply to a specific user
func (c *Context) ReplyTo(to, text string) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSend(c.Channel.ID, fmt.Sprintf("<@%s> %s", to, text))
}

// Reply to a user with an embed object
func (c *Context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendComplex(c.Channel.ID, &discordgo.MessageSend{Content: "<@" + c.User.ID + ">", Embed: embed})
}

// Reply to a user with a file object
func (c *Context) ReplyFile(name string, r io.Reader) (*discordgo.Message, error) {
	return c.Session.ChannelMessageSendComplex(c.Channel.ID, &discordgo.MessageSend{Content: "<@" + c.User.ID + ">", File: &discordgo.File{Name: name, Reader: r}})
}

// Find and return a named argument
func (c *Context) Argument(name string) string {
	if arg, exists := c.route.Arguments[name]; exists {
		if arg.Index > len(c.Arguments) - 1 {
			return ""
		}

		return c.Arguments[arg.Index]
	}
	return ""
}

// UserArgument parses and returns a *discordgo.User for the specified argument name
func (c *Context) UserArgument(name string) *discordgo.User {
	if arg, exists := c.route.Arguments[name]; exists {
		if arg.Index > len(c.Arguments) - 1 {
			return nil
		}

		m := userMentionRegexp.FindStringSubmatch(c.Arguments[arg.Index])

		if m == nil {
			return nil
		}

		u, err := c.Session.User(m[1])

		if err != nil {
			return nil
		}

		return u
	}
	return nil
}

// ChannelArgument returns the first found Channel from the argument name
func (c *Context) ChannelArgument(name string) *discordgo.Channel {
	return c.ChannelArgumentType(name, -1)
}

// ChannelArgumentType returns the first found Channel of the specified type
func (c *Context) ChannelArgumentType(name string, t discordgo.ChannelType) *discordgo.Channel {
	if arg, exists := c.route.Arguments[name]; exists {
		if arg.Index > len(c.Arguments) - 1 {
			return nil
		}

		channelName := c.Arguments[arg.Index]

		m := channelMentionRegexp.FindStringSubmatch(channelName)

		if m != nil {
			c, err := c.Session.Channel(m[1])

			if err != nil {
				return nil
			}

			return c
		}

		for _, ch := range c.Guild.Channels {
			if strings.ToLower(ch.Name) == strings.ToLower(channelName) && (t == -1 || ch.Type == t) {
				c, err := c.Session.Channel(ch.ID)

				if err != nil {
					return nil
				}

				return c
			}
		}

		return nil
	}
	return nil
}