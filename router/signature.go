package router

import (
	"github.com/diamondburned/arikawa/v3/discord"
	"strings"
)

type ArgumentType int

// DiscordType returns the Discord CommandOptionType for an argument
func (t ArgumentType) DiscordType() discord.CommandOptionType {
	switch t {
	case ArgumentTypeInt:
		return discord.IntegerOption
	case ArgumentTypeBool:
		return discord.BooleanOption
	case ArgumentTypeUserMention:
		return discord.UserOption
	case ArgumentTypeChannelMention:
		return discord.ChannelOption
	default:
		return discord.StringOption
	}
}

const (
	ArgumentTypeBasic ArgumentType = iota
	ArgumentTypeInt
	ArgumentTypeFloat
	ArgumentTypeBool
	ArgumentTypeEmoji
	ArgumentTypeUserMention
	ArgumentTypeChannelMention
)

const (
	argInt   = "int"
	argFloat = "float"
	argBool  = "bool"
)

// parseSignature parses a route's signature
func parseSignature(r *Route, signature string) *Route {
	r.Name = signature
	r.Usage = signature

	if idx := strings.Index(signature, " "); idx != -1 {
		r.Name = signature[0:idx]

		signature = signature[idx+1:]

		// Parse out command arguments, example:
		// test <arg1> [optional arg2]
		// Walk through string, match < and >, [ and ]
		r.Arguments = make(map[string]*Argument)

		str := signature

		var name string
		var f []string
		var index int

		for {
			if len(str) == 0 {
				break
			}

			ch := str[0]

			if ch == '<' || ch == '[' {
				// Scan until closing arrow or end of string
				for i := 1; i < len(str); i++ {
					if (str[i] == '>' || str[i] == ']') && str[i-1] != '\\' {
						name = str[1:i]
						if i+2 < len(str) {
							str = str[i+2:]
						} else {
							str = ""
						}

						required := false

						if ch == '<' {
							required = true

							r.RequiredArgumentCount++
						}

						t := ArgumentTypeBasic

						f = strings.Fields(name)

						if name[0] == ':' {
							t = ArgumentTypeEmoji
							name = name[1:]
						} else if name[0] == '@' {
							t = ArgumentTypeUserMention
							name = name[1:]
						} else if name[0] == '#' {
							t = ArgumentTypeChannelMention
							name = name[1:]
						} else if len(f) > 1 {
							switch f[1] {
							case argInt:
								t = ArgumentTypeInt
								name = f[0]
							case argFloat:
								t = ArgumentTypeFloat
								name = f[0]
							case argBool:
								t = ArgumentTypeBool
								name = f[0]
							}
						}

						r.Arguments[name] = &Argument{
							Index:    index,
							Name:     name,
							Required: required,
							Type:     t,
						}

						index++

						break
					}
				}
			}
		}

		r.ArgumentCount = len(r.Arguments)
	}

	return r
}
