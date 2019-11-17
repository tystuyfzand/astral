package router

import (
	"github.com/bwmarrin/discordgo"
)

type FilterFunc func(string) string

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
	Filters        []FilterFunc
}

// Create a new Context from the session and event
func ContextFrom(session *discordgo.Session, event *discordgo.MessageCreate, r *Route, command string, args []string, argString string) (*Context, error) {
	// Find the channel for the event, which doesn't have a built-in discordgo equivalent of .Guild()
	c, err := channel(session, event.ChannelID)

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
		User:           event.Author,
		Command:        command,
		ArgumentString: argString,
		Arguments:      args,
		ArgumentCount:  len(args),
		Vars:           make(map[string]interface{}),
		Filters:        make([]FilterFunc, 0),
	}

	return ctx, nil
}

// Add a filter to the context
func (c *Context) Use(f ...FilterFunc) *Context {
	c.Filters = append(c.Filters, f...)
	return c
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

// Retrieve a channel from state or rest api
func channel(session *discordgo.Session, channelID string) (c *discordgo.Channel, err error) {
	if session.StateEnabled {
		c, err = session.State.Channel(channelID)

		if err == nil && c != nil {
			return
		}
	}

	c, err = session.Channel(channelID)
	return
}
