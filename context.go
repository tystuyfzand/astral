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

func WithArgumentString(argumentString string) ContextOption {
	return func(ctx *Context) {
		ctx.ParseArguments()
	}
}

type ContextOptions struct {
	Route     *Route
	Client    Client
	Server    Server
	Channel   Channel
	User      User
	Message   Message
	Responder Responder
}

// NewContext creates a new context with the specified options
func NewContext(opts ContextOptions, opt ...ContextOption) *Context {
	c := &Context{
		VariableBag: NewVariableBag(),
		Route:       opts.Route,
		Client:      opts.Client,
		Server:      opts.Server,
		Channel:     opts.Channel,
		Message:     opts.Message,
		User:        opts.User,
		responder:   opts.Responder,
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
