package astral

// A middleware handler
type MiddlewareFunc func(Handler) Handler
