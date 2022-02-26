package cooldown

import (
	"github.com/diamondburned/timedmap"
	"golang.org/x/time/rate"
	"meow.tf/astral"
	"meow.tf/astral/middleware"
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

var (
	rateLimiters = timedmap.New()
	cl           *timedmap.Cleaner
)

// Initialize the cleaner for the rateLimiters map
func init() {
	cl = timedmap.NewCleaner(60 * time.Minute)
	cl.AddCleanable(rateLimiters)
	cl.Start()
}

// Handler is a type which wraps the specific functionality of the rate limiter
type Handler struct {
	m         *timedmap.Map
	limit     int
	timeFrame time.Duration
	flags     int
	catch     middleware.CatchFunc
}

// limitorOrNew Finds or creates a new limiter for the specified ctx
func (h *Handler) limiterOrNew(ctx *astral.Context) *rate.Limiter {
	key := limiterKey(ctx, h.flags)

	var limiter *rate.Limiter

	if v, exists := h.m.Get(key); exists {
		limiter = v.Value.(*rate.Limiter)

		if v.ExpiryTime().Before(time.Now().Add(h.timeFrame)) {
			h.m.Extend(key, h.timeFrame)
		}
	} else {
		limiter = rate.NewLimiter(rate.Limit(float64(h.limit)/h.timeFrame.Seconds()), h.limit)

		h.m.Set(key, limiter, h.timeFrame*2)
	}

	return limiter
}

// Middleware function for the handler
func (h *Handler) Middleware(fn astral.Handler) astral.Handler {
	return func(ctx *astral.Context) {
		limiter := h.limiterOrNew(ctx)

		if !limiter.Allow() {
			if h.catch != nil {
				h.catch(ctx)
			}
			return
		}

		fn(ctx)
	}
}

// New Create a new limiter which does nothing when the limit is hit
func New(limit int, timeFrame time.Duration, flags int) astral.MiddlewareFunc {
	return NewWithCatch(limit, timeFrame, flags, nil)
}

// NewWithCatch Create a new limiter which calls "catch" if set when the limit is hit.
func NewWithCatch(limit int, timeFrame time.Duration, flags int, catch middleware.CatchFunc) astral.MiddlewareFunc {
	m := timedmap.New()

	cl.AddCleanable(m)

	h := &Handler{m, limit, timeFrame, flags, catch}

	return h.Middleware
}

// NewShared Create a new limiter shared across the application
func NewShared(limit int, timeFrame time.Duration, flags int) astral.MiddlewareFunc {
	return NewSharedWithCatch(limit, timeFrame, flags, nil)
}

// NewSharedWithCatch Create a new limiter shared across the application which calls "catch" if set when the limit is hit.
func NewSharedWithCatch(limit int, timeFrame time.Duration, flags int, catch middleware.CatchFunc) astral.MiddlewareFunc {
	h := &Handler{rateLimiters, limit, timeFrame, flags, catch}

	return h.Middleware
}

// Construct the rate limiter key from the context given the set of flags
func limiterKey(ctx *astral.Context, flags int) string {
	k := make([]string, 0)

	if flags&Command != 0 {
		k = append(k, "command", ctx.Command)
	}

	if flags&User != 0 {
		k = append(k, "user", ctx.User.ID.String())
	}

	if flags&Channel != 0 {
		k = append(k, "channel", ctx.Channel.ID.String())
	}

	if flags&Server != 0 {
		k = append(k, "guild", ctx.Guild.ID.String())
	}

	if flags&Global != 0 {
		k = append(k, "global")
	}

	return strings.Join(k, "_")
}
