package arikawa

import (
	"strings"

	"github.com/auroradevllc/astral/v3"
	"github.com/auroradevllc/astral/v3/arguments"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/mavolin/dismock/v3/pkg/dismock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contexts", func() {
	var (
		m *dismock.Mocker
		s *state.State
	)
	BeforeEach(func() {
		m, s = dismock.NewState(GinkgoT())
	})
	Context("Message contexts", func() {
		var (
			evt = &gateway.MessageCreateEvent{
				Message: discord.Message{
					ChannelID: 1234,
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
		It("Should construct a message context from a mock event", func() {
			r = r.On("test", nil)

			ctx, err := ContextFrom(s, evt, r, []string{""})

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

			ctx, err := ContextFrom(s, evt, r, []string{ch.Mention()})

			Expect(err).To(BeNil())
			Expect(ctx.Message).ToNot(BeNil())
			Expect(string(ctx.ChannelArg("channel").ID())).To(Equal(ch.ID.String()))
		})
		It("Should parse arguments using user discord endpoint", func() {
			u := discord.User{
				ID:            12345,
				Username:      "testing",
				Discriminator: "0001",
			}

			m.User(u)

			r = r.On("test <@user>", nil)

			ctx, err := ContextFrom(s, evt, r, []string{u.Mention()})

			Expect(err).To(BeNil())
			Expect(ctx.Message).ToNot(BeNil())
			Expect(string(ctx.UserArg("user").ID())).To(Equal(u.ID.String()))
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

			r = r.On("test <@user> <#channel> [test]", nil)

			ctx, err := ContextFrom(s, evt, r, []string{u.Mention(), ch.Mention(), "test"})

			Expect(err).To(BeNil())
			Expect(ctx.Message).ToNot(BeNil())
			Expect(string(ctx.UserArg("user").ID())).To(Equal(u.ID.String()))
			Expect(string(ctx.ChannelArg("channel").ID())).To(Equal(ch.ID.String()))
		})

		Context("Parsing", func() {
			It("Should correctly parse arguments with spaces/subcommands", func() {
				str := "test arguments with spaces and subcommands"

				r.On("test", nil).On("arguments", func(ctx *astral.Context) {
					// Nothing
				})

				args := arguments.Parse(str)

				match := r.Find(args...)

				Expect(match).ToNot(BeNil())

				level := len(match.Path())

				var command string

				if len(args) > 1 {
					command, args = strings.Join(args[:level], " "), args[level:]
				} else {
					command = str
					args = []string{}
				}

				Expect(command).To(Equal("test arguments"))
				Expect(args).To(Equal([]string{"with", "spaces", "and", "subcommands"}))

				ctx, err := ContextFrom(s, evt, r, args)

				Expect(err).To(BeNil())
				Expect(ctx).ToNot(BeNil())
			})
		})
	})
})
