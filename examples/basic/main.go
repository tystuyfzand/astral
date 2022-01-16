package main

import (
	"context"
	"flag"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"log"
	"meow.tf/astral/arguments"
	"meow.tf/astral/middleware"
	"meow.tf/astral/middleware/cooldown"
	"meow.tf/astral/router"
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

	route *router.Route
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

	route = router.New()

	ping := route.On("ping", func(ctx *router.Context) {
		ctx.Reply("pong!")
	}).Desc("Tests ping")

	ping = ping.Export(true)

	ping.On("pong", func(ctx *router.Context) {
		ctx.Reply("I love ping pong!")
	}).Desc("Tests pong")

	// Test for registering commands with arguments
	route.Group(func(r *router.Route) {
		r.Export(true)

		r.On("testing <type> <channel> [#discord channel] [message]", func(ctx *router.Context) {
			ctx.Replyf("Arg1: %s, Arg2: %s", ctx.Arg("type"), ctx.Arg("channel"))
		}).Desc("Testing command")
	})

	// Test for NSFW middleware
	route.Group(func(r *router.Route) {
		r.Use(middleware.RequireNSFW(middleware.CatchReply("You have to be in an nsfw channel for this!")))

		r.On("nsfw", func(ctx *router.Context) {
			ctx.Reply("That's LEWD!")
		})
	})

	// Test for cooldown/rate limiting middleware
	route.Group(func(r *router.Route) {
		reply := middleware.CatchReply("You're doing that too often! SLOW DOWN!")

		r.Use(cooldown.NewWithCatch(2, time.Minute, cooldown.User, reply))

		r.On("test", func(ctx *router.Context) {
			ctx.Reply("REPLY!")
		})
	})

	// Test for aliasing
	route.Group(func(r *router.Route) {
		r.On("testalias", func(ctx *router.Context) {
			ctx.Reply("Called " + ctx.Command)
		}).Alias("alias")
	})

	err := s.Open(context.Background())

	if err != nil {
		log.Fatalln("Unable to connect to Discord:", err)
	}

	log.Println("Ready.")

	if *flagGuildID != 0 {
		_, err := router.RegisterGuildCommands(route, s, discord.AppID(*flagAppID), discord.GuildID(*flagGuildID))

		if err != nil {
			log.Fatalln(err)
		}
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

		idx := strings.Index(str, " ")

		var argString string

		if idx == -1 {
			argString = ""
		} else {
			argString = strings.TrimSpace(str[idx+1:])
		}

		if len(args) > 1 {
			_, args = args[0], args[1:]
		} else {
			args = []string{}
		}

		ctx, err := router.ContextFrom(s, evt, match, args, argString)

		if err != nil {
			log.Println("Unable to create context:", err)
			return
		}

		go match.Call(ctx)
	}
}

func interactionHandler(s *state.State) func(evt *gateway.InteractionCreateEvent) {
	return func(evt *gateway.InteractionCreateEvent) {
		data, ok := evt.Data.(*discord.CommandInteraction)

		if !ok {
			return
		}

		// Find root command
		match := route.Find(data.Name)

		// Find subcommands if necessary
		match = match.FindInteraction(data.Options)

		if match == nil {
			log.Println("No match for command args")
			return
		}

		ctx, err := router.ContextFromInteraction(s, evt, match)

		if err != nil {
			log.Println("Unable to create context:", err)
			return
		}

		go match.Call(ctx)
	}
}
