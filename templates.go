package astral

import (
	"errors"
	"io"
	"io/fs"
)

var ErrNoTemplate = errors.New("no template found")

// Templates is a resolver func that's used to get a template for a specific type
type Templates func(clientType ClientType) (string, error)

// TemplateMap is a simple, easy to use map of type -> string.
// This allows you to inline templates using the Go multi-line strings easily, at the cost of some flexibility.
// You could also use the embed package and embed into string values, then pass them into the map!
func TemplateMap(m map[ClientType]string) Templates {
	return func(clientType ClientType) (string, error) {
		v, ok := m[clientType]

		if !ok {
			return "", ErrNoTemplate
		}

		return v, nil
	}
}

// TemplateFs allows you to use a fs.FS for defining templates.
// This could be an embed.FS or any other implementation, allowing great flexibility.
// All files MUST have the ".hcl" suffix!
func TemplateFs(f fs.FS) Templates {
	return func(clientType ClientType) (string, error) {
		f, err := f.Open(string(clientType) + ".hcl")

		if err != nil {
			return "", err
		}

		defer f.Close()

		b, err := io.ReadAll(f)

		if err != nil {
			return "", err
		}

		return string(b), nil
	}
}
