package router

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"testing"
)

func testInteractionData() []discord.CommandInteractionOption {
	return []discord.CommandInteractionOption{
		{
			Name: "something",
			Options: []discord.CommandInteractionOption{
				{
					Name: "cool",
					Options: []discord.CommandInteractionOption{
						{Name: "key", Value: []byte("\"value\"")},
					},
				},
			},
		},
	}
}

func testInteractionData2() []discord.CommandInteractionOption {
	return []discord.CommandInteractionOption{
		{
			Name: "something",
		},
	}
}

func TestInteractionOptionValue(t *testing.T) {
	val := "\"value\""

	if val[0] == '"' && val[len(val)-1] == '"' {
		val = val[1 : len(val)-1]
	}

	t.Log(val)
}

func TestOptions(t *testing.T) {
	path := []string{"test", "something", "cool"}

	opts := optionsFromPath(path[1:], testInteractionData())

	if opts == nil {
		t.Fatal("Expected options to not be nil")
	}

	t.Log(opts)
}

func TestRoutePath(t *testing.T) {
	r := New()

	r.On("test", nil).On("something", nil)

	route := r.FindInteraction("test", testInteractionData2())

	path := route.Path()

	if len(path) < 2 || path[0] != "test" && path[1] != "something" {
		t.Fatal("Expected path to be something, cool - got:", path)
	}

	t.Log(route.Name)
	t.Log(path)
}

func TestAutocomplete(t *testing.T) {
	data := discord.AutocompleteInteraction{
		Name: "autocomplete",
		Options: []discord.AutocompleteOption{
			{Type: discord.StringOptionType, Name: "test", Focused: true},
		},
	}

	r := New()

	respondTest := func(ctx *Context) {
		ctx.Reply("You chose: " + ctx.Arg("test"))
	}

	autocompleteFill := func(ctx *Context, option discord.AutocompleteOption) []StringChoice {
		choices := []StringChoice{
			{Name: "Test", Value: "test"},
		}

		if option.Value != "" {
			choices = append(choices, StringChoice{
				Name:  option.Value,
				Value: option.Value,
			})
		}

		return choices
	}

	auto := r.On("autocomplete <test>", respondTest).Argument("test", func(arg *Argument) {
		arg.Description = "Test Arg"
	}).Autocomplete("test", autocompleteFill).Export(true).Desc("Autocomplete test")

	auto.On("nested <test>", respondTest).Autocomplete("test", autocompleteFill)

	match, opts := r.FindAutocomplete(data.Name, data.Options)

	if match == nil {
		t.Fatal("Unable to find match")
	}

	t.Log(opts)

	data = discord.AutocompleteInteraction{
		Name: "autocomplete",
		Options: []discord.AutocompleteOption{
			{
				Type:  discord.SubcommandOptionType,
				Name:  "nested",
				Value: "",
				Options: []discord.AutocompleteOption{
					{Type: discord.StringOptionType, Name: "test", Focused: true, Value: "test123"},
				},
			},
		},
	}

	match, opts = r.FindAutocomplete(data.Name, data.Options)

	if match == nil {
		t.Fatal("Unable to find match")
	}

	t.Log(opts)
}
