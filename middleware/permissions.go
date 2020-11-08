package middleware

import (
	"github.com/diamondburned/arikawa/v2/discord"
	"meow.tf/astral/router"
)

// Permission validates the permission level
func Permission(permission discord.Permissions) router.MiddlewareFunc {
	return func(fn router.Handler) router.Handler {
		return func(ctx *router.Context) {
			member, err := ctx.Session.Member(ctx.Guild.ID, ctx.User.ID)

			if err != nil {
				return
			}

			p := discord.CalcOverwrites(*ctx.Guild, *ctx.Channel, *member)

			if !p.Has(permission) {
				return // No permission
			}

			fn(ctx)
		}
	}
}
