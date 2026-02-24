package discord

import "github.com/diamondburned/arikawa/v3/discord"

func NewUser(u *discord.User) *User {
	return &User{User: u}
}

type User struct {
	*discord.User
}

func (u *User) ID() string {
	return u.User.ID.String()
}

func (u *User) Name() string {
	return u.User.Username
}
