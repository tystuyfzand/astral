package router

import (
	"errors"
	"github.com/diamondburned/arikawa/v3/discord"
	emoji "github.com/tmdvs/Go-Emoji-Utils"
	"strconv"
)

var (
	ErrNoUser    = errors.New("no user found")
	ErrNoChannel = errors.New("no channel found")
)

// Find the specified argument nand return the information and value
func (c *Context) arg(name string) (*Argument, interface{}) {
	if arg, exists := c.route.Arguments[name]; exists {
		return arg, c.Arguments[arg.Name]
	}

	panic("undefined argument " + name)
}

func (c *Context) convertArg(arg *Argument, val interface{}) (interface{}, error) {
	switch arg.Type {
	case ArgumentTypeInt:
		switch v := val.(type) {
		case int:
			return int64(v), nil
		case int32:
			return int64(v), nil
		case int64:
			return v, nil
		}

		v, err := strconv.ParseInt(val.(string), 10, 64)

		if err != nil {
			return nil, err
		}

		return v, nil
	case ArgumentTypeFloat:
		switch v := val.(type) {
		case float32:
			return float64(v), nil
		case float64:
			return v, nil
		}

		return strconv.ParseFloat(val.(string), 64)
	case ArgumentTypeBool:
		if v, ok := val.(bool); ok {
			return v, nil
		}

		return strconv.ParseBool(val.(string))
	case ArgumentTypeUserMention:
		var sf discord.Snowflake
		var ok bool

		if sf, ok = val.(discord.Snowflake); !ok {
			m := userMentionRegexp.FindStringSubmatch(val.(string))

			if m == nil {
				return nil, ErrNoUser
			}

			var err error
			sf, err = discord.ParseSnowflake(m[1])

			if err != nil {
				return nil, err
			}
		}

		return c.Session.User(discord.UserID(sf))
	case ArgumentTypeChannelMention:
		var sf discord.Snowflake
		var ok bool

		if sf, ok = val.(discord.Snowflake); !ok {
			m := channelMentionRegexp.FindStringSubmatch(val.(string))

			if m == nil {
				return nil, ErrNoChannel
			}

			var err error
			sf, err = discord.ParseSnowflake(m[1])

			if err != nil {
				return nil, err
			}
		}

		return c.Session.Channel(discord.ChannelID(sf))
	case ArgumentTypeEmoji:
		m := emojiRegexp.FindStringSubmatch(val.(string))

		if m != nil {
			sf, err := discord.ParseSnowflake(m[3])

			if err != nil {
				return nil, err
			}

			return &discord.Emoji{
				ID:       discord.EmojiID(sf),
				Name:     m[2],
				Animated: m[1] == "a",
			}, nil
		}

		result, err := emoji.LookupEmoji(val.(string))

		if err != nil {
			return nil, err
		}

		return &discord.Emoji{
			Name: result.Value,
		}, nil
	}

	return val, nil
}

// Arg finds and returns a named argument as a string
func (c *Context) Arg(name string) string {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeBasic {
		panic("Trying to use a non-string argument as string")
	}

	return val.(string)
}

// IntArg finds and returns a named int argument
func (c *Context) IntArg(name string) int64 {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeInt {
		panic("Trying to use a non-int argument as int")
	}

	return val.(int64)
}

// FloatArg finds and returns a named float argument
func (c *Context) FloatArg(name string) float64 {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeFloat {
		panic("Trying to use a non-float argument as float")
	}

	return val.(float64)
}

// BoolArg finds and returns a named bool argument
func (c *Context) BoolArg(name string) bool {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeBool {
		panic("Trying to use a non-bool argument as bool")
	}

	return val.(bool)
}

// UserArg finds and returns a named User argument
func (c *Context) UserArg(name string) *discord.User {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeUserMention {
		panic("Trying to use a non-user argument as user")
	}

	return val.(*discord.User)
}

// ChannelArg finds and returns a named Channel argument
func (c *Context) ChannelArg(name string) *discord.Channel {
	return c.ChannelArgType(name, 255)
}

// ChannelArgType finds and returns Channel argument with a specified type
func (c *Context) ChannelArgType(name string, t discord.ChannelType) *discord.Channel {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeChannelMention {
		panic("Trying to use a non-channel argument as channel")
	}

	ch := val.(*discord.Channel)

	if t != 255 && ch.Type != t {
		return nil
	}

	return ch
}

// EmojiArg finds and returns an argument as an emoji
func (c *Context) EmojiArg(name string) *discord.Emoji {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeEmoji {
		panic("Trying to use a non-emoji argument as emoji")
	}

	return val.(*discord.Emoji)
}
