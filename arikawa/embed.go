package arikawa

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/diamondburned/arikawa/v3/discord"
)

// Embed is a discord.Embed, but with hcl tags to decode into a struct.
// After decoding, mapstructure is used to decode it into a discord.Embed
type Embed struct {
	Title       *string `hcl:"title"`
	Description *string `hcl:"description"`

	URL *discord.URL `hcl:"url"`

	Timestamp *discord.Timestamp `hcl:"timestamp"`

	// Decoded as a string due to hex issues, can be #003471 or 0x003471 as a string
	Color *string `hcl:"color"`

	Footer    *EmbedFooter    `hcl:"footer,block"`
	Image     *EmbedImage     `hcl:"image,block"`
	Thumbnail *EmbedThumbnail `hcl:"thumbnail,block"`
	Video     *EmbedVideo     `hcl:"video,block"`
	Provider  *EmbedProvider  `hcl:"provider,block"`
	Author    *EmbedAuthor    `hcl:"author,block"`
	Fields    []EmbedField    `hcl:"field,block"`
}

type EmbedFooter struct {
	Text *string      `hcl:"text"`
	Icon *discord.URL `hcl:"icon_url"`
}

// EmbedImage is the large image of an embed.
type EmbedImage struct {
	URL    discord.URL `hcl:"url"`
	Height *uint       `hcl:"height"`
	Width  *uint       `hcl:"width"`
}

// EmbedThumbnail is the small image of an embed. It often appears on the right.
type EmbedThumbnail struct {
	URL    discord.URL `hcl:"url"`
	Height *uint       `hcl:"height"`
	Width  *uint       `hcl:"width"`
}

// EmbedVideo is the video of an embed.
type EmbedVideo struct {
	URL    discord.URL `hcl:"url"`
	Height *uint       `hcl:"height"`
	Width  *uint       `hcl:"width"`
}

type EmbedProvider struct {
	Name string      `hcl:"name"`
	URL  discord.URL `hcl:"url"`
}

type EmbedAuthor struct {
	Name string       `hcl:"name"`
	URL  *discord.URL `hcl:"url"`
	Icon *discord.URL `hcl:"icon_url"`
}

type EmbedField struct {
	Name   string `hcl:"name,label"`
	Value  string `hcl:"value"`
	Inline *bool  `hcl:"inline"`
}

func stringToDiscordColor(f reflect.Type, t reflect.Type, data any) (any, error) {
	if f.Kind() == reflect.String && t == reflect.TypeOf(discord.Color(0)) {
		s := data.(string)
		// support both "0x" and "#" prefixes
		if strings.HasPrefix(s, "#") {
			s = "0x" + s[1:]
		}
		val, err := strconv.ParseInt(s, 0, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid color %q: %w", s, err)
		}
		return discord.Color(val), nil
	}

	return data, nil
}
