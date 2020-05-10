package cooldown

import (
	"github.com/diamondburned/timedmap"
	"golang.org/x/time/rate"
	"meow.tf/astral/middleware"
	"meow.tf/astral/router"
	"strings"
	"time"
)

const (
	Command = 1 << iota
	User
	Channel
	Server
	Global
)

var rateLimiters = timedmap.New()

// Initialize the cleaner for the rateLimiters map
func init() {
	cl := timedmap.NewCleaner(60 * time.Minute)
	cl.AddCleanable(rateLimiters)
	cl.Start()
}

// Create a new limiter which does nothing when the limit is hit
func New(limit int, timeFrame time.Duration, flags int) router.MiddlewareFunc {
	return NewWithCatch(limit, timeFrame, flags, nil)
}

// Create a new limiter which calls "catch" if set when the limit is hit.
func NewWithCatch(limit int, timeFrame time.Duration, flags int, catch middleware.CatchFunc) router.MiddlewareFunc {
	return func(fn router.Handler) router.Handler {
		return func(ctx *router.Context) {
			limiter := limiterOrNew(limiterKey(ctx, flags), limit, timeFrame)

			if !limiter.Allow() {
				if catch != nil {
					catch(ctx)
				}
				return
			}

			fn(ctx)
		}
	}
}

// Find or create a new rate.Limiter, extending the time on the map if necessary.
func limiterOrNew(key string, limit int, timeFrame time.Duration) *rate.Limiter {
	var limiter *rate.Limiter

	if v := rateLimiters.GetValue(key); v != nil {
		limiter = v.(*rate.Limiter)

		if exp, exists := rateLimiters.GetExpires(key); exists && exp.Before(time.Now().Add(timeFrame)) {
			rateLimiters.Extend(key, timeFrame)
		}
	} else {
		limiter = rate.NewLimiter(rate.Limit(float64(limit)/timeFrame.Seconds()), limit)

		rateLimiters.Set(key, limiter, timeFrame*2)
	}

	return limiter
}

// Construct the rate limiter key from the context given the set of flags
func limiterKey(ctx *router.Context, flags int) string {
	k := make([]string, 0)

	if flags&Command != 0 {
		k = append(k, ctx.Command)
	}

	if flags&User != 0 {
		k = append(k, ctx.User.ID)
	}

	if flags&Channel != 0 {
		k = append(k, ctx.Channel.ID)
	}

	if flags&Server != 0 {
		k = append(k, ctx.Guild.ID)
	}

	if flags&Global != 0 {
		k = append(k, "global")
	}

	return strings.Join(k, "_")
}
