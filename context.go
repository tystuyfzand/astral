package astral

import (
	"github.com/auroradevllc/astral/v3/adapter"
)

// Context is the base "context" object.
// It contains all fields that are present on both Messages and Interactions.
type Context struct {
	*VariableBag

	Route          *Route
	Session        adapter.Client
	Server         adapter.Server
	Channel        adapter.Channel
	Message        adapter.Message
	User           adapter.User
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
