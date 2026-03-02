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

type Message interface {
	ID() MessageID
	Content() string
}

type Server interface {
	ID() ServerID

	Name() string

	// IconURL returns a URL to the server's icon, or empty if not set
	IconURL() string

	// OwnerID retrieves the owner's user id
	OwnerID() UserID

	// Channel looks up/retrieves a Channel interface wrapper
	Channel(id ChannelID) (Channel, error)

	// Channels retrieves all channels in the server
	Channels() ([]Channel, error)

	// Role looks up/retrieves a Role object
	Role(id RoleID) (Role, error)

	// Roles retrieves all roles in the server
	Roles() ([]Role, error)

	// Emoji looks up/retrieves either an Emoji (built-in) or custom emoji
	Emoji(id EmojiID) (Emoji, error)

	// Emojis lists all emojis in a server
	Emojis() ([]Emoji, error)
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

type Member interface {
	ID() MemberID

	UserID() UserID

	// Name represents a server member specific name override
	Name() string

	// Roles retrieves role ids from a member
	Roles() []RoleID
}
