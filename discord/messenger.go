package discord

type Messenger interface {
	Send(message MessageBuilder, client ClientInterface) (*Message, error)
	Update(target Object, message MessageBuilder, client ClientInterface) (*Message, error)
	Delete(message Object, client ClientInterface) error
}
