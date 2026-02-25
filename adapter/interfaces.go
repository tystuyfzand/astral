package adapter

type ID string

type ChannelType int

const (
	ChannelTypeUnknown ChannelType = iota
	ChannelTypeText
	ChannelTypeCategory
	ChannelTypeVoice
	ChannelTypeVideo
	ChannelTypeDirect
)

type Client interface {
	// User looks up/retrieves a User interface wrapper
	User(id ID) (User, error)

	// Channel looks up/retrieves a Channel interface wrapper
	Channel(id ID) (Channel, error)

	// Server looks up/retrieves a Server interface wrapper (Guild, Planet, etc)
	Server(id ID) (Server, error)

	// Interface returns the underlying object backing this client
	Interface() any
}

type Message interface {
	Content() string
}

type Server interface {
	ID() ID
	Name() string
	OwnerID() ID
}

type Channel interface {
	ID() ID
	Name() string
	Type() ChannelType

	// Permissions returns channel permissions as a bitmask
	// TODO: Better method of handling this?
	Permissions() uint64
}

type User interface {
	ID() ID
	Name() string
}
