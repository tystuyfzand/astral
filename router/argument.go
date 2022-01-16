package router

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"strconv"
)

// Argument type contains defined arguments, parsed from the command signature
type Argument struct {
	Index       int
	Name        string
	Description string
	Required    bool
	Type        ArgumentType
	Options     []string
}

func (a *Argument) choices() interface{} {
	if len(a.Options) == 0 {
		return nil
	}

	switch a.Type {
	case ArgumentTypeInt:
	}

	return nil
}

func (a *Argument) integerChoices() []discord.IntegerChoice {
	choices := make([]discord.IntegerChoice, len(a.Options))

	for i, choice := range a.Options {
		v, err := strconv.Atoi(choice)

		if err != nil {
			continue
		}

		choices[i] = discord.IntegerChoice{
			Name:  choice,
			Value: v,
		}
	}

	return choices
}

func (a *Argument) numberChoices() []discord.NumberChoice {
	choices := make([]discord.NumberChoice, len(a.Options))

	for i, choice := range a.Options {
		v, err := strconv.ParseFloat(choice, 64)

		if err != nil {
			continue
		}

		choices[i] = discord.NumberChoice{
			Name:  choice,
			Value: v,
		}
	}

	return choices
}

func (a *Argument) stringChoices() []discord.StringChoice {
	choices := make([]discord.StringChoice, len(a.Options))

	for i, choice := range a.Options {
		choices[i] = discord.StringChoice{
			Name:  choice,
			Value: choice,
		}
	}

	return choices
}
