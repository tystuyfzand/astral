package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/samber/lo"
)

var (
	_ astral.User   = (*User)(nil)
	_ astral.Member = (*Member)(nil)
)

func NewUser(u *discord.User) *User {
	return &User{User: u}
}

type User struct {
	*discord.User
}

func (u *User) ID() astral.UserID {
	return astral.UserID(u.User.ID.String())
}

func (u *User) Username() string {
	return u.User.Username
}

func (u *User) Name() string {
	return u.DisplayOrUsername()
}

func (u *User) IsBot() bool {
	return u.User.Bot
}

type Role struct {
	*discord.Role
}

func NewRole(r *discord.Role) *Role {
	return &Role{Role: r}
}

func (r *Role) ID() astral.RoleID {
	return astral.RoleID(r.Role.ID.String())
}

func (r *Role) Name() string {
	return r.Role.Name
}

func NewMember(m *discord.Member) *Member {
	return &Member{Member: m}
}

type Member struct {
	*discord.Member
}

func (m *Member) ID() astral.MemberID {
	return astral.MemberID(m.Member.User.ID.String())
}

func (m *Member) UserID() astral.UserID {
	return astral.UserID(m.Member.User.ID.String())
}

func (m *Member) Name() string {
	if m.Member.Nick != "" {
		return m.Member.Nick
	}

	return m.Member.User.Username
}

func (m *Member) Roles() []astral.RoleID {
	return lo.Map(m.Member.RoleIDs, func(id discord.RoleID, _ int) astral.RoleID {
		return astral.RoleID(id.String())
	})
}
