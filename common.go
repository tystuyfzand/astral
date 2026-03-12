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

type CreateRoleData struct {
	Name        string
	Color       Color
	Mentionable bool
}

type File struct {
	Name   string
	Size   int64
	Reader io.Reader
}

// Embed is a special struct which allows us to send multi-platform embeds
// When you provide templates, the client's Type() func is used to resolve a template
// Once a template is resolved, hcl is used to decode a template into an embed, embedding data
type Embed struct {
	// Data can be a struct, map, or slice
	Data any

	// Templates are the list of templates this embed can be rendered as
	// You can provide multiple templates, or just one.
	// When a template is not found, an error will be thrown.
	Templates Templates
}
