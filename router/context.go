package router

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"io"
)

type Responder interface {
	Usage(usage ...string) (*discord.Message, error)
	Send(text string) (*discord.Message, error)
	Sendf(format string, a ...interface{}) (*discord.Message, error)
	SendFile(name string, r io.Reader) (*discord.Message, error)
	Reply(text string) (*discord.Message, error)
	Replyf(format string, a ...interface{}) (*discord.Message, error)
	ReplyTo(to discord.UserID, text string) (*discord.Message, error)
	ReplyEmbed(embed *discord.Embed) (*discord.Message, error)
	ReplyFile(name string, r io.Reader) (*discord.Message, error)
}

// Context is the base "context" object.
// It contains all fields that are present on both Messages and Interactions.
type Context struct {
	*VariableBag

	route          *Route
	Session        *state.State
	Event          *gateway.MessageCreateEvent
	Interaction    *gateway.InteractionCreateEvent
	Guild          *discord.Guild
	Channel        *discord.Channel
	Message        discord.Message
	User           discord.User
	Prefix         string
	Command        string
	ArgumentString string
	Arguments      []string
	ArgumentCount  int
	responder      Responder
}

// ContextFrom creates a new MessageContext from the session and event
func ContextFrom(state *state.State, event *gateway.MessageCreateEvent, r *Route, args []string, argString string) (*Context, error) {
	// Find the channel for the event, which doesn't have a built-in discordgo equivalent of .Guild()
	c, err := state.Channel(event.ChannelID)

	if err != nil {
		return nil, err
	}

	var g *discord.Guild

	if c.Type != discord.DirectMessage {
		// Find the guild for that channel. This uses State if enabled.
		g, err = state.Guild(c.GuildID)

		if err != nil {
			return nil, err
		}
	}

	ctx := &Context{
		VariableBag: NewVariableBag(),

		route:          r,
		Session:        state,
		Guild:          g,
		Channel:        c,
		User:           event.Author,
		Arguments:      args,
		ArgumentCount:  len(args),
		Event:          event,
		Message:        event.Message,
		ArgumentString: argString,
	}

	ctx.responder = &MessageResponder{ctx}

	return ctx, nil
}
