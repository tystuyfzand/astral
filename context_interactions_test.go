package astral

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
			Type: discord.SubcommandOptionType,
		},
	}
}

var _ = ginkgo.Describe("Context Interactions", func() {
	ginkgo.Context("Options", func() {
		ginkgo.It("Should retrieve options from an option slice", func() {
			path := []string{"test", "something", "cool"}

			opts := optionsFromPath(path[1:], testInteractionData())

			Expect(opts).ToNot(BeNil())
		})
	})
	ginkgo.Context("Routes", func() {
		ginkgo.It("Should retrieve the proper path from an interaction", func() {
			r := New()

			r.On("test", nil).On("something", nil)

			route := r.FindInteraction("test", testInteractionData2())

			path := route.Path()

			Expect(path).To(Equal([]string{"test", "something"}))
			Expect(route.Name).To(Equal("something"))
		})
	})
	ginkgo.Context("Autocomplete", func() {

	})
})

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
