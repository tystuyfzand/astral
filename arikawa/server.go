package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/samber/lo"
)

type Server struct {
	*discord.Guild
	state *state.State
}

func (s *Server) ID() astral.ServerID {
	return astral.ServerID(s.Guild.ID.String())
}

func (s *Server) Name() string {
	return s.Guild.Name
}

func (s *Server) IconURL() string {
	return s.Guild.IconURL()
}

func (s *Server) OwnerID() astral.UserID {
	return astral.UserID(s.Guild.OwnerID.String())
}

func (s *Server) Channel(id astral.ChannelID) (astral.Channel, error) {
	channel, err := s.state.Channel(ChannelID(id))

	if err != nil {
		return nil, err
	}

	return NewChannel(s.state, channel), nil
}

func (s *Server) Channels() ([]astral.Channel, error) {
	ch, err := s.state.Channels(s.Guild.ID)

	if err != nil {
		return nil, err
	}

	return lo.Map(ch, func(c discord.Channel, _ int) astral.Channel {
		return NewChannel(s.state, &c)
	}), nil
}

func (s *Server) Emoji(id astral.EmojiID) (astral.Emoji, error) {
	emoji, err := s.state.Emoji(s.Guild.ID, EmojiID(id))

	if err != nil {
		return nil, err
	}

	return NewEmoji(emoji), nil
}

func (s *Server) Emojis() ([]astral.Emoji, error) {
	emoji, err := s.state.Emojis(s.Guild.ID)

	if err != nil {
		return nil, err
	}

	return lo.Map(emoji, func(e discord.Emoji, _ int) astral.Emoji {
		return NewEmoji(&e)
	}), nil
}

func (s *Server) Role(id astral.RoleID) (astral.Role, error) {
	role, err := s.state.Role(s.Guild.ID, RoleID(id))

	if err != nil {
		return nil, err
	}

	return NewRole(role), nil
}

func (s *Server) Roles() ([]astral.Role, error) {
	roles, err := s.state.Roles(s.Guild.ID)

	if err != nil {
		return nil, err
	}

	return lo.Map(roles, func(r discord.Role, _ int) astral.Role {
		return NewRole(&r)
	}), nil
}

func NewServer(state *state.State, guild *discord.Guild) *Server {
	return &Server{
		Guild: guild,
		state: state,
	}
}
