package middleware

import (
	"meow.tf/astral/router"
)

// Permission validates the permission level
func Permission(permission int) router.MiddlewareFunc {
	return func(fn router.Handler) router.Handler {
		return func(ctx *router.Context) {
			p, err := ctx.Session.UserChannelPermissions(ctx.User.ID, ctx.Channel.ID)

			if err != nil {
				return
			}

			if p & permission == 0 {
				return // No permission
			}

			fn(ctx)
		}
	}
}
