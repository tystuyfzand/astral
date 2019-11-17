package router

import (
	"regexp"
)

var (
	userMentionRegexp    = regexp.MustCompile("<@!?(\\d+)>")
	channelMentionRegexp = regexp.MustCompile("<#(\\d+)>")
)

type Handler func(*Context)

// Route type contains information about a route, such as middleware, routes, etc
type Route struct {
	handler    Handler
	middleware []MiddlewareFunc
	Routes     map[string]*Route

	Name                  string
	Usage                 string
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

func New() *Route {
	return &Route{
		middleware: make([]MiddlewareFunc, 0),
		Routes:     make(map[string]*Route),
	}
}

func (r *Route) Add(n *Route) *Route {
	r.Routes[n.Name] = n
	return r
}

func (r *Route) On(signature string, f Handler) *Route {
	rt := New()
	rt.handler = f
	parseSignature(rt, signature)
	r.Routes[rt.Name] = rt.Use(r.middleware...)
	return rt
}

func (r *Route) Group(fn func(*Route)) *Route {
	rt := New()
	fn(rt)
	for _, sub := range rt.Routes {
		r.Add(sub)
	}
	return r
}

func (r *Route) Use(f ...MiddlewareFunc) *Route {
	if r.middleware == nil {
		r.middleware = f
	} else {
		r.middleware = append(r.middleware, f...)
	}

	return r
}

func (r *Route) Find(args ...string) *Route {
	if len(args) > 0 {
		if subRoute, ok := r.Routes[args[0]]; ok {
			args = args[1:]
			return subRoute.Find(args...)
		}
	}

	if r.handler == nil {
		return nil
	}

	return r
}

func (r *Route) Call(ctx *Context) {
	if ctx.ArgumentCount > 0 {
		if subRoute := r.Find(ctx.Arguments[1:]...); subRoute != nil {
			ctx.Prefix = ctx.Prefix + " " + ctx.Arguments[0]
			ctx.Arguments = ctx.Arguments[1:]
			ctx.ArgumentCount = len(ctx.Arguments)
			subRoute.Call(ctx)
			return
		}
	}

	ctx.route = r

	if r.ArgumentCount > 0 {
		// Arguments are cached, construct usage
		if err := r.Validate(ctx); err != nil {
			if err == UsageError {
				ctx.Reply("Usage: " + ctx.Prefix + r.Usage)
			} else {
				ctx.Reply(err.Error())
			}
			return
		}
	}

	handler := r.handler

	for _, v := range r.middleware {
		handler = v(handler)
	}

	handler(ctx)
}
