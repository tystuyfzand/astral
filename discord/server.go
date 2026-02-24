package discord

import "github.com/diamondburned/arikawa/v3/discord"

type Server struct {
	*discord.Guild
}

func (s *Server) ID() string {
	return s.Guild.ID.String()
}

func (s *Server) Name() string {
	return s.Guild.Name
}
