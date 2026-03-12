package valour

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/auroradevllc/valourgo"
)

const Valour = astral.ClientType("valour")

type Client struct {
	client valour.Client
}

func (c *Client) User(id astral.UserID) (astral.User, error) {
	return nil, nil
}

func (c *Client) Server(id astral.ServerID) (astral.Server, error) {
	planet, err := c.client.Planet(PlanetID(id))

	if err != nil {
		return nil, err
	}

	return NewServer(c.client, planet), nil
}

func (c *Client) Servers() ([]astral.Server, error) {
	planets, err := c.client.Planets()

	if err != nil {
		return nil, err
	}

	return mapRefItems(planets, func(in *valour.Planet) astral.Server {
		return NewServer(c.client, in)
	}), nil
}

func (c *Client) Interface() any {
	return c.client
}

func (c *Client) Type() astral.ClientType {
	return Valour
}

var _ astral.Client = (*Client)(nil)

func NewClient(client valour.Client) *Client {
	return &Client{client: client}
}

func mapItems[V any, O any](in []V, fn func(in V) O) []O {
	out := make([]O, len(in))

	for i := range in {
		out[i] = fn(in[i])
	}

	return out
}

func mapRefItems[V any, O any](in []V, fn func(in *V) O) []O {
	out := make([]O, len(in))

	for i := range in {
		out[i] = fn(&in[i])
	}

	return out
}
