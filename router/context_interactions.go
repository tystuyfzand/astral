package router

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"strconv"
	"strings"
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
			// TODO: Is this even valid?
			// This skips values which are simply subcommand types...
			if opt.Value == nil {
				continue
			}

			for _, arg := range r.Arguments {
				argName := strings.ToLower(commandNameRe.ReplaceAllString(strings.ToLower(arg.Name), ""))

				if argName != opt.Name {
					continue
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

				break
			}
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
