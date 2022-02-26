package astral

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"strconv"
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

// Autocomplete registers an autocomplete handler for this argument
func (a *Argument) Autocomplete(f AutocompleteHandler) *Argument {
	a.autocomplete = f
	return a
}

func (a *Argument) integerChoices() []discord.IntegerChoice {
	choices := make([]discord.IntegerChoice, len(a.Choices))

	for i, choice := range a.Choices {
		v, err := strconv.Atoi(choice.Value)

		if err != nil {
			continue
		}

		choices[i] = discord.IntegerChoice{
			Name:  choice.Name,
			Value: v,
		}
	}

	return choices
}

func (a *Argument) numberChoices() []discord.NumberChoice {
	choices := make([]discord.NumberChoice, len(a.Choices))

	for i, choice := range a.Choices {
		v, err := strconv.ParseFloat(choice.Value, 64)

		if err != nil {
			continue
		}

		choices[i] = discord.NumberChoice{
			Name:  choice.Name,
			Value: v,
		}
	}

	return choices
}

func (a *Argument) stringChoices() []discord.StringChoice {
	choices := make([]discord.StringChoice, len(a.Choices))

	for i, choice := range a.Choices {
		choices[i] = discord.StringChoice{
			Name:  choice.Name,
			Value: choice.Value,
		}
	}

	return choices
}
