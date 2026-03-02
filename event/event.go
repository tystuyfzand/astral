package event

import "github.com/auroradevllc/astral/v3"

type ServerCreateEvent struct {
	Server astral.Server
}

type ServerUpdateEvent struct {
	Server astral.Server
}

type ServerDeleteEvent struct {
	Server astral.Server
}

type MemberJoinEvent struct {
	Server astral.Server
	Member astral.Member
}

type MemberLeaveEvent struct {
	Server astral.Server
	Member astral.Member
}
