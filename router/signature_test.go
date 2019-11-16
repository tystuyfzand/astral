package router

import "testing"

func TestParseSignature(t *testing.T) {
	r := New()

	parseSignature(r, "test <stringarg> <:emojiarg> <@mentionarg> <#channelarg> <intarg int> <floatarg float> <boolarg bool> [optional]")

	if r.ArgumentCount < 8 {
		t.Fatal("Expected 8 arguments")
	}

	if optional, exists := r.Arguments["optional"]; !exists || optional.Required {
		t.Fatal("Expected optional argument to not be required")
	}

	if basic, exists := r.Arguments["stringarg"]; !exists || basic.Type != ArgumentTypeBasic {
		t.Fatal("Expected stringarg to be type Basic")
	}

	if emoji, exists := r.Arguments["emojiarg"]; !exists || emoji.Type != ArgumentTypeEmoji {
		t.Fatal("Expected emojiarg to be type Emoji")
	}

	if mention, exists := r.Arguments["mentionarg"]; !exists || mention.Type != ArgumentTypeUserMention {
		t.Fatal("Expected mentionarg to be type UserMention")
	}

	if channel, exists := r.Arguments["channelarg"]; !exists || channel.Type != ArgumentTypeChannelMention {
		t.Fatal("Expected channelarg to be type Channel")
	}

	if intArg, exists := r.Arguments["intarg"]; !exists || intArg.Type != ArgumentTypeInt {
		t.Fatal("Expected intarg to be type int")
	}

	if floatArg, exists := r.Arguments["floatarg"]; !exists || floatArg.Type != ArgumentTypeFloat {
		t.Fatal("Expected floatarg to be type float")
	}

	if boolArg, exists := r.Arguments["boolarg"]; !exists || boolArg.Type != ArgumentTypeBool {
		t.Fatal("Expected boolarg to be type bool")
	}
}
