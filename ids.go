package astral

type ID string

const (
	NullID        = ID("")
	NullServerID  = ServerID("")
	NullMessageID = MessageID("")
	NullChannelID = ChannelID("")
	NullRoleID    = RoleID("")
	NullUserID    = UserID("")
	NullEmojiID   = EmojiID("")
)

func (id ID) IsValid() bool {
	return id != ""
}

type ServerID ID

func (id ServerID) IsValid() bool {
	return ID(id).IsValid()
}

type ChannelID ID

func (id ChannelID) IsValid() bool {
	return ID(id).IsValid()
}

type RoleID ID

func (id RoleID) IsValid() bool {
	return ID(id).IsValid()
}

type MessageID ID

func (id MessageID) IsValid() bool {
	return ID(id).IsValid()
}

type UserID ID

func (id UserID) IsValid() bool {
	return ID(id).IsValid()
}

type MemberID ID

func (id MemberID) IsValid() bool {
	return ID(id).IsValid()
}

type EmojiID ID

func (id EmojiID) IsValid() bool {
	return ID(id).IsValid()
}
