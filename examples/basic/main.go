package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"log"
	"meow.tf/astral/v2"
	"meow.tf/astral/v2/arguments"
	"meow.tf/astral/v2/middleware"
	"meow.tf/astral/v2/middleware/cooldown"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	flagToken   = flag.String("token", "", "Discord bot token")
	flagPrefix  = flag.String("prefix", "!", "Command prefix")
	flagAppID   = flag.Int64("appID", 0, "App ID for commands")
	flagGuildID = flag.Int64("guildID", 0, "Guild ID for commands")

	route *astral.Route
)

func main() {
	flag.Parse()

	intents := []gateway.Intents{
		gateway.IntentGuilds,
		gateway.IntentGuildMessages,
	}

	s := state.NewWithIntents("Bot "+*flagToken, intents...)

	s.AddHandler(messageCreateHandler(s))
	s.AddHandler(interactionHandler(s))

	route = astral.New()

	ping := route.On("ping", func(ctx *astral.Context) {
		ctx.Reply("pong!")
	}).Desc("Tests ping")

	ping = ping.Export(true)

	ping.On("pong", func(ctx *astral.Context) {
		ctx.Reply("I love ping pong!")
	}).Desc("Tests pong")

	// Test for registering commands with arguments
	route.Group(func(r *astral.Route) {
		r.On("testing <type> <channel> [#discord channel] [message]", func(ctx *astral.Context) {
			ctx.Replyf("Arg1: %s, Arg2: %s", ctx.Arg("type"), ctx.Arg("channel"))
		}).Desc("Testing command")
	})

	// Test for NSFW middleware
	route.Group(func(r *astral.Route) {
		r.Use(middleware.RequireNSFW(middleware.CatchReply("You have to be in an nsfw channel for this!")))

		r.On("nsfw", func(ctx *astral.Context) {
			ctx.Reply("That's LEWD!")
		})
	})

	// Test for cooldown/rate limiting middleware
	route.Group(func(r *astral.Route) {
		reply := middleware.CatchReply("You're doing that too often! SLOW DOWN!")

		r.Use(cooldown.NewWithCatch(2, time.Minute, cooldown.User, reply))

		r.On("test", func(ctx *astral.Context) {
			ctx.Reply("REPLY!")
		})
	})

	// Test for aliasing
	route.Group(func(r *astral.Route) {
		r.On("testalias", func(ctx *astral.Context) {
			ctx.Reply("Called " + ctx.Command)
		}).Alias("alias")
	})

	route.On("nesting", nil).On("level1 <test>", func(ctx *astral.Context) {
		ctx.Reply("Argument: " + ctx.Arg("test"))
	})

	// Test for autocomplete
	route.On("autocomplete <test>", func(ctx *astral.Context) {
		ctx.Reply("You chose: " + ctx.Arg("test"))
	}).Argument("test", func(arg *astral.Argument) {
		arg.Description = "Test Arg"
	}).Autocomplete("test", func(ctx *astral.Context, option discord.AutocompleteOption) []astral.StringChoice {
		choices := []astral.StringChoice{
			{Name: "Test", Value: "test"},
		}

		if option.Value != "" {
			choices = append(choices, astral.StringChoice{
				Name:  option.Value,
				Value: option.Value,
			})
		}

		return choices
	}).Export(true).Desc("Autocomplete test")

	err := s.Open(context.Background())

	if err != nil {
		log.Fatalln("Unable to connect to Discord:", err)
	}

	log.Println("Ready.")

	if *flagGuildID != 0 {
		log.Println("Registering guild commands")

		cmds, err := astral.RegisterGuildCommands(route, s, discord.AppID(*flagAppID), discord.GuildID(*flagGuildID))

		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Done. Commands:", len(cmds))
	}

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-interrupt
}

func messageCreateHandler(s *state.State) func(evt *gateway.MessageCreateEvent) {
	return func(evt *gateway.MessageCreateEvent) {
		str := evt.Content

		prefix := *flagPrefix

		if !strings.HasPrefix(str, prefix) {
			return
		}

		str = strings.TrimPrefix(str, prefix)

		args := arguments.Parse(str)

		match := route.Find(args...)

		if match == nil {
			log.Println("No match for command args", args)
			return
		}

		level := len(match.Path())

		var command string

		if len(args) > 1 {
			command, args = strings.Join(args[:level], " "), args[level:]
		} else {
			command = str
			args = []string{}
		}

		ctx, err := astral.ContextFrom(s, evt, match, args)

		if err != nil {
			log.Println("Unable to create context:", err)
			return
		}

		ctx.Command = command

		go match.Call(ctx)
	}
}

func interactionHandler(s *state.State) func(evt *gateway.InteractionCreateEvent) {
	return func(evt *gateway.InteractionCreateEvent) {
		switch data := evt.Data.(type) {
		case *discord.CommandInteraction:
			b, _ := json.MarshalIndent(data, "", "\t")

			log.Println(string(b))
			// Find root command
			match := route.FindInteraction(data.Name, data.Options)

			if match == nil {
				log.Println("No match for command args")
				return
			}

			ctx, err := astral.ContextFromInteraction(s, evt, match)

			if err != nil {
				log.Println("Unable to create context:", err)
				return
			}

			go match.Call(ctx)
		case *discord.AutocompleteInteraction:
			// Find root command
			match, opts := route.FindAutocomplete(data.Name, data.Options)

			if match == nil {
				log.Println("No match for command args")
				return
			}

			ctx, err := astral.ContextFromInteraction(s, evt, match)

			if err != nil {
				log.Println("Unable to create context:", err)
				return
			}

			log.Println("Calling autocomplete")

			err = match.CallAutocomplete(ctx, opts)

			if err != nil {
				log.Println("Error calling:", err)
			}
		}
	}
}
