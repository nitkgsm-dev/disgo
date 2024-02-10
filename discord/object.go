package discord

import "github.com/disgoorg/snowflake/v2"

type Object interface {
	ID() snowflake.ID
}
