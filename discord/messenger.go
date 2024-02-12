package discord

type Messenger[T any] interface {
	Send(message MessageBuilder, client T) (*Message, error)
	Update(target Object, message MessageBuilder, client T) (*Message, error)
	Delete(message Object, client T) error
}
