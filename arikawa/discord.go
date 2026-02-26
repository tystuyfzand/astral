package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

type Client struct {
	state *state.State
}

func (c *Client) User(id astral.ID) (astral.User, error) {
	user, err := c.state.User(UserID(id))

	if err != nil {
		return nil, err
	}

	return NewUser(user), nil
}

func (c *Client) Server(id astral.ID) (astral.Server, error) {
	guild, err := c.state.Guild(GuildID(id))

	if err != nil {
		return nil, err
	}

	return NewServer(c.state, guild), nil
}

func (c *Client) Interface() any {
	return c.state
}

func New(state *state.State) astral.Client {
	return &Client{state: state}
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

	ctx := astral.NewContext(
		r,
		astral.WithClient(&Client{state: state}),
		astral.WithServer(NewServer(state, g)),
		astral.WithChannel(NewChannel(state, c)),
		astral.WithUser(NewUser(&event.Author)),
		astral.WithMessage(NewMessage(&event.Message)),
		astral.WithResponder(&MessageResponder{
			event:   event,
			state:   state,
			channel: c,
			user:    event.Author,
			message: event.Message,
		}),
	)

	if err := ctx.ParseArguments(args); err != nil {
		return nil, err
	}

	return ctx, nil
}
