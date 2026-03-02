package event

import (
	"github.com/auroradevllc/astral/v3"
	"github.com/auroradevllc/handler"
)

// Handler is an event handler for all platforms, transforming events into normalized ones
type Handler struct {
	*handler.Handler
	client astral.Client
}

func NewHandler() *Handler {
	return &Handler{
		Handler: handler.New(),
	}
}
