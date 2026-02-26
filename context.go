package astral

import (
	"sync"

	"golang.org/x/sync/errgroup"
)

// Context is the base "context" object.
// It contains all fields that are present on both Messages and Interactions.
type Context struct {
	*VariableBag

	Route          *Route
	Client         Client
	Server         Server
	Channel        Channel
	Message        Message
	User           User
	Prefix         string
	Command        string
	ArgumentString string
	Arguments      map[string]interface{}
	ArgumentCount  int
	responder      Responder
}

type ContextOption func(*Context)

func WithResponder(r Responder) ContextOption {
	return func(ctx *Context) {
		ctx.responder = r
	}
}

func WithClient(c Client) ContextOption {
	return func(ctx *Context) {
		ctx.Client = c
	}
}

func WithServer(ss Server) ContextOption {
	return func(ctx *Context) {
		ctx.Server = ss
	}
}

func WithChannel(c Channel) ContextOption {
	return func(ctx *Context) {
		ctx.Channel = c
	}
}

func WithUser(u User) ContextOption {
	return func(ctx *Context) {
		ctx.User = u
	}
}

func WithMessage(m Message) ContextOption {
	return func(ctx *Context) {
		ctx.Message = m
	}
}

func NewContext(route *Route, opt ...ContextOption) *Context {
	c := &Context{
		VariableBag: NewVariableBag(),
		Route:       route,
	}

	for _, opt := range opt {
		opt(c)
	}

	return c
}

// convertedArg is an internal struct used to pass argument conversion off to a goroutine
type convertedArg struct {
	argument *Argument
	val      interface{}
}

func (c *Context) ParseArguments(args []string) error {
	wg := new(errgroup.Group)

	out := make(map[string]interface{})

	var outLock sync.Mutex

	convertArg := func(arg *Argument) func() error {
		return func() error {
			convertedVal, err := c.ConvertArg(arg, args[arg.Index])

			if err != nil {
				return err
			}

			outLock.Lock()
			out[arg.Name] = convertedVal
			outLock.Unlock()
			return nil
		}
	}

	for _, arg := range c.Route.Arguments {
		if len(args) <= arg.Index {
			continue
		}

		wg.Go(convertArg(arg))
	}

	err := wg.Wait()

	if err != nil {
		return err
	}

	c.Arguments = out

	return nil
}
