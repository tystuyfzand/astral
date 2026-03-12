package valour

import (
	"github.com/auroradevllc/astral/v3"
	valour "github.com/auroradevllc/valourgo"
)

var _ astral.Member = (*Member)(nil)

type Member struct {
	*valour.Member
}

func (m *Member) ID() astral.MemberID {
	return astral.MemberID(m.Member.ID.String())
}

func (m *Member) UserID() astral.UserID {
	return astral.UserID(m.Member.UserID.String())
}

func (m *Member) Name() string {
	if m.Member.Nickname != nil {
		return *m.Member.Nickname
	}

	return m.User.Name
}

func (m *Member) Roles() []astral.RoleID {
	return nil
}
