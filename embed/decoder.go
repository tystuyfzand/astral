package embed

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/auroradevllc/astral/v3"
	"github.com/go-viper/mapstructure/v2"
)

var (
	decoderMap     = make(map[astral.ClientType]Decoder)
	decoderMapLock sync.RWMutex
)

func RegisterDecoder(clientType astral.ClientType, decoder Decoder) {
	decoderMapLock.Lock()
	defer decoderMapLock.Unlock()

	decoderMap[clientType] = decoder
}

func HasDecoder(clientType astral.ClientType) bool {
	decoderMapLock.RLock()
	defer decoderMapLock.RUnlock()

	_, ok := decoderMap[clientType]
	return ok
}

func GetDecoder(clientType astral.ClientType) Decoder {
	decoderMapLock.RLock()
	defer decoderMapLock.RUnlock()

	return decoderMap[clientType]
}

type Decoder interface {
	Output() any
}

func NewSimpleDecoder[V any]() *SimpleDecoder {
	return &SimpleDecoder{
		outStruct: reflect.TypeOf((*V)(nil)).Elem(),
	}
}

type SimpleDecoder struct {
	outStruct reflect.Type
}

// Output MUST return type of *outStruct
func (s *SimpleDecoder) Output() any {
	return reflect.New(s.outStruct).Interface()
}

type DecoderOption func(*mapstructure.DecoderConfig)

func WithDecodeHook(hook ...mapstructure.DecodeHookFunc) DecoderOption {
	return func(d *mapstructure.DecoderConfig) {
		if len(hook) > 1 {
			d.DecodeHook = mapstructure.ComposeDecodeHookFunc(hook...)
		} else {
			d.DecodeHook = hook[0]
		}
	}
}

func WithTagName(name string) DecoderOption {
	return func(d *mapstructure.DecoderConfig) {
		d.TagName = name
	}
}

func NewComplexDecoder[V any, M any](opts ...DecoderOption) *ComplexDecoder {
	config := &mapstructure.DecoderConfig{}

	for _, opt := range opts {
		opt(config)
	}

	return &ComplexDecoder{
		hclStruct:     reflect.TypeOf((*V)(nil)).Elem(),
		outStruct:     reflect.TypeOf((*M)(nil)).Elem(),
		decoderConfig: config,
	}
}

type ComplexDecoder struct {
	hclStruct     reflect.Type
	outStruct     reflect.Type
	decoderConfig *mapstructure.DecoderConfig
}

// Output MUST return type of *hclStruct
func (c *ComplexDecoder) Output() any {
	return reflect.New(c.hclStruct).Interface()
}

// Map MUST return type of *outStruct
func (c *ComplexDecoder) Map(in any) (any, error) {
	out := reflect.New(c.outStruct).Interface()

	// Clone our decoder config and assign an output
	config := *c.decoderConfig
	config.Result = out

	decoder, err := mapstructure.NewDecoder(&config)

	if err != nil {
		return nil, fmt.Errorf("failed creating decoder: %w", err)
	}

	err = decoder.Decode(in)

	if err != nil {
		return nil, fmt.Errorf("failed decoding: %w", err)
	}

	return out, nil
}
