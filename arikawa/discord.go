package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/auroradevllc/astral/v3/embed"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

const Discord astral.ClientType = "discord"

func init() {
	embed.RegisterDecoder(Discord, embed.NewComplexDecoder[Embed, discord.Embed]())
}

func NewClient(state *state.State) *Client {
	return &Client{
		state: state,
	}
}

type Client struct {
	state *state.State
}

func (c *Client) Servers() ([]astral.Server, error) {
	guilds, err := c.state.Guilds()

	if err != nil {
		return nil, err
	}

	servers := make([]astral.Server, len(guilds))

	for i, guild := range guilds {
		servers[i] = NewServer(c.state, &guild)
	}

	return servers, nil
}

func (c *Client) User(id astral.UserID) (astral.User, error) {
	user, err := c.state.User(UserID(id))

	if err != nil {
		return nil, err
	}

	return NewUser(user), nil
}

func (c *Client) Server(id astral.ServerID) (astral.Server, error) {
	guild, err := c.state.Guild(GuildID(id))

	if err != nil {
		return nil, err
	}

	return NewServer(c.state, guild), nil
}

func (c *Client) Interface() any {
	return c.state
}

func (c *Client) Type() astral.ClientType {
	return Discord
}

// ContextFrom creates a new MessageContext from the session and event
func ContextFrom(state *state.State, event *gateway.MessageCreateEvent, r *astral.Route, args []string) (*astral.Context, error) {
	// Find the channel for the event, which doesn't have a built-in discordgo equivalent of .Guild()
	c, err := state.Channel(event.ChannelID)

	if err != nil {
		return nil, err
	}

	var g *discord.Guild

	if c.Type != discord.DirectMessage {
		// Find the guild for that channel. This uses State if enabled.
		g, err = state.Guild(c.GuildID)

		if err != nil {
			return nil, err
		}
	}

	ctx := astral.NewContext(astral.ContextOptions{
		Route:   r,
		Client:  NewClient(state),
		Server:  NewServer(state, g),
		Channel: NewChannel(state, c),
		User:    NewUser(&event.Author),
		Message: NewMessage(&event.Message),
		Responder: &MessageResponder{
			event:   event,
			state:   state,
			channel: c,
			user:    event.Author,
			message: event.Message,
		},
	})

	if err := ctx.ParseArguments(args); err != nil {
		return nil, err
	}

	return ctx, nil
}
