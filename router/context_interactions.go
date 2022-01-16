package router

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"strconv"
)

// ContextFromInteraction creates a new Context from an interaction event
func ContextFromInteraction(state *state.State, event *gateway.InteractionCreateEvent, r *Route) (*Context, error) {
	// Find the guild for that channel. This uses State if enabled.
	c, err := state.Channel(event.ChannelID)

	if err != nil {
		return nil, err
	}

	g, err := state.Guild(event.GuildID)

	if err != nil {
		return nil, err
	}

	args := make([]string, r.ArgumentCount)

	event.Data.InteractionType()
	switch event.Data.InteractionType() {
	case discord.CommandInteractionType:
		data := event.Data.(*discord.CommandInteraction)

		for _, opt := range data.Options {
			arg, ok := r.Arguments[opt.Name]

			if !ok {
				return nil, fmt.Errorf("%s is not a valid argument name", opt.Name)
			}

			var val string

			switch arg.Type {
			case ArgumentTypeInt:
				v, err := opt.IntValue()

				if err != nil {
					return nil, err
				}

				val = strconv.FormatInt(v, 10)
			case ArgumentTypeUserMention:
				v, err := opt.SnowflakeValue()

				if err != nil {
					return nil, err
				}

				val = discord.UserID(v).Mention()
			case ArgumentTypeChannelMention:
				v, err := opt.SnowflakeValue()

				if err != nil {
					return nil, err
				}

				val = discord.ChannelID(v).Mention()
			default:
				val = opt.Value.String()
			}

			args[arg.Index] = val
		}
	}

	ctx := &Context{
		VariableBag: NewVariableBag(),

		route:          r,
		Session:        state,
		Guild:          g,
		Channel:        c,
		User:           event.Member.User,
		Arguments:      args,
		ArgumentCount:  len(args),
		Interaction:    event,
		ArgumentString: "",
	}

	ctx.responder = &InteractionResponder{ctx}

	return ctx, nil
}
