package astral

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

	// Server looks up/retrieves a Server interface wrapper (Guild, Planet, etc)
	Server(id ID) (Server, error)

	// SendMessage sends a basic text message to the specified channel
	SendMessage(channel ID, msg string) (Message, error)

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

	// Channel looks up/retrieves a Channel interface wrapper
	Channel(id ID) (Channel, error)

	// Emoji looks up/retrieves either an Emoji (built-in) or custom emoji
	Emoji(id ID) (Emoji, error)

	// Role looks up/retrieves a Role object
	Role(id ID) (Role, error)
}

type Channel interface {
	ID() ID
	Name() string
	Type() ChannelType
}

type User interface {
	ID() ID
	Name() string
}

type Emoji interface {
	ID() ID
	Name() string
	IsAnimated() bool
}

type Role interface {
	ID() ID
	Name() string
}
