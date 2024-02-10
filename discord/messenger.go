package discord

import "github.com/disgoorg/disgo/bot"

type Messenger interface {
	Send(message MessageBuilder, client bot.Client) (*Message, error)
	Update(target Object, message MessageBuilder, client bot.Client) (*Message, error)
	Delete(message Object, client bot.Client) error
}
