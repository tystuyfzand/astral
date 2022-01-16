package router

import (
	"encoding/json"
	"errors"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
	"log"
	"regexp"
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

// RegisterCommands registers all sub routes as interaction/slash commands
func RegisterCommands(r *Route, s *state.State, appID discord.AppID) ([]*discord.Command, error) {
	return RegisterGuildCommands(r, s, appID, discord.NullGuildID)
}

// RegisterGuildCommands registers all sub routes as interaction/slash commands to a guild
func RegisterGuildCommands(r *Route, s *state.State, appID discord.AppID, guildID discord.GuildID) ([]*discord.Command, error) {
	commands := make([]*discord.Command, 0)

	// Pull existing commands
	existing, err := s.GuildCommands(appID, guildID)

	if err != nil {
		return nil, err
	}

	existingMap := make(map[string]discord.Command)

	for _, cmd := range existing {
		existingMap[cmd.Name] = cmd

		log.Println("Found exiting command", cmd.Name)
	}

	for _, sub := range r.routes {
		if !sub.export {
			continue
		}

		var command *discord.Command

		if cmd, exists := existingMap[sub.Name]; exists {
			command, err = sub.UpdateCommand(s, appID, cmd.ID, guildID)
		} else {
			command, err = sub.RegisterCommand(s, appID, guildID)
		}

		if err != nil {
			return nil, registrationError{err, sub}
		}

		commands = append(commands, command)
	}

	return commands, nil
}

func (r *Route) toCommandData() api.CreateCommandData {
	data := api.CreateCommandData{
		Name:        r.Name,
		Description: r.Description,
	}

	if len(r.routes) > 0 {
		options := make([]discord.CommandOption, len(r.routes))

		i := 0

		for _, route := range r.routes {
			values := make([]discord.CommandOptionValue, 0)

			for k, value := range argsFromRoute(route) {
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
		data.Options = argsFromRoute(r)
	}

	return data
}

// RegisterCommand registers a single command, with sub routes as subcommands.
func (r *Route) RegisterCommand(s *state.State, appID discord.AppID, guildID discord.GuildID) (*discord.Command, error) {
	data := r.toCommandData()

	if guildID != discord.NullGuildID {
		return s.CreateGuildCommand(appID, guildID, data)
	}

	return s.CreateCommand(appID, data)
}

// UpdateCommand registers a single command, with sub routes as subcommands.
func (r *Route) UpdateCommand(s *state.State, appID discord.AppID, commandID discord.CommandID, guildID discord.GuildID) (*discord.Command, error) {
	data := r.toCommandData()

	b, _ := json.Marshal(data)

	log.Println(string(b))

	if guildID != discord.NullGuildID {
		return s.EditGuildCommand(appID, guildID, commandID, data)
	}

	return s.EditCommand(appID, commandID, data)
}

var (
	commandNameRe = regexp.MustCompile("[^\\w-]")
)

// argsFromRoute takes a route's arguments and translates them into a discord.CommandOption
func argsFromRoute(r *Route) []discord.CommandOption {
	options := make([]discord.CommandOption, len(r.Arguments))

	for _, arg := range r.Arguments {
		argName := strings.ToLower(commandNameRe.ReplaceAllString(strings.ToLower(arg.Name), ""))

		switch arg.Type {
		case ArgumentTypeInt:
			options[arg.Index] = &discord.IntegerOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: "The " + arg.Name + " argument",
			}
		case ArgumentTypeFloat:
			options[arg.Index] = &discord.NumberOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: "The " + arg.Name + " argument",
			}
		case ArgumentTypeBool:
			options[arg.Index] = &discord.BooleanOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: "The " + arg.Name + " argument",
			}
		case ArgumentTypeUserMention:
			options[arg.Index] = &discord.UserOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: "The " + arg.Name + " argument",
			}
		case ArgumentTypeChannelMention:
			options[arg.Index] = &discord.ChannelOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: "The " + arg.Name + " argument",
			}
		case ArgumentTypeBasic:
			options[arg.Index] = &discord.StringOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: "The " + arg.Name + " argument",
			}
		}
	}

	return options
}
