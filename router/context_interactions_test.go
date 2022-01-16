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
