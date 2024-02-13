package discord

import "github.com/disgoorg/snowflake/v2"

type Messenger interface {
	Send(message MessageBuilder) (*Message, error)
	Update(target snowflake.ID, message MessageBuilder) (*Message, error)
	Delete(message snowflake.ID) error
}
