package astral

import "errors"

var ErrRoleNotFound = errors.New("role not found")

// FindRoleByName finds a server role by name
func FindRoleByName(s Server, name string) (Role, error) {
	roles, err := s.Roles()

	if err != nil {
		return nil, err
	}

	for _, r := range roles {
		if r.Name() == name {
			return r, nil
		}
	}

	return nil, ErrRoleNotFound
}
