package router

import (
	"github.com/bwmarrin/discordgo"
	"io"
)

// A sender function to send a Discord message.
type SendFunc func(*Context, Reply) (*discordgo.Message, error)

// A middleware function for catching/modifying replies
type ContextMiddlewareFunc func(SendFunc) SendFunc

// A generic reply to Discord/Users
type Reply interface {
}

// A reply containing only text
type TextReply struct {
	Reply
	Text string
}

// A reply containing a message Embed
type EmbedReply struct {
	Reply
	Embed *discordgo.MessageEmbed
}

// A reply containing a File
type FileReply struct {
	Reply
	Name   string
	Reader io.Reader
}
