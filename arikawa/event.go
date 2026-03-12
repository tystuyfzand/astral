package arikawa

import (
	"fmt"

	"github.com/auroradevllc/astral/v3"
	"github.com/auroradevllc/astral/v3/event"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

type Mapper struct {
	state *state.State
}

func mapMemberJoin(e *gateway.GuildMemberAddEvent) *event.MemberJoinEvent {
	return &event.MemberJoinEvent{
		ServerID: astral.ServerID(e.GuildID.String()),
		Member:   NewMember(&e.Member),
	}
}

func mapMemberLeave(e *gateway.GuildMemberRemoveEvent) *event.MemberLeaveEvent {
	return &event.MemberLeaveEvent{
		ServerID: astral.ServerID(e.GuildID.String()),
		UserID:   astral.UserID(e.User.ID.String()),
	}
}

func mapUserUpdate(e *gateway.UserUpdateEvent) *event.UserUpdateEvent {
	return &event.UserUpdateEvent{
		User: NewUser(&e.User),
	}
}

func mapGuildUpdate(s *state.State, e *gateway.GuildUpdateEvent) *event.ServerUpdateEvent {
	return &event.ServerUpdateEvent{
		Server: NewServer(s, &e.Guild),
	}
}

func mapGuildDelete(e *gateway.GuildDeleteEvent) *event.ServerDeleteEvent {
	return &event.ServerDeleteEvent{
		ServerID: astral.ServerID(e.ID.String()),
	}
}

func mapChannelCreate(s *state.State, e *gateway.ChannelCreateEvent) *event.ChannelCreateEvent {
	return &event.ChannelCreateEvent{
		Channel: NewChannel(s, &e.Channel),
	}
}

func mapChannelUpdate(s *state.State, e *gateway.ChannelUpdateEvent) *event.ChannelUpdateEvent {
	return &event.ChannelUpdateEvent{
		Channel: NewChannel(s, &e.Channel),
	}
}

func mapChannelDelete(s *state.State, e *gateway.ChannelDeleteEvent) *event.ChannelDeleteEvent {
	return &event.ChannelDeleteEvent{
		Channel: NewChannel(s, &e.Channel),
	}
}

func mapMessageCreate(e *gateway.MessageCreateEvent) *event.MessageCreateEvent {
	return &event.MessageCreateEvent{
		Message: NewMessage(&e.Message),
	}
}

func mapMessageUpdate(e *gateway.MessageUpdateEvent) *event.MessageUpdateEvent {
	return &event.MessageUpdateEvent{
		Message: NewMessage(&e.Message),
	}
}

func mapMessageDelete(e *gateway.MessageDeleteEvent) *event.MessageDeleteEvent {
	return &event.MessageDeleteEvent{
		ID:        astral.MessageID(e.ID.String()),
		ChannelID: astral.ChannelID(e.ChannelID.String()),
		ServerID:  astral.ServerID(e.GuildID.String()),
	}
}

func mapMessageReactionAdd(e *gateway.MessageReactionAddEvent) *event.MessageReactionAddEvent {
	return &event.MessageReactionAddEvent{
		ServerID:  astral.ServerID(e.GuildID.String()),
		UserID:    astral.UserID(e.UserID.String()),
		ChannelID: astral.ChannelID(e.ChannelID.String()),
		MessageID: astral.MessageID(e.MessageID.String()),
		Emoji:     NewEmoji(&e.Emoji),
	}
}

func mapMessageReactionRemove(e *gateway.MessageReactionRemoveEvent) *event.MessageReactionRemoveEvent {
	return &event.MessageReactionRemoveEvent{
		ServerID:  astral.ServerID(e.GuildID.String()),
		UserID:    astral.UserID(e.UserID.String()),
		ChannelID: astral.ChannelID(e.ChannelID.String()),
		MessageID: astral.MessageID(e.MessageID.String()),
		Emoji:     NewEmoji(&e.Emoji),
	}
}

func mapGuildRoleCreate(e *gateway.GuildRoleCreateEvent) *event.ServerRoleCreateEvent {
	return &event.ServerRoleCreateEvent{
		ServerID: astral.ServerID(e.GuildID.String()),
		Role:     NewRole(&e.Role),
	}
}

func mapGuildRoleUpdate(e *gateway.GuildRoleUpdateEvent) *event.ServerRoleUpdateEvent {
	return &event.ServerRoleUpdateEvent{
		ServerID: astral.ServerID(e.GuildID.String()),
		Role:     NewRole(&e.Role),
	}
}

func mapGuildRoleDelete(e *gateway.GuildRoleDeleteEvent) *event.ServerRoleDeleteEvent {
	return &event.ServerRoleDeleteEvent{
		ServerID: astral.ServerID(e.GuildID.String()),
		RoleID:   astral.RoleID(e.RoleID.String()),
	}
}

func NewEventMapper(v any) *EventMapper {
	return &EventMapper{
		client: v,
	}
}

type EventMapper struct {
	client any
}

func (m *EventMapper) Map(v any) (any, error) {
	switch e := v.(type) {
	case *gateway.GuildMemberAddEvent:
		return mapMemberJoin(e), nil
	case *gateway.GuildMemberRemoveEvent:
		return mapMemberLeave(e), nil
	case *gateway.UserUpdateEvent:
		return mapUserUpdate(e), nil
	case *gateway.GuildUpdateEvent:
		return mapGuildUpdate(stateFrom(m.client, e.Guild.ID), e), nil
	case *gateway.GuildDeleteEvent:
		return mapGuildDelete(e), nil
	case *gateway.ChannelCreateEvent:
		return mapChannelCreate(stateFrom(m.client, e.GuildID), e), nil
	case *gateway.ChannelUpdateEvent:
		return mapChannelUpdate(stateFrom(m.client, e.GuildID), e), nil
	case *gateway.ChannelDeleteEvent:
		return mapChannelDelete(stateFrom(m.client, e.GuildID), e), nil
	case *gateway.MessageCreateEvent:
		return mapMessageCreate(e), nil
	case *gateway.MessageUpdateEvent:
		return mapMessageUpdate(e), nil
	case *gateway.MessageDeleteEvent:
		return mapMessageDelete(e), nil
	case *gateway.MessageReactionAddEvent:
		return mapMessageReactionAdd(e), nil
	case *gateway.MessageReactionRemoveEvent:
		return mapMessageReactionRemove(e), nil
	case *gateway.GuildRoleCreateEvent:
		return mapGuildRoleCreate(e), nil
	case *gateway.GuildRoleUpdateEvent:
		return mapGuildRoleUpdate(e), nil
	case *gateway.GuildRoleDeleteEvent:
		return mapGuildRoleDelete(e), nil
	}

	return nil, fmt.Errorf("unknown event type: %T", v)
}

func stateFrom(i any, id discord.GuildID) *state.State {
	if sh, ok := i.(*ShardManager); ok {
		if id.IsValid() {
			return sh.FromGuildID(id)
		}

		return sh.Shard(0)
	}

	return i.(*state.State)
}
