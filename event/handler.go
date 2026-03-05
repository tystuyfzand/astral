package event

import (
	"errors"

	"github.com/auroradevllc/handler"
)

var ErrNoMapper = errors.New("no mapper found")

// Mapper takes in an event and resolves the associated event
type Mapper interface {
	// Map can take in an event and map/resolve any event info needed
	Map(v any) (any, error)
}

type MultiMapper struct {
	mappers []Mapper
}

func (m *MultiMapper) Map(v any) (any, error) {
	for _, mapper := range m.mappers {
		v, err := mapper.Map(v)

		if v == nil || err != nil {
			continue
		}

		return v, nil
	}

	return nil, ErrNoMapper
}

// Handler is an event handler for all platforms, transforming events into normalized ones
type Handler struct {
	*handler.Handler
	mapper Mapper
}

// NewHandler creates a new event handler with associated mapper
func NewHandler(f Mapper) *Handler {
	return &Handler{
		Handler: handler.New(),
		mapper:  f,
	}
}

// Handle takes in events and remaps them
// This has no return as we just ignore mapping errors... is this the right thing to do?
func (h *Handler) Handle(e any) {
	out, err := h.mapper.Map(e)

	if err != nil {
		return
	}

	go h.Call(out)
}
