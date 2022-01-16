package router

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"strconv"
)

// StringChoice is a basic wrapper for name/value choices
type StringChoice struct {
	Name  string
	Value string
}

// Argument type contains defined arguments, parsed from the command signature
type Argument struct {
	Index       int
	Name        string
	Description string
	Required    bool
	Type        ArgumentType
	Choices     []StringChoice
	Min         interface{}
	Max         interface{}
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
