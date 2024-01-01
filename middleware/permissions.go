package middleware

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/samber/lo"
	"meow.tf/astral/v2"
)

// Permission validates the permission level
func Permission(permission discord.Permissions) astral.MiddlewareFunc {
	return func(fn astral.Handler) astral.Handler {
		return func(ctx *astral.Context) {
			member, err := ctx.Session.Member(ctx.Guild.ID, ctx.User.ID)

			if err != nil {
				return
			}

			if ctx.Guild.Roles == nil {
				roles, err := ctx.Session.Roles(ctx.Guild.ID)

				if err != nil {
					return
				}

				ctx.Guild.Roles = roles
			}

			roles := lo.Filter(ctx.Guild.Roles, func(role discord.Role, _ int) bool {
				return lo.Contains(member.RoleIDs, role.ID)
			})

			p := discord.CalcOverrides(*ctx.Guild, *ctx.Channel, *member, roles)

			if !p.Has(permission) {
				return // No permission
			}

			fn(ctx)
		}
	}
}
