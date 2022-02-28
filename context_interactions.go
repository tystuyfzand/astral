package astral

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
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
		if option.Type == discord.SubcommandOptionType || option.Type == discord.SubcommandGroupOptionType {
			return option.Name, option.Options
		}
	}

	return "", nil
}

// FindAutocomplete finds a route path from an autocomplete interaction
func (r *Route) FindAutocomplete(parentRoute string, options []discord.AutocompleteOption) (*Route, []discord.AutocompleteOption) {
	if len(r.routes) < 1 {
		return r, nil
	}

	opts := options

	currentRoute := r.routes[parentRoute]

	if currentRoute == nil {
		return nil, nil
	}

	var routeName string
	var focused bool

	for opts != nil {
		routeName, opts, focused = recurseAutocompleteOptions(opts)

		if routeName != "" {
			if newRoute, exists := currentRoute.routes[routeName]; exists {
				currentRoute = newRoute
			} else {
				break
			}
		}

		if focused {
			break
		}
	}

	return currentRoute, opts
}

func recurseAutocompleteOptions(options []discord.AutocompleteOption) (string, []discord.AutocompleteOption, bool) {
	for _, option := range options {
		if option.Type == discord.SubcommandOptionType {
			return option.Name, option.Options, false
		}

		foundFocused := option.Focused

		if option.Options != nil {
			for _, opt := range option.Options {
				if opt.Focused {
					foundFocused = true
					break
				}
			}
		}

		if foundFocused {
			return "", options, option.Focused
		}
	}

	return "", nil, false
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

	ctx := &Context{
		VariableBag: NewVariableBag(),

		route:          r,
		Session:        state,
		Guild:          g,
		Channel:        c,
		User:           event.Member.User,
		Interaction:    event,
		ArgumentString: "",
	}

	ctx.responder = &InteractionResponder{ctx}

	switch data := event.Data.(type) {
	case *discord.CommandInteraction:
		path := r.Path()
		path = path[1:]

		wg := new(errgroup.Group)

		out := make(map[string]interface{})
		var outLock sync.Mutex

		checkArg := func(opt discord.CommandInteractionOption) func() error {
			return func() error {
				for _, arg := range r.Arguments {
					argName := strings.ToLower(commandNameRe.ReplaceAllString(strings.ToLower(arg.Name), ""))

					if argName != opt.Name {
						continue
					}

					switch arg.Type {
					case ArgumentTypeInt:
						v, err := opt.IntValue()

						if err != nil {
							return err
						}

						outLock.Lock()
						out[arg.Name] = v
						outLock.Unlock()
					case ArgumentTypeUserMention:
						v, err := opt.SnowflakeValue()

						if err != nil {
							return err
						}

						convertedVal, err := ctx.convertArg(arg, v)

						if err != nil {
							return err
						}

						outLock.Lock()
						out[arg.Name] = convertedVal
						outLock.Unlock()
					case ArgumentTypeChannelMention:
						v, err := opt.SnowflakeValue()

						if err != nil {
							return err
						}

						convertedVal, err := ctx.convertArg(arg, v)

						if err != nil {
							return err
						}

						outLock.Lock()
						out[arg.Name] = convertedVal
						outLock.Unlock()
					default:
						val := opt.Value.String()

						if val[0] == '"' && val[len(val)-1] == '"' {
							val = val[1 : len(val)-1]
						}

						convertedVal, err := ctx.convertArg(arg, val)

						if err != nil {
							return err
						}

						outLock.Lock()
						out[arg.Name] = convertedVal
						outLock.Unlock()
					}

					break
				}

				return nil
			}
		}

		for _, opt := range optionsFromPath(path, data.Options) {
			wg.Go(checkArg(opt))
		}

		err := wg.Wait()

		if err != nil {
			return nil, err
		}

		ctx.Arguments = out
	case *discord.AutocompleteInteraction:
	}

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
