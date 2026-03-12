package event

import "github.com/auroradevllc/astral/v3"

type ServerCreateEvent struct {
	astral.Server
}

type ServerUpdateEvent struct {
	astral.Server
}

type ServerDeleteEvent struct {
	ServerID astral.ServerID
}

type MemberJoinEvent struct {
	ServerID astral.ServerID
	Member   astral.Member
}

type MemberLeaveEvent struct {
	ServerID astral.ServerID
	UserID   astral.UserID
}

type UserUpdateEvent struct {
	astral.User
}

func (*UserUpdateEvent) Username() string {
	//TODO implement me
	panic("implement me")
}

type ChannelCreateEvent struct {
	astral.Channel
}

type ChannelUpdateEvent struct {
	astral.Channel
}

type ChannelDeleteEvent struct {
	astral.Channel
}

type MessageCreateEvent struct {
	astral.Message
}

type MessageUpdateEvent struct {
	astral.Message
}

type MessageDeleteEvent struct {
	ID        astral.MessageID
	ChannelID astral.ChannelID
	ServerID  astral.ServerID
}

type MessageReactionAddEvent struct {
	ServerID  astral.ServerID
	UserID    astral.UserID
	ChannelID astral.ChannelID
	MessageID astral.MessageID
	Emoji     astral.Emoji
}

type MessageReactionRemoveEvent struct {
	ServerID  astral.ServerID
	UserID    astral.UserID
	ChannelID astral.ChannelID
	MessageID astral.MessageID
	Emoji     astral.Emoji
}

type ServerRoleCreateEvent struct {
	ServerID astral.ServerID
	Role     astral.Role
}

type ServerRoleUpdateEvent struct {
	ServerID astral.ServerID
	Role     astral.Role
}

type ServerRoleDeleteEvent struct {
	ServerID astral.ServerID
	RoleID   astral.RoleID
}
