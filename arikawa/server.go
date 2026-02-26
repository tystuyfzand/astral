package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/state"
)

type Server struct {
	*discord.Guild
	state *state.State
}

func (s *Server) ID() astral.ID {
	return astral.ID(s.Guild.ID.String())
}

func (s *Server) Name() string {
	return s.Guild.Name
}

func (s *Server) OwnerID() astral.ID {
	return astral.ID(s.Guild.OwnerID.String())
}

func (s *Server) Channel(id astral.ID) (astral.Channel, error) {
	channel, err := s.state.Channel(ChannelID(id))

	if err != nil {
		return nil, err
	}

	return NewChannel(s.state, channel), nil
}

func (s *Server) Emoji(id astral.ID) (astral.Emoji, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) Role(id astral.ID) (astral.Role, error) {
	role, err := s.state.Role(s.Guild.ID, RoleID(id))

	if err != nil {
		return nil, err
	}

	return NewRole(role), nil
}

func NewServer(state *state.State, guild *discord.Guild) *Server {
	return &Server{
		Guild: guild,
		state: state,
	}
}
