package discord

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"golang.org/x/sync/errgroup"
	"sync"
)

// ContextFrom creates a new MessageContext from the session and event
func ContextFrom(state *state.State, event *gateway.MessageCreateEvent, r *astral.Route, args []string) (*Context, error) {
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

	ctx := astral.NewContext(
		r,
		astral.WithResponder(&MessageResponder{ctx})
	)
	ctx := &astral.Context{
		VariableBag: astral.NewVariableBag(),

		Route:   r,
		Session: state,
		Server:  &Server{Guild: g},
		Channel: &Channel{Channel: c},
		User:    NewUser(&event.Author),
		Message: NewMessage(event.Message),
	}

	ctx.responder = &MessageResponder{
		ctx: ctx,
		Event: event,
	}

	wg := new(errgroup.Group)

	out := make(map[string]interface{})

	var outLock sync.Mutex

	convertArg := func(arg *Argument) func() error {
		return func() error {
			convertedVal, err := ctx.convertArg(arg, args[arg.Index])

			if err != nil {
				return err
			}

			outLock.Lock()
			out[arg.Name] = convertedVal
			outLock.Unlock()
			return nil
		}
	}

	for _, arg := range r.Arguments {
		if len(args) <= arg.Index {
			continue
		}

		wg.Go(convertArg(arg))
	}

	err = wg.Wait()

	if err != nil {
		return nil, err
	}

	ctx.Arguments = out

	return ctx, nil
}
