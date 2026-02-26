package astral

import (
	"strconv"
)

type AutocompleteChoice struct {
	Name  string
	Value string
}

// AutocompleteHandler is a handler for autocomplete events.
type AutocompleteHandler func(*Context, AutocompleteOption) []StringChoice

// StringChoice is a basic wrapper for name/value choices
type StringChoice struct {
	Name  string
	Value string
}

// Argument type contains defined arguments, parsed from the command signature
type Argument struct {
	autocomplete AutocompleteHandler
	Index        int
	Name         string
	Description  string
	Required     bool
	Type         ArgumentType
	Choices      []StringChoice // TODO: maybe swap to a generic Choice with int/float/channel/user support?
	Min          interface{}
	Max          interface{}
}

func (a *Argument) HasAutocomplete() bool {
	return a.autocomplete != nil
}

func (a *Argument) CallAutocomplete(ctx *Context, option AutocompleteOption) []StringChoice {
	return a.autocomplete(ctx, option)
}

// AutocompleteOption is a value passed from the client noting the current value, used for autocomplete
// Note that the value here will always be a string and can be converted using the other functions.
type AutocompleteOption struct {
	Name  string
	Value string
}

func (o AutocompleteOption) IntValue() (int64, error) {
	return strconv.ParseInt(o.Value, 10, 64)
}

func (o AutocompleteOption) BoolValue() (bool, error) {
	return strconv.ParseBool(o.Value)
}

func (o AutocompleteOption) FloatValue() (float64, error) {
	return strconv.ParseFloat(o.Value, 64)
}

func (o AutocompleteOption) IDValue() (ID, error) {
	return ID(o.Value), nil
}
