package router

import "testing"

func TestRoute_Path(t *testing.T) {
	parent := New()

	p := parent.On("test", nil).On("something", nil).On("deeper", nil).Path()

	if len(p) < 3 {
		t.Fatal("Expected 3 return items")
	}

	if p[0] != "test" {
		t.Fatal("Expected element 0 to be test")
	}

	if p[1] != "something" {
		t.Fatal("Expected element 1 to be something")
	}

	if p[2] != "deeper" {
		t.Fatal("Expected element 2 to be deeper")
	}
}

func TestRoute_Validate(t *testing.T) {
	parent := New()

	parent = parent.On("content", nil)

	r := parent.On("track <type> <name>", nil)
	r.Arguments["type"].Choices = []StringChoice{
		{Name: "Twitch", Value: "twitch"},
	}

	if err := r.Validate(&Context{
		ArgumentCount: 2,
		Arguments:     []string{"twitch", "test"},
	}); err != nil {
		t.Fatal(err)
	}
}
