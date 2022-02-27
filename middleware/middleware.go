package middleware

import (
	"meow.tf/astral/v2"
)

const (
	ctxPrefix = "middleware."
	ctxError  = ctxPrefix + "err"
)

// CatchFunc function called if one of the middleware experiences an error
// Can be left as nil
type CatchFunc func(ctx *astral.Context)

// CatchReply returns a function that prints the message you pass it
func CatchReply(message string) func(ctx *astral.Context) {
	return func(ctx *astral.Context) {
		ctx.Reply(message)
	}
}

// callCatch calls a catch function with an error
func callCatch(ctx *astral.Context, fn CatchFunc, err error) {
	if fn == nil {
		return
	}
	ctx.Set(ctxError, err)
	fn(ctx)
}

// RecoverFunc is a function called if a handler is recovered
type RecoverFunc func(ctx *astral.Context, v interface{})

// Recoverer is a middleware to catch panics inside calls.
// Usually, this is best to handle to make sure your code is working right, HOWEVER
// this is useful to catch errors and log them instead of fatally erroring.
func Recoverer(rec RecoverFunc) astral.MiddlewareFunc {
	return func(fn astral.Handler) astral.Handler {
		return func(ctx *astral.Context) {
			defer func() {
				if r := recover(); r != nil {
					rec(ctx, r)
				}
			}()

			fn(ctx)
		}
	}
}
