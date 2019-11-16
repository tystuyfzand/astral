package middleware

import "meow.tf/astral/router"

const (
	ctxPrefix  = "middleware."
	ctxError   = ctxPrefix + "err"
)

// CatchFunc function called if one of the middleware experiences an error
// Can be left as nil
type CatchFunc func(ctx *router.Context)

// CatchReply returns a function that prints the message you pass it
func CatchReply(message string) func(ctx *router.Context) {
	return func(ctx *router.Context) {
		ctx.Reply(message)
	}
}

// callCatch calls a catch function with an error
func callCatch(ctx *router.Context, fn CatchFunc, err error) {
	if fn == nil {
		return
	}
	ctx.Set(ctxError, err)
	fn(ctx)
}