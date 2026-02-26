package arikawa

import (
	"strings"
	"sync"

	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"golang.org/x/sync/errgroup"
)

// ContextFromInteraction creates a new Context from an interaction event
func ContextFromInteraction(state *state.State, event *gateway.InteractionCreateEvent, r *astral.Route) (*astral.Context, error) {
	// Find the guild for that channel. This uses State if enabled.
	c, err := state.Channel(event.ChannelID)

	if err != nil {
		return nil, err
	}

	g, err := state.Guild(event.GuildID)

	if err != nil {
		return nil, err
	}

	ctx := astral.NewContext(r,
		astral.WithClient(&Client{state: state}),
		astral.WithServer(NewServer(state, g)),
		astral.WithChannel(NewChannel(state, c)),
		astral.WithUser(NewUser(event.User)),
		astral.WithResponder(&InteractionResponder{
			state: state,
		}),
	)

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
					case astral.ArgumentTypeInt:
						v, err := opt.IntValue()

						if err != nil {
							return err
						}

						outLock.Lock()
						out[arg.Name] = v
						outLock.Unlock()
					case astral.ArgumentTypeUserMention:
						v, err := opt.SnowflakeValue()

						if err != nil {
							return err
						}

						convertedVal, err := ctx.ConvertArg(arg, v)

						if err != nil {
							return err
						}

						outLock.Lock()
						out[arg.Name] = convertedVal
						outLock.Unlock()
					case astral.ArgumentTypeChannelMention:
						v, err := opt.SnowflakeValue()

						if err != nil {
							return err
						}

						convertedVal, err := ctx.ConvertArg(arg, v)

						if err != nil {
							return err
						}

						outLock.Lock()
						out[arg.Name] = convertedVal
						outLock.Unlock()
					case astral.ArgumentTypeRole:
						v, err := opt.SnowflakeValue()

						if err != nil {
							return err
						}

						convertedVal, err := ctx.ConvertArg(arg, v)

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

						convertedVal, err := ctx.ConvertArg(arg, val)

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
