package astral

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/mavolin/dismock/v3/pkg/dismock"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
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

var _ = ginkgo.Describe("Context Interactions", func() {
	var (
		m *dismock.Mocker
		s *state.State
	)
	ginkgo.BeforeEach(func() {
		m, s = dismock.NewState(ginkgo.GinkgoT())
	})
	ginkgo.Context("Options", func() {
		ginkgo.It("Should retrieve options from an option slice", func() {
			path := []string{"test", "something", "cool"}

			opts := optionsFromPath(path[1:], testInteractionData())

			Expect(opts).ToNot(BeNil())
		})
	})
	ginkgo.Context("Routes", func() {
		ginkgo.It("Should retrieve the proper path from an interaction", func() {
			r := New()

			r.On("test", nil).On("something", nil)

			route := r.FindInteraction("test", testInteractionData2())

			path := route.Path()

			Expect(path).To(Equal([]string{"test", "something"}))
			Expect(route.Name).To(Equal("something"))
		})
	})
	ginkgo.Context("Autocomplete", func() {

	})
	ginkgo.Context("Context creation", func() {
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

			r *Route
		)
		ginkgo.BeforeEach(func() {
			r = New()

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
		ginkgo.It("Should construct an interaction context from a mock event", func() {
			r = r.On("test", nil)

			ctx, err := ContextFromInteraction(s, evt, r)

			Expect(err).To(BeNil())
			Expect(ctx.Message).ToNot(BeNil())
		})
		ginkgo.It("Should parse arguments using channel discord endpoint", func() {
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
		ginkgo.It("Should parse arguments using user discord endpoint", func() {
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
		ginkgo.It("Should parse arguments using role discord endpoint", func() {
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
		ginkgo.It("Should parse multiple arguments simultaneously", func() {
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
		ginkgo.It("Should parse nested arguments properly", func() {

		})
	})
})

func TestAutocomplete(t *testing.T) {
	data := discord.AutocompleteInteraction{
		Name: "autocomplete",
		Options: []discord.AutocompleteOption{
			{Type: discord.StringOptionType, Name: "test", Focused: true},
		},
	}

	r := New()

	respondTest := func(ctx *Context) {
		ctx.Reply("You chose: " + ctx.Arg("test"))
	}

	autocompleteFill := func(ctx *Context, option discord.AutocompleteOption) []StringChoice {
		choices := []StringChoice{
			{Name: "Test", Value: "test"},
		}

		if option.Value != "" {
			choices = append(choices, StringChoice{
				Name:  option.Value,
				Value: option.Value,
			})
		}

		return choices
	}

	auto := r.On("autocomplete <test>", respondTest).Argument("test", func(arg *Argument) {
		arg.Description = "Test Arg"
	}).Autocomplete("test", autocompleteFill).Export(true).Desc("Autocomplete test")

	auto.On("nested <test>", respondTest).Autocomplete("test", autocompleteFill)

	match, opts := r.FindAutocomplete(data.Name, data.Options)

	if match == nil {
		t.Fatal("Unable to find match")
	}

	t.Log(opts)

	data = discord.AutocompleteInteraction{
		Name: "autocomplete",
		Options: []discord.AutocompleteOption{
			{
				Type:  discord.SubcommandOptionType,
				Name:  "nested",
				Value: "",
				Options: []discord.AutocompleteOption{
					{Type: discord.StringOptionType, Name: "test", Focused: true, Value: "test123"},
				},
			},
		},
	}

	match, opts = r.FindAutocomplete(data.Name, data.Options)

	if match == nil {
		t.Fatal("Unable to find match")
	}

	t.Log(opts)
}
