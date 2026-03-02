package astral

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
