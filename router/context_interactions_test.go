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
			Options: []discord.CommandInteractionOption{
				{
					Name: "cool",
				},
			},
		},
	}
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

	r.On("something", nil).On("cool", nil)

	route := r.FindInteraction(testInteractionData2())

	path := route.Path()

	if len(path) < 2 || path[0] != "something" || path[1] != "cool" {
		t.Fatal("Expected path to be something, cool - got:", path)
	}

	t.Log(route.Name)
	t.Log(path)
}
