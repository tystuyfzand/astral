package embed

import (
	"errors"
	"fmt"

	"github.com/auroradevllc/astral/v3"
	"github.com/go-viper/mapstructure/v2"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

var (
	ErrNoDecoder = errors.New("no decoder registered for type")
)

func Render[V any](clientType astral.ClientType, e astral.Embed) (*V, error) {
	decoder := GetDecoder(clientType)

	if decoder == nil {
		return nil, ErrNoDecoder
	}

	tpl, err := e.Templates(clientType)

	if err != nil {
		return nil, fmt.Errorf("unable to read template for %s: %w", clientType, err)
	}

	fields, err := extractFields(e.Data)

	if err != nil {
		return nil, err
	}

	ctx := &hcl.EvalContext{
		Variables: fields,
	}

	// reflect.New().Interface() -> new output variable
	o := decoder.Output()

	if err := hclsimple.Decode(string(clientType)+".hcl", []byte(tpl), ctx, o); err != nil {
		return nil, err
	}

	// If the value is an intermediate mapping (arikawa.Embed -> discord.Embed), map the value
	if c, ok := decoder.(*ComplexDecoder); ok {
		o, err = c.Map(o)

		if err != nil {
			return nil, err
		}
	}

	return o.(*V), nil
}

func extractFields(obj any) (map[string]cty.Value, error) {
	if obj == nil {
		return nil, nil
	}

	var m map[string]any

	if err := mapstructure.Decode(obj, &m); err != nil {
		return nil, err
	}

	out := make(map[string]cty.Value)

	for k, v := range m {
		ty, err := gocty.ImpliedType(v)

		if err != nil {
			return nil, err
		}

		val, err := gocty.ToCtyValue(v, ty)

		if err != nil {
			return nil, err
		}

		out[k] = val
	}

	return out, nil
}
