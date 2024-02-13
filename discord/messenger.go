package discord

type Messenger interface {
	Send(message MessageBuilder) (*Message, error)
	Update(target Object, message MessageBuilder) (*Message, error)
	Delete(message Object) error
}
