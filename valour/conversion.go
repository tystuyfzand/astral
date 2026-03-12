package valour

import (
	"github.com/auroradevllc/astral/v3"
	valour "github.com/auroradevllc/valourgo"
)

func PlanetID(id astral.ServerID) valour.PlanetID {
	sf, err := valour.ParseSnowflake[valour.PlanetID](string(id))

	if err != nil {
		return valour.NullPlanetID
	}

	return sf
}

func ChannelID(id astral.ChannelID) valour.ChannelID {
	sf, err := valour.ParseSnowflake[valour.ChannelID](string(id))

	if err != nil {
		return valour.NullChannelID
	}

	return sf
}

func RoleID(id astral.RoleID) valour.RoleID {
	sf, err := valour.ParseSnowflake[valour.RoleID](string(id))

	if err != nil {
		return valour.NullRoleID
	}

	return sf
}
