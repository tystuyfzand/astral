package arikawa

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/httputil"
)

type registrationError struct {
	cause error
	route *astral.Route
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
	route *astral.Route
}

func (e commandDescriptionError) Error() string {
	return "invalid command description for " + strings.Join(e.route.Path(), "->") + ": " + e.route.Description
}

type argDescriptionError struct {
	route *astral.Route
	arg   *astral.Argument
}

func (e argDescriptionError) Error() string {
	return "invalid argument description for " + strings.Join(e.route.Path(), "->") + " arg " + e.arg.Name + ": " + e.arg.Description
}

type argTypeError struct {
	route *astral.Route
	arg   *astral.Argument
}

func (e argTypeError) Error() string {
	return "invalid argument type for " + strings.Join(e.route.Path(), "->") + " arg " + e.arg.Name + ": " + strconv.Itoa(int(e.arg.Type))
}

func NewInteractionHandler(state *state.State, appID discord.AppID) *InteractionHandler {
	return &InteractionHandler{
		state: state,
		appID: appID,
	}
}

type InteractionHandler struct {
	state  *state.State
	appID  discord.AppID
	parent *astral.Route
}

// RegisterCommands registers all sub routes as interaction/slash commands
func (i *InteractionHandler) RegisterCommands(r *astral.Route) ([]discord.Command, error) {
	return i.RegisterGuildCommands(r, discord.NullGuildID)
}

// RegisterGuildCommands registers all sub routes as interaction/slash commands to a guild
func (i *InteractionHandler) RegisterGuildCommands(r *astral.Route, guildID discord.GuildID) ([]discord.Command, error) {
	commands := make([]api.CreateCommandData, 0)

	for _, sub := range r.Children() {
		if !sub.IsExported() {
			continue
		}

		data, err := i.toCommandData(sub)

		if err != nil {
			return nil, err
		}

		commands = append(commands, data)
	}

	if guildID.IsValid() {
		return i.state.BulkOverwriteGuildCommands(i.appID, guildID, commands)
	}

	return i.state.BulkOverwriteCommands(i.appID, commands)
}

func (i *InteractionHandler) toCommandData(r *astral.Route) (api.CreateCommandData, error) {
	data := api.CreateCommandData{
		Name:        r.Name,
		Description: r.Description,
	}

	if r.Description == "" {
		return data, commandDescriptionError{route: r}
	}

	childLen := len(r.Children())
	if childLen > 0 {
		options := make([]discord.CommandOption, childLen)

		x := 0

		for _, route := range r.Children() {
			inputValues, err := i.argsFromRoute(route)

			if err != nil {
				return data, err
			}

			values := make([]discord.CommandOptionValue, len(inputValues))

			for k, value := range inputValues {
				values[k] = value.(discord.CommandOptionValue)
			}

			options[x] = &discord.SubcommandOption{
				OptionName:  route.Name,
				Options:     values,
				Required:    false,
				Description: route.Description,
			}

			x++
		}

		data.Options = options
	} else {
		args, err := i.argsFromRoute(r)

		if err != nil {
			return data, err
		}

		data.Options = args
	}

	return data, nil
}

// RegisterCommand registers a single command, with sub routes as subcommands.
func (i *InteractionHandler) RegisterCommand(r *astral.Route, guildID discord.GuildID) (*discord.Command, error) {
	data, err := i.toCommandData(r)

	if err != nil {
		return nil, err
	}

	if guildID != discord.NullGuildID {
		return i.state.CreateGuildCommand(i.appID, guildID, data)
	}

	return i.state.CreateCommand(i.appID, data)
}

// UpdateCommand registers a single command, with sub routes as subcommands.
func (i *InteractionHandler) UpdateCommand(r *astral.Route, commandID discord.CommandID, guildID discord.GuildID) (*discord.Command, error) {
	data, err := i.toCommandData(r)

	if err != nil {
		return nil, err
	}

	if guildID != discord.NullGuildID {
		return i.state.EditGuildCommand(i.appID, guildID, commandID, data)
	}

	return i.state.EditCommand(i.appID, commandID, data)
}

var (
	commandNameRe = regexp.MustCompile("[^\\w-]")
)

// argsFromRoute takes a route's arguments and translates them into a discord.CommandOption
func (i *InteractionHandler) argsFromRoute(r *astral.Route) ([]discord.CommandOption, error) {
	options := make([]discord.CommandOption, len(r.Arguments))

	for _, arg := range r.Arguments {
		argName := strings.ToLower(commandNameRe.ReplaceAllString(strings.ToLower(arg.Name), ""))

		if arg.Description == "" {
			return nil, argDescriptionError{route: r, arg: arg}
		}

		switch arg.Type {
		case astral.ArgumentTypeInt:
			opt := &discord.IntegerOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}

			//if arg.autocomplete != nil {
			//	opt.Autocomplete = true
			//}

			if len(arg.Choices) > 0 {
				opt.Choices = i.integerChoices(arg)
			}

			options[arg.Index] = opt
		case astral.ArgumentTypeFloat:
			opt := &discord.NumberOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}

			//if arg.autocomplete != nil {
			//	opt.Autocomplete = true
			//}

			if len(arg.Choices) > 0 {
				opt.Choices = i.numberChoices(arg)
			}

			options[arg.Index] = opt
		case astral.ArgumentTypeBool:
			options[arg.Index] = &discord.BooleanOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}
		case astral.ArgumentTypeUserMention:
			options[arg.Index] = &discord.UserOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}
		case astral.ArgumentTypeChannelMention:
			options[arg.Index] = &discord.ChannelOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}
		case astral.ArgumentTypeRole:
			options[arg.Index] = &discord.RoleOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}
		case astral.ArgumentTypeEmoji, astral.ArgumentTypeBasic:
			opt := &discord.StringOption{
				OptionName:  argName,
				Required:    arg.Required,
				Description: arg.Description,
			}

			//if arg.autocomplete != nil {
			//	opt.Autocomplete = true
			//}

			if len(arg.Choices) > 0 {
				opt.Choices = i.stringChoices(arg)
			}

			options[arg.Index] = opt
		default:
			return nil, argTypeError{route: r, arg: arg}
		}
	}

	return options, nil
}

var (
	ErrUnknownOption   = errors.New("unknown option")
	ErrNotAutocomplete = errors.New("option is not registered to autocomplete")
)

// CallAutocomplete calls the autocomplete handler for a route's argument
func (i *InteractionHandler) CallAutocomplete(ctx *astral.Context, interaction *discord.InteractionEvent, options []discord.AutocompleteOption) error {
	opt := focusedOption(options)

	if opt == nil {
		return ErrUnknownOption
	}

	//arg, exists := ctx.Route.Arguments[opt.Name]

	//if !exists {
	//	return ErrUnknownOption
	//}

	//if arg.autocomplete == nil {
	//	return ErrNotAutocomplete
	//}

	var ret []astral.StringChoice
	// ret := arg.autocomplete(ctx, *opt)

	if ret != nil {
		choices := make(api.AutocompleteStringChoices, len(ret))

		for i, choice := range ret {
			choices[i] = discord.StringChoice{
				Name:  choice.Name,
				Value: choice.Value,
			}
		}

		return i.state.RespondInteraction(interaction.ID, interaction.Token, api.InteractionResponse{
			Type: api.AutocompleteResult,
			Data: &api.InteractionResponseData{
				Choices: &choices,
			},
		})
	}

	return nil
}

func focusedOption(options []discord.AutocompleteOption) *discord.AutocompleteOption {
	for _, opt := range options {
		if opt.Focused {
			return &opt
		}
	}

	return nil
}

// Autocomplete registers an autocomplete handler for this argument
//func (i *InteractionHandler) Autocomplete(f AutocompleteHandler) *Argument {
//	a.autocomplete = f
//	return a
//}

func (i *InteractionHandler) integerChoices(a *astral.Argument) []discord.IntegerChoice {
	choices := make([]discord.IntegerChoice, len(a.Choices))

	for i, choice := range a.Choices {
		v, err := strconv.Atoi(choice.Value)

		if err != nil {
			continue
		}

		choices[i] = discord.IntegerChoice{
			Name:  choice.Name,
			Value: v,
		}
	}

	return choices
}

func (i *InteractionHandler) numberChoices(a *astral.Argument) []discord.NumberChoice {
	choices := make([]discord.NumberChoice, len(a.Choices))

	for i, choice := range a.Choices {
		v, err := strconv.ParseFloat(choice.Value, 64)

		if err != nil {
			continue
		}

		choices[i] = discord.NumberChoice{
			Name:  choice.Name,
			Value: v,
		}
	}

	return choices
}

func (i *InteractionHandler) stringChoices(a *astral.Argument) []discord.StringChoice {
	choices := make([]discord.StringChoice, len(a.Choices))

	for i, choice := range a.Choices {
		choices[i] = discord.StringChoice{
			Name:  choice.Name,
			Value: choice.Value,
		}
	}

	return choices
}

// FindInteraction finds a route path from a command interaction
func (i *InteractionHandler) FindInteraction(parentRoute string, options []discord.CommandInteractionOption) *astral.Route {
	children := i.parent.Children()

	if len(children) < 1 {
		return i.parent
	}

	opts := options

	currentRoute := children[parentRoute]

	if currentRoute == nil {
		return nil
	}

	var routeName string

	for opts != nil {
		routeName, opts = recurseOptions(opts)

		if routeName != "" {
			if newRoute, exists := currentRoute.Children()[routeName]; exists {
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
func (i *InteractionHandler) FindAutocomplete(parentRoute string, options []discord.AutocompleteOption) (*astral.Route, []discord.AutocompleteOption) {
	children := i.parent.Children()

	if len(children) < 1 {
		return i.parent, nil
	}

	opts := options

	currentRoute := children[parentRoute]

	if currentRoute == nil {
		return nil, nil
	}

	var routeName string
	var focused bool

	for opts != nil {
		routeName, opts, focused = recurseAutocompleteOptions(opts)

		if routeName != "" {
			if newRoute, exists := currentRoute.Children()[routeName]; exists {
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
