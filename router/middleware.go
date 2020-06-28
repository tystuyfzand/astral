package router

// A middleware handler
type MiddlewareFunc func(Handler) Handler
