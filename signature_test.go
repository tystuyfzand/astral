package astral

import "testing"

func TestParseSignature(t *testing.T) {
	r := New()

	parseSignature(r, "test <string arg> <:emoji arg> <@mention arg> <#channel arg> <intarg int min:1> [optional val int min:1] <floatarg float> <boolarg bool> [optional]")

	if r.ArgumentCount < 8 {
		t.Fatal("Expected 8 arguments")
	}

	if optional, exists := r.Arguments["optional"]; !exists || optional.Required {
		t.Fatal("Expected optional argument to not be required")
	}

	if basic, exists := r.Arguments["string arg"]; !exists || basic.Type != ArgumentTypeBasic {
		t.Fatal("Expected string arg to be type Basic")
	}

	if emoji, exists := r.Arguments["emoji arg"]; !exists || emoji.Type != ArgumentTypeEmoji {
		t.Fatal("Expected emoji arg to be type Emoji")
	}

	if mention, exists := r.Arguments["mention arg"]; !exists || mention.Type != ArgumentTypeUserMention {
		t.Fatal("Expected mention arg to be type UserMention")
	}

	if channel, exists := r.Arguments["channel arg"]; !exists || channel.Type != ArgumentTypeChannelMention {
		t.Fatal("Expected channel arg to be type Channel")
	}

	if intArg, exists := r.Arguments["intarg"]; !exists || intArg.Type != ArgumentTypeInt {
		t.Fatal("Expected intarg to be type int")
	}

	if optionalIntArg, exists := r.Arguments["optional val"]; !exists || optionalIntArg.Type != ArgumentTypeInt {
		t.Fatal("Expected optional val to be type int")
	}

	if floatArg, exists := r.Arguments["floatarg"]; !exists || floatArg.Type != ArgumentTypeFloat {
		t.Fatal("Expected floatarg to be type float")
	}

	if boolArg, exists := r.Arguments["boolarg"]; !exists || boolArg.Type != ArgumentTypeBool {
		t.Fatal("Expected boolarg to be type bool")
	}
}
