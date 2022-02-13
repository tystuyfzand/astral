package router

import (
	"errors"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"regexp"
	"strconv"
	"strings"
)

type registrationError struct {
	cause error
	route *Route
}

func (r registrationError) Unwrap() error {
	return r.cause
}

func (r registrationError) Error() string {
	var e *httputil.HTTPError

	if errors.As(r.cause, &e) {
		return "command registration returned http error:" + string(e.Body)
	}

	return "command registration error on " + r.route.Name + ": " + r.cause.Error()
}

type commandDescriptionError struct {
	route *Route
}

func (e commandDescriptionError) Error() string {
	return "invalid command description for " + strings.Join(e.route.Path(), "->") + ": " + e.route.Description
}

type argDescriptionError struct {
	route *Route
	arg   *Argument
}

func (e argDescriptionError) Error() string {
	return "invalid argument description for " + strings.Join(e.route.Path(), "->") + " arg " + e.arg.Name + ": " + e.arg.Description
}

type argTypeError struct {
	route *Route
	arg   *Argument
}

func (e argTypeError) Error() string {
	return "invalid argument type for " + strings.Join(e.route.Path(), "->") + " arg " + e.arg.Name + ": " + strconv.Itoa(int(e.arg.Type))
}

// RegisterCommands registers all sub routes as interaction/slash commands
func RegisterCommands(r *Route, s *state.State, appID discord.AppID) ([]discord.Command, error) {
	return RegisterGuildCommands(r, s, appID, discord.NullGuildID)
}

// RegisterGuildCommands registers all sub routes as interaction/slash commands to a guild
func RegisterGuildCommands(r *Route, s *state.State, appID discord.AppID, guildID discord.GuildID) ([]discord.Command, error) {
	commands := make([]api.CreateCommandData, 0)

	for _, sub := range r.routes {
		if !sub.export {
			continue
		}

		data, err := sub.toCommandData()

		if err != nil {
			return nil, err
		}

		commands = append(commands, data)
	}

	if guildID.IsValid() {
		return s.BulkOverwriteGuildCommands(appID, guildID, commands)
	}

	return s.BulkOverwriteCommands(appID, commands)
}

func (r *Route) toCommandData() (api.CreateCommandData, error) {
	data := api.CreateCommandData{
		Name:        r.Name,
		Description: r.Description,
	}

	if r.Description == "" {
		return data, commandDescriptionError{route: r}
	}

	if len(r.routes) > 0 {
		options := make([]discord.CommandOption, len(r.routes))

		i := 0

		for _, route := range r.routes {
			inputValues, err := argsFromRoute(route)

			if err != nil {
				return data, err
			}

			values := make([]discord.CommandOptionValue, len(inputValues))

			for k, value := range inputValues {
				values[k] = value.(discord.CommandOptionValue)
			}

			options[i] = &discord.SubcommandOption{
				OptionName:  route.Name,
				Options:     values,
				Required:    false,
				Description: route.Description,
			}

			i++
		}

		data.Options = options
	} else {
		args, err := argsFromRoute(r)

		if err != nil {
			return data, err
		}

		data.Options = args
	}

	return data, nil
}

// RegisterCommand registers a single command, with sub routes as subcommands.
func (r *Route) RegisterCommand(s *state.State, appID discord.AppID, guildID discord.GuildID) (*discord.Command, error) {
	data, err := r.toCommandData()

	if err != nil {
		return nil, err
	}

	if guildID != discord.NullGuildID {
		return s.CreateGuildCommand(appID, guildID, data)
	}

	return s.CreateCommand(appID, data)
}

// UpdateCommand registers a single command, with sub routes as subcommands.
func (r *Route) UpdateCommand(s *state.State, appID discord.AppID, commandID discord.CommandID, guildID discord.GuildID) (*discord.Command, error) {
	data, err := r.toCommandData()

	if err != nil {
		return nil, err
	}

	if guildID != discord.NullGuildID {
		return s.EditGuildCommand(appID, guildID, commandID, data)
	}

	return s.EditCommand(appID, commandID, data)
}

var (
	commandNameRe = regexp.MustCompile("[^\\w-]")
)

// argsFromRoute takes a route's arguments and translates them into a discord.CommandOption
func argsFromRoute(r *Route) ([]discord.CommandOption, error) {
	options := make([]discord.CommandOption, len(r.Arguments))

	for _, arg := range r.Arguments {
		argName := strings.ToLower(commandNameRe.ReplaceAllString(strings.ToLower(arg.Name), ""))

		if arg.Description == "" {
			return nil, argDescriptionError{route: r, arg: arg}
		}

		switch arg.Type {
		case ArgumentTypeInt:
			opt := &discord.IntegerOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}

			if len(arg.Choices) > 0 {
				opt.Choices = arg.integerChoices()
			}

			options[arg.Index] = opt
		case ArgumentTypeFloat:
			opt := &discord.NumberOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}

			if len(arg.Choices) > 0 {
				opt.Choices = arg.numberChoices()
			}

			options[arg.Index] = opt
		case ArgumentTypeBool:
			options[arg.Index] = &discord.BooleanOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}
		case ArgumentTypeUserMention:
			options[arg.Index] = &discord.UserOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}
		case ArgumentTypeChannelMention:
			options[arg.Index] = &discord.ChannelOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}
		case ArgumentTypeBasic:
			opt := &discord.StringOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}

			if arg.autocomplete != nil {
				opt.Autocomplete = true
			}

			if len(arg.Choices) > 0 {
				opt.Choices = arg.stringChoices()
			}

			options[arg.Index] = opt
		default:
			return nil, argTypeError{route: r, arg: arg}
		}
	}

	return options, nil
}
