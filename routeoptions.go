package astral

type RouteOption func(r *Route)

// WithDescription sets a route's description
func WithDescription(desc string) RouteOption {
	return func(r *Route) {
		r.Description = desc
	}
}

// WithAliases sets a route's alias(es)
func WithAliases(aliases ...string) RouteOption {
	return func(r *Route) {
		for _, alias := range aliases {
			r.Alias(alias)
		}
	}
}

// Export sets a route's export option
func Export() RouteOption {
	return func(r *Route) {
		r.export = true
	}
}
