package valour

import (
	"github.com/auroradevllc/astral/v3"
	valour "github.com/auroradevllc/valourgo"
)

var _ astral.Server = (*Server)(nil)

type Server struct {
	*valour.Planet
	client valour.Client
}

func (s *Server) ID() astral.ServerID {
	return astral.ServerID(s.Planet.ID.String())
}

func (s *Server) Name() string {
	return s.Planet.Name
}

func (s *Server) IconURL() string {
	return ""
}

func (s *Server) OwnerID() astral.UserID {
	return astral.UserID(s.Planet.OwnerID.String())
}

func (s *Server) Channel(id astral.ChannelID) (astral.Channel, error) {
	ch, err := s.client.Channel(s.Planet.ID, ChannelID(id))

	if err != nil {
		return nil, err
	}

	return NewChannel(s.client, ch), nil
}

func (s *Server) Channels() ([]astral.Channel, error) {
	channels, err := s.client.Channels(s.Planet.ID)

	if err != nil {
		return nil, err
	}

	return mapRefItems(channels, func(in *valour.Channel) astral.Channel {
		return NewChannel(s.client, in)
	}), nil
}

func (s *Server) Role(id astral.RoleID) (astral.Role, error) {
	role, err := s.client.Role(s.Planet.ID, RoleID(id))

	if err != nil {
		return nil, err
	}

	return NewRole(role), nil
}

func (s *Server) Roles() ([]astral.Role, error) {
	roles, err := s.client.Roles(s.Planet.ID)

	if err != nil {
		return nil, err
	}

	return mapRefItems(roles, NewRole), nil
}

func (s *Server) Emoji(id astral.EmojiID) (astral.Emoji, error) {
	panic("implement me")
}

func (s *Server) Emojis() ([]astral.Emoji, error) {
	panic("implement me")
}

func (s *Server) Member(id astral.UserID) (astral.Member, error) {
	member, err := s.client.MemberByUser(s.Planet.ID, UserID(id))

	if err != nil {
		return nil, err
	}

	return
}

func (s *Server) CreateRole(data astral.CreateRoleData) (astral.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Server) DeleteRole(id astral.RoleID) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) AddRole(memberID astral.MemberID, roleID astral.RoleID) error {
	//TODO implement me
	panic("implement me")
}

func (s *Server) RemoveRole(memberID astral.MemberID, roleID astral.RoleID) error {
	//TODO implement me
	panic("implement me")
}

func NewServer(client valour.Client, planet *valour.Planet) *Server {
	return &Server{
		Planet: planet,
		client: client,
	}
}
