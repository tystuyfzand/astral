package router

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"golang.org/x/sync/errgroup"
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
	Arguments      map[string]interface{}
	ArgumentCount  int
	responder      Responder
}

// convertedArg is an internal struct used to pass argument conversion off to a goroutine
type convertedArg struct {
	argument *Argument
	val      interface{}
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

		route:   r,
		Session: state,
		Guild:   g,
		Channel: c,
		User:    event.Author,
		Event:   event,
		Message: event.Message,
	}

	ctx.responder = &MessageResponder{ctx}

	argCh := make(chan convertedArg)

	wg := new(errgroup.Group)

	convertArg := func(arg *Argument) func() error {
		return func() error {
			convertedVal, err := ctx.convertArg(arg, args[arg.Index])

			if err != nil {
				return err
			}

			argCh <- convertedArg{arg, convertedVal}
			return nil
		}
	}

	for _, arg := range r.Arguments {
		if len(args) < arg.Index {
			continue
		}

		wg.Go(convertArg(arg))
	}

	err = wg.Wait()

	if err != nil {
		return nil, err
	}

	close(argCh)

	out := make(map[string]interface{})

	for converted := range argCh {
		out[converted.argument.Name] = converted.val
	}

	ctx.Arguments = out

	return ctx, nil
}
