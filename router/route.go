package router

import (
	"regexp"
	"strings"
)

var (
	userMentionRegexp    = regexp.MustCompile("<@!?(\\d+)>")
	channelMentionRegexp = regexp.MustCompile("<#(\\d+)>")
)

type Handler func(*Context)

// Options for FindComplex. Default is just Args in Find.
type FindOpts struct {
	Args      []string
	MatchCase bool
}

// Route type contains information about a route, such as middleware, routes, etc
type Route struct {
	parent     *Route
	handler    Handler
	middleware []MiddlewareFunc
	Routes     map[string]*Route
	aliases    map[string]string

	Name                  string
	Usage                 string
	Description           string
	Arguments             map[string]*Argument
	ArgumentCount         int
	RequiredArgumentCount int
}

// Argument type contains defined arguments, parsed from the command signature
type Argument struct {
	Index    int
	Name     string
	Required bool
	Type     int
}

// Create a new, empty route.
func New() *Route {
	return &Route{
		middleware: make([]MiddlewareFunc, 0),
		Routes:     make(map[string]*Route),
		aliases:    make(map[string]string),
	}
}

// Add a sub route to this route.
func (r *Route) Add(n *Route) *Route {
	r.Routes[n.Name] = n
	return r
}

// Set this route's description
func (r *Route) Desc(description string) *Route {
	r.Description = description
	return r
}

// Add an alias to the parent route for the current route.
func (r *Route) Alias(alias string) *Route {
	if r.parent != nil {
		r.parent.aliases[alias] = r.Name
	}
	return r
}

// Adds a handler for a specific command.
// Signature can be a simple command, or a string like the following:
//  	command <arg1> <arg2> [arg3] [#channel] [@user]
// The library will automatically parse and validate the required arguments.
// <> means an argument will be required, [] says it's optional
// As well as required and optional types, you can use # and @ to signify
// That routes must match a valid user or channel.
func (r *Route) On(signature string, f Handler) *Route {
	rt := New()
	rt.parent = r
	rt.handler = f
	parseSignature(rt, signature)
	r.Routes[rt.Name] = rt.Use(r.middleware...)
	return rt
}

// Create a temporary route to use for registering sub routes.
// All routes will be copied into this route, with middleware applied.
func (r *Route) Group(fn func(*Route)) *Route {
	rt := New()
	rt.Use(r.middleware...)
	fn(rt)

	for _, sub := range rt.Routes {
		r.Add(sub)
	}

	for alias, name := range rt.aliases {
		r.aliases[alias] = name
	}

	return r
}

// Apply middleware to this route. All sub-routes will also inherit this middleware.
func (r *Route) Use(f ...MiddlewareFunc) *Route {
	if r.middleware == nil {
		r.middleware = f
	} else {
		r.middleware = append(r.middleware, f...)
	}

	return r
}

// Find a route by arguments
func (r *Route) Find(args ...string) *Route {
	return r.FindComplex(FindOpts{Args: args})
}

// Find route by options, including args, case sensitive matching, etc
func (r *Route) FindComplex(opts FindOpts) *Route {
	if len(opts.Args) > 0 {
		routeName := opts.Args[0]

		if !opts.MatchCase {
			routeName = strings.ToLower(routeName)
		}

		if alias, ok := r.aliases[routeName]; ok {
			routeName = alias
		}

		if subRoute, ok := r.Routes[routeName]; ok {
			opts.Args = opts.Args[1:]
			return subRoute.FindComplex(opts)
		}
	}

	if r.handler == nil {
		return nil
	}

	return r
}

// Execute a route.
// Handlers are called synchronously.
// Sub-routes will be walked until the stack is empty or a match couldn't be found.
func (r *Route) Call(ctx *Context) error {
	if ctx.ArgumentCount > 0 {
		if subRoute := r.Find(ctx.Arguments[1:]...); subRoute != nil && subRoute != r {
			ctx.Prefix = ctx.Prefix + " " + ctx.Arguments[0]
			ctx.Arguments = ctx.Arguments[1:]
			ctx.ArgumentCount = len(ctx.Arguments)
			return subRoute.Call(ctx)
		}
	}

	ctx.route = r

	if r.ArgumentCount > 0 {
		// Arguments are cached, construct usage
		if err := r.Validate(ctx); err != nil {
			if err == UsageError {
				_, err = ctx.Reply("Usage: " + ctx.Prefix + r.Usage)
			} else {
				_, err = ctx.Reply(err.Error())
			}
			return err
		}
	}

	handler := r.handler

	for _, v := range r.middleware {
		handler = v(handler)
	}

	handler(ctx)

	return nil
}
