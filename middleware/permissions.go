package middleware

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"meow.tf/astral"
)

// Permission validates the permission level
func Permission(permission discord.Permissions) astral.MiddlewareFunc {
	return func(fn astral.Handler) astral.Handler {
		return func(ctx *astral.Context) {
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
