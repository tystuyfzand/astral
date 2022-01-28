package router

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"strconv"
	"strings"
)

// FindInteraction finds a route path from a command interaction
func (r *Route) FindInteraction(parentRoute string, options []discord.CommandInteractionOption) *Route {
	if len(r.routes) < 1 {
		return r
	}

	opts := options

	currentRoute := r.routes[parentRoute]

	if currentRoute == nil {
		return nil
	}

	var routeName string

	for opts != nil {
		routeName, opts = recurseOptions(opts)

		if routeName != "" {
			if newRoute, exists := currentRoute.routes[routeName]; exists {
				currentRoute = newRoute
			} else {
				break
			}
		}
	}

	return currentRoute
}

func recurseOptions(options []discord.CommandInteractionOption) (string, []discord.CommandInteractionOption) {
	for _, option := range options {
		if option.Value == nil {
			return option.Name, option.Options
		}
	}

	return "", nil
}

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

		path := r.Path()
		path = path[1:]

		for _, opt := range optionsFromPath(path, data.Options) {
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

					if val[0] == '"' && val[len(val)-1] == '"' {
						val = val[1 : len(val)-1]
					}
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

func optionsFromPath(path []string, options []discord.CommandInteractionOption) []discord.CommandInteractionOption {
	if len(path) < 1 {
		return options
	}

	for _, opt := range options {
		if opt.Name == path[0] {
			// Recurse deeper until we're at path depth (path < 1)
			return optionsFromPath(path[1:], opt.Options)
		}
	}

	return nil
}
