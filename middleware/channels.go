package middleware

import (
	"errors"

	"github.com/auroradevllc/astral/v3"
)

// Errors
var (
	ErrChannelNotNSFW = errors.New("this command can only be used in an NSFW channel")
	ErrChannelType    = errors.New("channel type does not match expected type")
)

// RequireNSFW requires a message to be sent from an NSFW channel
func RequireNSFW(catch CatchFunc) astral.MiddlewareFunc {
	return func(fn astral.Handler) astral.Handler {
		return func(ctx *astral.Context) {
			if !ctx.Channel.IsNSFW() {
				callCatch(ctx, catch, ErrChannelNotNSFW)
				return
			}
			fn(ctx)
		}
	}
}

// ChannelType requires the specific channel type from the message
func ChannelType(t astral.ChannelType, catch CatchFunc) astral.MiddlewareFunc {
	return func(fn astral.Handler) astral.Handler {
		return func(ctx *astral.Context) {
			if ctx.Channel.Type() != t {
				callCatch(ctx, catch, ErrChannelType)
				return
			}
			fn(ctx)
		}
	}
}
