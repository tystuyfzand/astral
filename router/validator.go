package router

import (
	"errors"
	"github.com/diamondburned/arikawa/v2/discord"
	"regexp"
	"strconv"
)

var (
	UsageError  = errors.New("usage")
	emojiRegexp = regexp.MustCompile("<(a?):(.+?):(\\d+)>")
)

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
			return errors.New("The " + arg.Name + " argument is required.")
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
	}

	return nil
}

// validateInt checks an integer argument to ensure it's a valid integer
func validateInt(ctx *Context, arg *Argument, argValue string) error {
	_, err := strconv.ParseInt(argValue, 10, 64)

	if err != nil {
		return errors.New(arg.Name + " must be an integer.")
	}

	return nil
}

// validateFloat checks an integer argument to ensure it's a valid float
func validateFloat(ctx *Context, arg *Argument, argValue string) error {
	_, err := strconv.ParseFloat(argValue, 64)

	if err != nil {
		return errors.New(arg.Name + " must be an float.")
	}

	return nil
}

// validateBool checks an integer argument to ensure it's a valid bool
func validateBool(ctx *Context, arg *Argument, argValue string) error {
	_, err := strconv.ParseBool(argValue)

	if err != nil {
		return errors.New(arg.Name + " must be a true/false value.")
	}

	return nil
}

// validateEmoji checks an integer argument to ensure it's a valid bool
func validateEmoji(ctx *Context, arg *Argument, argValue string) error {
	if !emojiRegexp.MatchString(argValue) {
		return errors.New(arg.Name + " must be a valid emoji.")
	}

	return nil
}

// validateUserMention checks a user mention argument to ensure the user exists
func validateUserMention(ctx *Context, arg *Argument, argValue string) error {
	m := userMentionRegexp.FindStringSubmatch(argValue)

	if m == nil {
		return errors.New(arg.Name + " must be a user.")
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
	return errors.New(arg.Name + " must be a user.")
}

// validateChannelMention checks a channel mention argument to ensure the channel exists
func validateChannelMention(ctx *Context, arg *Argument, argValue string) error {
	m := channelMentionRegexp.FindStringSubmatch(argValue)

	if m == nil {
		return errors.New(arg.Name + " must be a channel.")
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
	return errors.New(arg.Name + " must be a channel.")
}
