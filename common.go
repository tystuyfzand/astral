package astral

import (
	"io"

	"github.com/auroradevllc/astral/v3/utils/option"
)

type MessageContent struct {
	Content string
	Embeds  []Embed
	Files   []File
}

type EditMessageContent struct {
	Content option.NullableString
	Embeds  *[]Embed
	Files   *[]File
}

type File struct {
	Name   string
	Size   int64
	Reader io.Reader
}

type Embed struct {
}
