package middleware

import (
	"errors"
	"meow.tf/astral/router"
)

// Errors
var (
	ErrOnCooldown     = errors.New("command is on cooldown")
	ErrChannelNotNSFW = errors.New("this command can only be used in an NSFW channel")
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