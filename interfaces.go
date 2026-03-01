package astral

type ChannelType int

const (
	ChannelTypeUnknown ChannelType = iota
	ChannelTypeText
	ChannelTypeCategory
	ChannelTypeVoice
	ChannelTypeVideo
	ChannelTypeDirect
)

type ClientType string

type Client interface {
	// User looks up/retrieves a User interface wrapper
	User(id UserID) (User, error)

	// Server looks up/retrieves a Server interface wrapper (Guild, Planet, etc)
	Server(id ServerID) (Server, error)

	// Servers retrieves all servers available to the client
	Servers() ([]Server, error)

	// Interface returns the underlying object backing this client
	Interface() any

	// Type returns the ClientType, used for embeds and identification
	Type() ClientType
}

type Message interface {
	ID() MessageID
	Content() string
}

type Server interface {
	ID() ServerID
	Name() string
	OwnerID() UserID

	// Channel looks up/retrieves a Channel interface wrapper
	Channel(id ChannelID) (Channel, error)

	// Emoji looks up/retrieves either an Emoji (built-in) or custom emoji
	Emoji(id EmojiID) (Emoji, error)

	// Role looks up/retrieves a Role object
	Role(id RoleID) (Role, error)
}

type Channel interface {
	ID() ChannelID
	Name() string
	Type() ChannelType

	// IsNSFW checks if a channel is NSFW
	IsNSFW() bool

	// SendMessage sends a basic text message to the channel
	SendMessage(msg string) (Message, error)

	// Send sends a MessageContent (a message, but with fields for content, embeds, files) to the channel
	Send(r MessageContent) (Message, error)

	// EditMessage edits a message with new content
	EditMessage(id MessageID, msg MessageContent) (Message, error)

	// DeleteMessage removes a message
	DeleteMessage(id MessageID) error
}

type User interface {
	ID() UserID
	Name() string

	// Mention returns a formatted message with a "tag" or mention to @user
	Mention() string
}

type Emoji interface {
	ID() EmojiID
	Name() string
	IsAnimated() bool
}

type Role interface {
	ID() RoleID
	Name() string
}
