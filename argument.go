package astral

import (
	"github.com/diamondburned/arikawa/v3/discord"
)

type AutocompleteChoice struct {
	Name  string
	Value string
}

// AutocompleteHandler is a handler for autocomplete events.
type AutocompleteHandler func(*Context, discord.AutocompleteOption) []StringChoice

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
	Choices      []StringChoice
	Min          interface{}
	Max          interface{}
}
