package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
)

func NewUser(u *discord.User) *User {
	return &User{User: u}
}

type User struct {
	*discord.User
}

func (u *User) ID() astral.ID {
	return astral.ID(u.User.ID.String())
}

func (u *User) Name() string {
	return u.User.Username
}

type Role struct {
	*discord.Role
}

func NewRole(r *discord.Role) *Role {
	return &Role{Role: r}
}

func (r *Role) ID() astral.ID {
	return astral.ID(r.Role.ID.String())
}

func (r *Role) Name() string {
	return r.Role.Name
}
