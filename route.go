package astral

import (
	"errors"
	"regexp"
	"strings"
)

var (
	userMentionRegexp    = regexp.MustCompile("<@!?(\\d+)>")
	channelMentionRegexp = regexp.MustCompile("<#(\\d+)>")
	roleMentionRegexp    = regexp.MustCompile("<@&(\\d+)>")
)

// Handler is a command handler.
type Handler func(*Context)

// FindOpts represents options for FindComplex. Default is just Path in Find.
type FindOpts struct {
	Path      []string
	MatchCase bool
}

// Route type contains information about a route, such as middleware, routes, etc
type Route struct {
	parent     *Route
	handler    Handler
	middleware []MiddlewareFunc
	routes     map[string]*Route
	aliases    map[string]string
	export     bool

	Name                  string
	Usage                 string
	Description           string
	Arguments             map[string]*Argument
	ArgumentCount         int
	RequiredArgumentCount int
}

// New creates a new, empty route.
func New() *Route {
	return &Route{
		middleware: make([]MiddlewareFunc, 0),
		routes:     make(map[string]*Route),
		aliases:    make(map[string]string),
	}
}

// Children returns the routes that this parent route has
func (r *Route) Children() map[string]*Route {
	return r.routes
}

// Path returns the route's full path
func (r *Route) Path() []string {
	path := []string{r.Name}

	parent := r.parent

	for parent != nil && parent.Name != "" {
		path = append(path, parent.Name)

		parent = parent.parent
	}

	// Reverse path
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

// Argument is a quick helper that lets you pull a route's argument into a func and modify it.
func (r *Route) Argument(name string, f func(*Argument)) *Route {
	if arg, ok := r.Arguments[name]; ok {
		f(arg)
	} else {
		panic("Unable to find argument " + name)
	}

	return r
}

// Autocomplete is a helper func to pass through autocomplete functions into options
func (r *Route) Autocomplete(name string, f AutocompleteHandler) *Route {
	if arg, ok := r.Arguments[name]; ok {
		arg.autocomplete = f
	} else {
		panic("Unable to find argument " + name)
	}

	return r
}

// Add adds a sub route to this route.
func (r *Route) Add(n *Route) *Route {
	n.parent = r
	r.routes[n.Name] = n
	return r
}

// Desc sets this route's description
func (r *Route) Desc(description string) *Route {
	r.Description = description
	return r
}

// Alias adds an alias to the parent route for the current route.
func (r *Route) Alias(alias string) *Route {
	if r.parent != nil {
		r.parent.aliases[alias] = r.Name
	}
	return r
}

// Export sets the route to be exported to either commands or guild commands.
func (r *Route) Export(export bool) *Route {
	r.export = export
	return r
}

func (r *Route) IsExported() bool {
	return r.export
}

// On adds a handler for a specific command.
// Signature can be a simple command, or a string like the following:
//
//	command <arg1> <arg2> [arg3] [#channel] [@user]
//
// The library will automatically parse and validate the required arguments.
// <> means an argument will be required, [] says it's optional
// As well as required and optional types, you can use # and @ to signify
// That routes must match a valid user or channel.
func (r *Route) On(signature string, f Handler) *Route {
	rt := New()
	rt.parent = r
	rt.handler = f
	rt.export = r.export
	parseSignature(rt, signature)
	r.routes[rt.Name] = rt.Use(r.middleware...)
	return rt
}

// Group creates a temporary route to use for registering sub routes.
// All routes will be copied into this route, with middleware applied.
func (r *Route) Group(fn func(*Route)) *Route {
	rt := New()
	rt.Use(r.middleware...)
	fn(rt)

	for _, sub := range rt.routes {
		r.Add(sub)
	}

	for alias, name := range rt.aliases {
		r.aliases[alias] = name
	}

	return r
}

// Use applies middleware to this route. All sub-routes will also inherit this middleware.
func (r *Route) Use(f ...MiddlewareFunc) *Route {
	if r.middleware == nil {
		r.middleware = f
	} else {
		r.middleware = append(r.middleware, f...)
	}

	return r
}

// Find a route by arguments
func (r *Route) Find(path ...string) *Route {
	return r.FindComplex(FindOpts{Path: path})
}

// FindComplex finds a route by options, including args, case sensitive matching, etc
func (r *Route) FindComplex(opts FindOpts) *Route {
	if len(opts.Path) > 0 {
		routeName := opts.Path[0]

		if !opts.MatchCase {
			routeName = strings.ToLower(routeName)
		}

		if alias, ok := r.aliases[routeName]; ok {
			routeName = alias
		}

		if subRoute, ok := r.routes[routeName]; ok {
			opts.Path = opts.Path[1:]
			return subRoute.FindComplex(opts)
		}
	}

	if r.handler == nil {
		return nil
	}

	return r
}

// Call executes a route.
// Handlers are called synchronously.
// Sub-routes will no longer be recursed automatically, and must be found using Find(...)
func (r *Route) Call(ctx *Context) error {
	ctx.Route = r

	if r.ArgumentCount > 0 {
		// Arguments are cached, construct usage
		if err := r.Validate(ctx); err != nil {
			if errors.Is(err, UsageError) {
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
