package discord

import (
	"github.com/auroradevllc/astral/v3/adapter"
	"github.com/diamondburned/arikawa/v3/discord"
)

type Server struct {
	*discord.Guild
}

func (s *Server) ID() adapter.ID {
	return adapter.ID(s.Guild.ID.String())
}

func (s *Server) Name() string {
	return s.Guild.Name
}

func (s *Server) OwnerID() adapter.ID {
	return adapter.ID(s.Guild.OwnerID.String())
}
