package router

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"github.com/diamondburned/arikawa/v2/session"
)

type Context struct {
	route          *Route
	Session        *session.Session
	Event          *gateway.MessageCreateEvent
	Guild          *discord.Guild
	Channel        *discord.Channel
	User           discord.User
	Prefix         string
	Command        string
	ArgumentString string
	Arguments      []string
	ArgumentCount  int
	Vars           map[string]interface{}
}

// Create a new Context from the session and event
func ContextFrom(session *session.Session, event *gateway.MessageCreateEvent, r *Route, command string, args []string, argString string) (*Context, error) {
	// Find the channel for the event, which doesn't have a built-in discordgo equivalent of .Guild()
	c, err := session.Channel(event.ChannelID)

	if err != nil {
		return nil, err
	}

	// Find the guild for that channel. This uses State if enabled.
	g, err := session.Guild(c.GuildID)

	if err != nil {
		return nil, err
	}

	ctx := &Context{
		route:          r,
		Session:        session,
		Event:          event,
		Guild:          g,
		Channel:        c,
		User:           event.Author,
		Command:        command,
		ArgumentString: argString,
		Arguments:      args,
		ArgumentCount:  len(args),
		Vars:           make(map[string]interface{}),
	}

	return ctx, nil
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
