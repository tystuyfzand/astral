package router

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io"
	"strings"
)

var (
	ErrEmptyText = errors.New("text is empty")
)

// Show context usage
func (c *Context) Usage(usage ...string) (*discordgo.Message, error) {
	if len(usage) == 0 {
		usage = []string{c.route.Usage}
	}

	usage[0] = strings.Replace(usage[0], "{command}", c.Command, -1)

	return c.Reply(usage[0])
}

// Send text to the originating channel
func (c *Context) Send(text string) (*discordgo.Message, error) {
	if text == "" {
		return nil, ErrEmptyText
	}

	return c.wrapMiddleware(&TextReply{Text: text})
}

// Send formattable text to the originating channel
func (c *Context) Sendf(format string, a ...interface{}) (*discordgo.Message, error) {
	return c.Send(fmt.Sprintf(format, a...))
}

// Send a file by name and read from r
func (c *Context) SendFile(name string, r io.Reader) (*discordgo.Message, error) {
	return c.Session.ChannelFileSend(c.Channel.ID, name, r)
}

// Reply with a user mention
func (c *Context) Reply(text string) (*discordgo.Message, error) {
	return c.Send(fmt.Sprintf("<@%s> %s", c.User.ID, text))
}

// Reply with formatted text
func (c *Context) Replyf(format string, a ...interface{}) (*discordgo.Message, error) {
	return c.Reply(fmt.Sprintf(format, a...))
}

// Reply to a specific user
func (c *Context) ReplyTo(to, text string) (*discordgo.Message, error) {
	return c.Send(fmt.Sprintf("<@%s> %s", to, text))
}

// Reply to a user with an embed object
func (c *Context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return c.wrapMiddleware(&EmbedReply{Embed: embed})
}

// Reply to a user with a file object
func (c *Context) ReplyFile(name string, r io.Reader) (*discordgo.Message, error) {
	return c.wrapMiddleware(&FileReply{Name: name, Reader: r})
}

// Wrap a reply with middleware
func (c *Context) wrapMiddleware(reply Reply) (*discordgo.Message, error) {
	handler := func(ctx *Context, reply Reply) (*discordgo.Message, error) {
		return ctx.SendReply(reply)
	}

	for _, m := range c.middleware {
		handler = m(handler)
	}

	return handler(c, reply)
}

// Send a reply (Text, Embed, File, etc)
func (c *Context) SendReply(reply Reply) (*discordgo.Message, error) {
	send := &discordgo.MessageSend{}
	switch v := reply.(type) {
	case *TextReply:
		send.Content = v.Text
	case *EmbedReply:
		send.Content = "<@" + c.User.ID + ">"
		send.Embed = v.Embed
	case *FileReply:
		send.Content = "<@" + c.User.ID + ">"
		send.File = &discordgo.File{Name: v.Name, Reader: v.Reader}
	}

	return c.Session.ChannelMessageSendComplex(c.Channel.ID, send)
}
