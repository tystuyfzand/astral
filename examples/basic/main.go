package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
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
	flagToken  = flag.String("token", "", "Discord bot token")
	flagPrefix = flag.String("prefix", "!", "Command prefix")

	route *router.Route
)

func main() {
	flag.Parse()

	dg, err := discordgo.New("Bot " + *flagToken)

	if err != nil {
		log.Fatalln("Unable to create discordgo instance:", err)
	}

	dg.AddHandler(messageCreate)

	route = router.New()

	ping := route.On("ping", func(ctx *router.Context) {
		ctx.Reply("pong!")
	})

	ping.On("pong", func(ctx *router.Context) {
		ctx.Reply("I love ping pong!")
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

	err = dg.Open()

	if err != nil {
		log.Fatalln("Unable to connect to Discord:", err)
	}

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-interrupt
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	str := m.Content

	prefix := *flagPrefix

	if !strings.HasPrefix(str, prefix) {
		return
	}

	str = strings.TrimPrefix(str, prefix)

	args := arguments.Parse(str)

	match := route.Find(args...)

	if match == nil {
		return
	}

	idx := strings.Index(str, " ")

	var argString string

	if idx == -1 {
		argString = ""
	} else {
		argString = strings.TrimSpace(str[idx+1:])
	}

	var command string

	if len(args) > 1 {
		command, args = args[0], args[1:]
	} else {
		command = str
		args = []string{}
	}

	ctx, err := router.ContextFrom(s, m, match, command, args, argString)

	if err != nil {
		return
	}

	ctx.Prefix = prefix

	go match.Call(ctx)
}
