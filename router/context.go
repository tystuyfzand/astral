package router

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"strconv"
	"strings"
)

type Context struct {
	route          *Route
	Session        *discordgo.Session
	Event          *discordgo.MessageCreate
	Guild          *discordgo.Guild
	Channel        *discordgo.Channel
	User           *discordgo.User
	Prefix         string
	Command        string
	ArgumentString string
	Arguments      []string
	ArgumentCount  int
	Vars           map[string]interface{}
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
