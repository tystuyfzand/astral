package astral

type Client interface {
	User(id string) (User, error)
	Channel(id string) (Channel, error)
}

type Message interface {
	Content() string
}

type Server interface {
	ID() string
	Name() string
}

type Channel interface {
	ID() string
	Name() string
}

type User interface {
	ID() string
	Name() string
}
