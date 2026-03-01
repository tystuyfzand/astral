# Embed System Overview

This library provides a **flexible system for rendering templates** into embeds. It allows a single method to support multiple types of embeds dynamically, without requiring separate manual construction for each type.

Decoders are registered per embed type, and templates are mapped to the corresponding decoder automatically.

## Templates

Maps an **astral.ClientType** to the corresponding template, which can be loaded using the `embed` package, or off a 
filesystem.

```go
e := astral.Embed{
    Templates: func(astral.ClientType) (string, error),
}
```

When you render an embed, the system automatically selects the correct decoder based on its type.

Data can be provided (as variables) using the `Data` field - this can be a map or a struct.

Example usage, using a map:

```go
template := loadTemplate("rich")
embed := astral.Embed{
    Data: map[string]string{
        "key": "value",
    },
    Templates: astral.TemplateMap(map[astral.ClientType]string{
        "discord": template,
    }),
}
```

No manual embed construction is needed — multiple embed types are handled in a single method.

See also `astral.TemplateFs` if you have many templates stored in a directory. You can use `fs.Sub` and `embed.Fs` to 
build them into your application!


### HCL Basics

[HCL]() is a powerful configuration language, but it fits very well into what we'd like to do here.

It was chosen for the following reasons:

- Top level variables which can be mapped directly into values
- Support for `block`s which allow sub-structs, enabling powerful embed configuration
- Easy to understand syntax, easy to provide examples for which can be modified
- Native variable support, allowing direct replacement values

An example of a template for Discord can be found in [arikawa/discord.hcl](../arikawa/discord.hcl), which defines most
fields.

Optional values are defined using pointers, so if you want a string to be optional (or empty) you must define it as a
*string, or the hcl decoder will assume it's required. The mapstructure implementation is smart enough to map `*string`
into `string` with an empty value when using complex decoders.

Example:

```hcl
field = "Value"
num = 12345

block "Test" {
  field = "Value " + variable
}
```

---

## SimpleDecoder

```go
NewSimpleDecoder[V any]
```

Use when the target struct `V` can be decoded directly (already has HCL tags or template-compatible fields). 

Example:

```go
embed.RegisterDecoder("basic", embed.NewSimpleDecoder[MyEmbed]())
```

Workflow:

1. The embed render func will find a template of type `basic`
2. The renderer will attempt to decode it into the struct that was registered (`MyEmbed`)
3. The resulting value, a pointer to the newly created `MyEmbed`, will be returned.

---

## ComplexDecoder

```go
NewComplexDecoder[V any, M any](opt ...DecoderOpt)
```

Use when the target struct `V` **cannot have template tags** (for example, `discord.Embed`). 

The template is first decoded into a **wrapper struct** `M` (with template-compatible fields), then mapped via 
**mapstructure** into `V`.  This allows dynamic handling of multiple embed types.

Example:

```go
embed.RegisterDecoder("discord", embed.NewComplexDecoder(arikawa.Embed, discord.Embed]())
```

Workflow:

1. Template → decoded into wrapper `M` with `hclsimple`
2. `mapstructure` config cloned, targeted at a new value V
3. Wrapper `M` decoded into value of V, returned

