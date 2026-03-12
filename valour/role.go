package valour

import (
	"github.com/auroradevllc/astral/v3"
	valour "github.com/auroradevllc/valourgo"
)

var _ astral.Role = (*Role)(nil)

func NewRole(role *valour.Role) astral.Role {
	return &Role{
		Role: role,
	}
}

type Role struct {
	*valour.Role
}

func (r *Role) ID() astral.RoleID {
	return astral.RoleID(r.Role.ID.String())
}

func (r *Role) Name() string {
	return r.Role.Name
}
