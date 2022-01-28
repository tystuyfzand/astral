package router

import (
	"errors"
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
	"regexp"
	"strconv"
)

var (
	UsageError      = errors.New("usage")
	ErrInvalidValue = errors.New("unknown argument value")
	emojiRegexp     = regexp.MustCompile("<(a?):(.+?):(\\d+)>")
)

// InvalidValueError is an error type thrown when a value is invalid/unknown
type InvalidValueError struct {
	Argument string
	Value    string
}

// Error constructs a string for the error with the argument and value
func (i InvalidValueError) Error() string {
	return "unknown argument value for " + i.Argument + ": " + i.Value
}

// Validate checks the context against the Route's defined arguments and ensures all required arguments
// and types are satisfied.
func (r *Route) Validate(ctx *Context) error {
	if ctx.ArgumentCount < r.RequiredArgumentCount {
		return UsageError
	}

	var argValue string
	var err error

	for _, arg := range r.Arguments {
		if ctx.ArgumentCount < arg.Index+1 {
			break
		}

		if !arg.Required {
			continue
		}

		argValue = ctx.Arguments[arg.Index]

		if argValue == "" {
			return fmt.Errorf("The %s argument is required.", arg.Name)
		}

		switch arg.Type {
		case ArgumentTypeInt:
			err = validateInt(ctx, arg, argValue)
		case ArgumentTypeFloat:
			err = validateFloat(ctx, arg, argValue)
		case ArgumentTypeBool:
			err = validateBool(ctx, arg, argValue)
		case ArgumentTypeEmoji:
			err = validateEmoji(ctx, arg, argValue)
		case ArgumentTypeUserMention:
			err = validateUserMention(ctx, arg, argValue)
		case ArgumentTypeChannelMention:
			err = validateChannelMention(ctx, arg, argValue)
		}

		if err != nil {
			return err
		}

		if len(arg.Choices) > 0 {
			// Ensure options contains value
			found := false

			for _, value := range arg.Choices {
				if value.Value == argValue {
					found = true
					break
				}
			}

			if !found {
				return InvalidValueError{Argument: arg.Name, Value: argValue}
			}
		}
	}

	return nil
}

// validateInt checks an integer argument to ensure it's a valid integer
func validateInt(ctx *Context, arg *Argument, argValue string) error {
	v, err := strconv.ParseInt(argValue, 10, 64)

	if err != nil {
		return fmt.Errorf("%s must be an integer.", arg.Name)
	}

	if arg.Min != nil && v < arg.Min.(int64) {
		return fmt.Errorf("%s must be larger than %d.", arg.Name, arg.Min)
	}

	if arg.Max != nil && v < arg.Max.(int64) {
		return fmt.Errorf("%s must be smaller than %d.", arg.Name, arg.Max)
	}

	return nil
}

// validateFloat checks an integer argument to ensure it's a valid float
func validateFloat(ctx *Context, arg *Argument, argValue string) error {
	v, err := strconv.ParseFloat(argValue, 64)

	if err != nil {
		return fmt.Errorf("%s must be a floating point number.", arg.Name)
	}

	if arg.Min != nil && v < arg.Min.(float64) {
		return fmt.Errorf("%s must be larger than %f.", arg.Name, arg.Min)
	}

	if arg.Max != nil && v < arg.Max.(float64) {
		return fmt.Errorf("%s must be smaller than %f.", arg.Name, arg.Max)
	}

	return nil
}

// validateBool checks an integer argument to ensure it's a valid bool
func validateBool(ctx *Context, arg *Argument, argValue string) error {
	_, err := strconv.ParseBool(argValue)

	if err != nil {
		return fmt.Errorf("%s must be a true/false value.", arg.Name)
	}

	return nil
}

// validateEmoji checks an integer argument to ensure it's a valid bool
func validateEmoji(ctx *Context, arg *Argument, argValue string) error {
	if emojiRegexp.MatchString(argValue) {
		return nil
	}

	_, err := emoji.LookupEmoji(argValue)

	if err == nil {
		return nil
	}

	return fmt.Errorf("%s must be a valid emoji.", arg.Name)
}

// validateUserMention checks a user mention argument to ensure the user exists
func validateUserMention(ctx *Context, arg *Argument, argValue string) error {
	m := userMentionRegexp.FindStringSubmatch(argValue)

	if m == nil {
		return fmt.Errorf("%s must be a valid user.", arg.Name)
	}

	sf, err := discord.ParseSnowflake(m[1])

	if err != nil {
		return err
	}

	member, err := ctx.Session.Member(ctx.Guild.ID, discord.UserID(sf))

	if member != nil && err == nil {
		return nil
	}

	// User is not in this guild/doesn't exist.
	return fmt.Errorf("%s must be a valid user.", arg.Name)
}

// validateChannelMention checks a channel mention argument to ensure the channel exists
func validateChannelMention(ctx *Context, arg *Argument, argValue string) error {
	m := channelMentionRegexp.FindStringSubmatch(argValue)

	if m == nil {
		return fmt.Errorf("%s must be a valid channel.", arg.Name)
	}

	sf, err := discord.ParseSnowflake(m[1])

	if err != nil {
		return err
	}

	c, err := ctx.Session.Channel(discord.ChannelID(sf))

	if c != nil && c.GuildID == ctx.Guild.ID {
		return nil
	}

	// Channel does not exist, or is not in this guild.
	return fmt.Errorf("%s must be a valid channel.", arg.Name)
}
