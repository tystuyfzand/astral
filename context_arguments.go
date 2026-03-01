package astral

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrNoUser    = errors.New("no user found")
	ErrNoChannel = errors.New("no channel found")
	ErrNoRole    = errors.New("no role found")
)

// Find the specified argument nand return the information and value
func (c *Context) arg(name string) (*Argument, interface{}) {
	if arg, exists := c.Route.Arguments[name]; exists {
		return arg, c.Arguments[arg.Name]
	}

	panic("undefined argument " + name)
}

func (c *Context) ConvertArg(arg *Argument, val any) (any, error) {
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

		v, err := strconv.ParseInt(valToString(val), 10, 64)

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

		return strconv.ParseFloat(valToString(val), 64)
	case ArgumentTypeBool:
		if v, ok := val.(bool); ok {
			return v, nil
		}

		return strconv.ParseBool(valToString(val))
	case ArgumentTypeUserMention:
		// Match Discord style <@ID> mentions
		m := userMentionRegexp.FindStringSubmatch(valToString(val))

		if m != nil {
			val = m[1]
		}

		// TODO: Values can be int64s/etc

		return c.Client.User(UserID(valToString(val)))
	case ArgumentTypeChannelMention:
		m := channelMentionRegexp.FindStringSubmatch(valToString(val))

		if m != nil {
			val = m[1]
		}

		// TODO: Values can be int64s/etc

		return c.Server.Channel(ChannelID(valToString(val)))
	case ArgumentTypeEmoji:
		return c.Server.Emoji(EmojiID(valToString(val)))
	case ArgumentTypeRole:
		return c.Server.Role(RoleID(valToString(val)))
	}

	return val, nil
}

type Stringable interface {
	String() string
}

func valToString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}

	if s, ok := v.(Stringable); ok {
		return s.String()
	}

	return fmt.Sprintf("%v", v)
}

// Arg finds and returns a named argument as a string
func (c *Context) Arg(name string) string {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeBasic {
		panic("Trying to use a non-string argument as string")
	}

	if val == nil {
		return ""
	}

	return val.(string)
}

// IntArg finds and returns a named int argument
func (c *Context) IntArg(name string) int64 {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeInt {
		panic("Trying to use a non-int argument as int")
	}

	if val == nil {
		return 0
	}

	return val.(int64)
}

// FloatArg finds and returns a named float argument
func (c *Context) FloatArg(name string) float64 {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeFloat {
		panic("Trying to use a non-float argument as float")
	}

	if val == nil {
		return 0
	}

	return val.(float64)
}

// BoolArg finds and returns a named bool argument
func (c *Context) BoolArg(name string) bool {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeBool {
		panic("Trying to use a non-bool argument as bool")
	}

	if val == nil {
		return false
	}

	return val.(bool)
}

// UserArg finds and returns a named User argument
func (c *Context) UserArg(name string) User {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeUserMention {
		panic("Trying to use a non-user argument as user")
	}

	if val == nil {
		return nil
	}

	return val.(User)
}

// ChannelArg finds and returns a named Channel argument
func (c *Context) ChannelArg(name string) Channel {
	return c.ChannelArgType(name, 255)
}

// ChannelArgType finds and returns Channel argument with a specified type
func (c *Context) ChannelArgType(name string, t ChannelType) Channel {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeChannelMention {
		panic("Trying to use a non-channel argument as channel")
	}

	if val == nil {
		return nil
	}

	ch := val.(Channel)

	if t != 255 && ch.Type() != t {
		return nil
	}

	return ch
}

// EmojiArg finds and returns an argument as an emoji
func (c *Context) EmojiArg(name string) Emoji {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeEmoji {
		panic("Trying to use a non-emoji argument as emoji")
	}

	if val == nil {
		return nil
	}

	return val.(Emoji)
}

// RoleArg finds and returns a named Role argument
func (c *Context) RoleArg(name string) Role {
	arg, val := c.arg(name)

	if arg.Type != ArgumentTypeRole {
		panic("Trying to use a non-role argument as role")
	}

	if val == nil {
		return nil
	}

	return val.(Role)
}
