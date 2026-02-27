package arikawa

import (
	"context"

	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/session/shard"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/samber/lo"
)

type ShardedClient struct {
	*ShardManager
}

func (c *ShardedClient) User(id astral.UserID) (astral.User, error) {
	user, err := c.Shard(0).User(UserID(id))

	if err != nil {
		return nil, err
	}

	return NewUser(user), nil
}

func (c *ShardedClient) Server(id astral.ServerID) (astral.Server, error) {
	guildId := GuildID(id)

	currentShard := c.FromGuildID(guildId)

	guild, err := currentShard.Guild(GuildID(id))

	if err != nil {
		return nil, err
	}

	return NewServer(currentShard, guild), nil
}

func (c *ShardedClient) Servers() ([]astral.Server, error) {
	var servers []astral.Server

	c.ForEach(func(s *state.State) {
		guilds, err := s.Guilds()

		if err != nil {
			return
		}

		servers = append(servers, lo.Map(guilds, func(i discord.Guild, _ int) astral.Server {
			return NewServer(s, &i)
		})...)
	})

	return servers, nil
}

func (c *ShardedClient) Interface() any {
	return c.ShardManager
}

func NewShardManager(token string, fn shard.NewShardFunc) (*ShardManager, error) {
	m := &ShardManager{
		fn:    fn,
		hooks: make([]ShardHook, 0),
	}

	var err error

	m.manager, err = shard.NewManager(token, m.newShard)

	if err != nil {
		return nil, err
	}

	return m, nil
}

type ShardHook func(*state.State)

type ShardManager struct {
	manager *shard.Manager

	fn    shard.NewShardFunc
	hooks []ShardHook
}

func (s *ShardManager) newShard(m *shard.Manager, id *gateway.Identifier) (shard.Shard, error) {
	sh, err := s.fn(m, id)

	if err != nil {
		return nil, err
	}

	for _, hook := range s.hooks {
		hook(sh.(*state.State))
	}

	return sh, err
}

func (s *ShardManager) Shard(i int) *state.State {
	sh := s.manager.Shard(i)

	return sh.(*state.State)
}

func (s *ShardManager) FromGuildID(guildID discord.GuildID) *state.State {
	v, _ := s.manager.FromGuildID(guildID)

	return v.(*state.State)
}

// Open opens all gateways handled by this Manager. If an error occurs, Open
// will attempt to close all previously opened gateways before returning.
func (s *ShardManager) Open(ctx context.Context) error {
	return s.manager.Open(ctx)
}

// Close closes all gateways handled by this Manager; it will stop rescaling if
// the manager is currently being rescaled. If an error occurs, Close will
// attempt to close all remaining gateways first, before returning.
func (s *ShardManager) Close() error {
	return s.manager.Close()
}

func (s *ShardManager) AddHook(hook ShardHook) {
	s.hooks = append(s.hooks, hook)

	// Call all hooks after appending
	// This forces existing shards to register as well
	s.ForEach(hook)
}

func (s *ShardManager) ForEach(hook ShardHook) {
	s.manager.ForEach(func(s shard.Shard) {
		hook(s.(*state.State))
	})
}

func (s *ShardManager) NumShards() int {
	return s.manager.NumShards()
}
