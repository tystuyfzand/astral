package arikawa

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/mavolin/dismock/v3/pkg/dismock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func testInteractionData() []discord.CommandInteractionOption {
	return []discord.CommandInteractionOption{
		{
			Name: "something",
			Options: []discord.CommandInteractionOption{
				{
					Name: "cool",
					Options: []discord.CommandInteractionOption{
						{Name: "key", Value: []byte("\"value\"")},
					},
				},
			},
		},
	}
}

func testInteractionData2() []discord.CommandInteractionOption {
	return []discord.CommandInteractionOption{
		{
			Name: "something",
			Type: discord.SubcommandOptionType,
		},
	}
}

var _ = Describe("Context Interactions", func() {
	var (
		m *dismock.Mocker
		s *state.State
		i *InteractionHandler
	)
	BeforeEach(func() {
		m, s = dismock.NewState(GinkgoT())

		i = &InteractionHandler{
			state: s,
		}
	})
	Context("Options", func() {
		It("Should retrieve options from an option slice", func() {
			path := []string{"test", "something", "cool"}

			opts := optionsFromPath(path[1:], testInteractionData())

			Expect(opts).ToNot(BeNil())
		})
	})
	Context("Routes", func() {
		It("Should retrieve the proper path from an interaction", func() {
			r := astral.New()

			r.On("test", nil).On("something", nil)

			route := i.FindInteraction("test", testInteractionData2())

			path := route.Path()

			Expect(path).To(Equal([]string{"test", "something"}))
			Expect(route.Name).To(Equal("something"))
		})
	})
	Context("Autocomplete", func() {

	})
	Context("Context creation", func() {
		var (
			evt = &gateway.InteractionCreateEvent{
				InteractionEvent: discord.InteractionEvent{
					GuildID:   1234,
					ChannelID: 1234,
					Member: &discord.Member{
						User: discord.User{
							ID:            1,
							Username:      "tester",
							Discriminator: "0001",
						},
					},
				},
			}

			r *astral.Route
		)
		BeforeEach(func() {
			r = astral.New()

			m.Channel(discord.Channel{
				ID:      1234,
				Name:    "test",
				Type:    discord.GuildText,
				GuildID: 1234,
			})

			m.Guild(discord.Guild{
				ID:   1234,
				Name: "Test Guild",
			})
		})
		It("Should construct an interaction context from a mock event", func() {
			r = r.On("test", nil)

			ctx, err := ContextFromInteraction(s, evt, r)

			Expect(err).To(BeNil())
			Expect(ctx.Message).ToNot(BeNil())
		})
		It("Should parse arguments using channel discord endpoint", func() {
			ch := discord.Channel{
				ID:   12345,
				Name: "test_argument",
				Type: discord.GuildText,
			}

			m.Channel(ch)

			r = r.On("test <#channel>", nil)

			chJson, _ := ch.ID.MarshalJSON()

			evt.Data = &discord.CommandInteraction{
				Options: []discord.CommandInteractionOption{
					{
						Type:  discord.ChannelOptionType,
						Name:  "channel",
						Value: chJson,
					},
				},
			}

			ctx, err := ContextFromInteraction(s, evt, r)

			Expect(err).To(BeNil())
			Expect(ctx.Message).ToNot(BeNil())
			Expect(ctx.ChannelArg("channel").ID).To(Equal(ch.ID))
		})
		It("Should parse arguments using user discord endpoint", func() {
			u := discord.User{
				ID:            12345,
				Username:      "testing",
				Discriminator: "0001",
			}

			m.User(u)

			r = r.On("test <@user>", nil)

			uJson, _ := u.ID.MarshalJSON()

			evt.Data = &discord.CommandInteraction{
				Options: []discord.CommandInteractionOption{
					{
						Type:  discord.UserOptionType,
						Name:  "user",
						Value: uJson,
					},
				},
			}

			ctx, err := ContextFromInteraction(s, evt, r)

			Expect(err).To(BeNil())
			Expect(ctx.Message).ToNot(BeNil())
			Expect(ctx.UserArg("user").ID).To(Equal(u.ID))
		})
		It("Should parse arguments using role discord endpoint", func() {
			guildId := discord.GuildID(123456)

			role := discord.Role{
				ID:   12345,
				Name: "testing",
			}

			m.Roles(guildId, []discord.Role{role})

			r = r.On("test <&role>", nil)

			rJSON, _ := role.ID.MarshalJSON()

			evt.Data = &discord.CommandInteraction{
				Options: []discord.CommandInteractionOption{
					{
						Type:  discord.RoleOptionType,
						Name:  "role",
						Value: rJSON,
					},
				},
			}

			ctx, err := ContextFromInteraction(s, evt, r)

			Expect(err).To(BeNil())
			Expect(ctx.Message).ToNot(BeNil())
			Expect(ctx.UserArg("role").ID).To(Equal(role.ID))
		})
		It("Should parse multiple arguments simultaneously", func() {
			u := discord.User{
				ID:            12345,
				Username:      "testing",
				Discriminator: "0001",
			}

			m.User(u)
			ch := discord.Channel{
				ID:   12345,
				Name: "test_argument",
				Type: discord.GuildText,
			}

			m.Channel(ch)

			uJson, _ := u.ID.MarshalJSON()
			chJson, _ := ch.ID.MarshalJSON()

			evt.Data = &discord.CommandInteraction{
				Options: []discord.CommandInteractionOption{
					{
						Type:  discord.UserOptionType,
						Name:  "user",
						Value: uJson,
					},
					{
						Type:  discord.ChannelOptionType,
						Name:  "channel",
						Value: chJson,
					},
				},
			}

			r = r.On("test <@user> <#channel> [test]", nil)

			ctx, err := ContextFromInteraction(s, evt, r)

			Expect(err).To(BeNil())
			Expect(ctx.Message).ToNot(BeNil())
			Expect(ctx.UserArg("user").ID).To(Equal(u.ID))
			Expect(ctx.ChannelArg("channel").ID).To(Equal(ch.ID))
		})
		It("Should parse nested arguments properly", func() {

		})
	})
})
