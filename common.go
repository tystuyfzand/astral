package astral

import (
	"io"
)

type Response struct {
	Content string
	Embeds  []Embed
	Files   []File
}

type File struct {
	Name   string
	Size   int64
	Reader io.Reader
}

type Embed struct {
}
