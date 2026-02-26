package astral

import (
	"fmt"
	"io"
	"strings"
)

// Responder represents an available responder
// This can send messages to channels, direct messages, interactions, etc.
type Responder interface {
	Usage(usage ...string) (Message, error)
	Send(text string) (Message, error)
	Sendf(format string, a ...interface{}) (Message, error)
	SendFile(name string, r io.Reader) (Message, error)
	Reply(text string) (Message, error)
	Replyf(format string, a ...interface{}) (Message, error)
	ReplyTo(to ID, text string) (Message, error)
	ReplyEmbed(embed Embed) (Message, error)
	ReplyFile(name string, r io.Reader) (Message, error)
	Respond(r Response) (Message, error)
	Acknowledge() error
	Error(message string) error
}

type DefaultResponder struct {
	ctx     *Context
	client  Client
	channel Channel
	user    User
}

func (r *DefaultResponder) Usage(usage ...string) (Message, error) {
	if len(usage) == 0 {
		usage = []string{r.ctx.Route.Usage}
	}

	usage[0] = strings.Replace(usage[0], "{command}", strings.Join(r.ctx.Route.Path(), " "), -1)

	return r.Reply(usage[0])
}

func (r *DefaultResponder) Send(text string) (Message, error) {
	return r.channel.SendMessage(text)
}

func (r *DefaultResponder) Sendf(format string, a ...interface{}) (Message, error) {
	return r.channel.SendMessage(fmt.Sprintf(format, a...))
}

func (r *DefaultResponder) SendFile(name string, reader io.Reader) (Message, error) {
	return r.Respond(Response{
		Files: []File{
			{
				Name:   name,
				Reader: reader,
			},
		},
	})
}

func (r *DefaultResponder) Reply(text string) (Message, error) {
	return r.Sendf("%s: %s", r.user.Mention(), text)
}

func (r *DefaultResponder) Replyf(format string, a ...interface{}) (Message, error) {
	return r.Reply(fmt.Sprintf(format, a...))
}

func (r *DefaultResponder) ReplyTo(to ID, text string) (Message, error) {
	user, err := r.client.User(to)

	if err != nil {
		return nil, err
	}

	return r.Sendf("%s: %s", user.Mention(), text)
}

func (r *DefaultResponder) ReplyEmbed(embed Embed) (Message, error) {
	//TODO implement me
	panic("implement me")
}

func (r *DefaultResponder) ReplyFile(name string, reader io.Reader) (Message, error) {
	return r.Respond(Response{
		Content: r.user.Mention(),
		Files: []File{
			{Name: name, Reader: reader},
		},
	})
}

func (r *DefaultResponder) Respond(res Response) (Message, error) {
	return r.channel.Send(res)
}

func (r *DefaultResponder) Acknowledge() error {
	return nil
}

func (r *DefaultResponder) Error(message string) error {
	return nil
}

func NewDefaultResponder(channel Channel, user User) Responder {
	return &DefaultResponder{
		channel: channel,
		user:    user,
	}
}
