package middleware

import (
	"errors"
	"github.com/diamondburned/arikawa/v3/discord"
	"meow.tf/astral/router"
)

// Errors
var (
	ErrChannelNotNSFW = errors.New("this command can only be used in an NSFW channel")
	ErrChannelType    = errors.New("channel type does not match expected type")
)

// RequireNSFW requires a message to be sent from an NSFW channel
func RequireNSFW(catch CatchFunc) router.MiddlewareFunc {
	return func(fn router.Handler) router.Handler {
		return func(ctx *router.Context) {
			if !ctx.Channel.NSFW {
				callCatch(ctx, catch, ErrChannelNotNSFW)
				return
			}
			fn(ctx)
		}
	}
}

// ChannelType requires the specific channel type from the message
func ChannelType(t discord.ChannelType, catch CatchFunc) router.MiddlewareFunc {
	return func(fn router.Handler) router.Handler {
		return func(ctx *router.Context) {
			if ctx.Channel.Type != t {
				callCatch(ctx, catch, ErrChannelType)
				return
			}
			fn(ctx)
		}
	}
}
