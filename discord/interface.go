package discord

import "github.com/disgoorg/snowflake/v2"

type ClientInterface interface {
	Rest() RestInterface
	ApplicationID() snowflake.ID
}

type RestInterface interface {
	GetWebhooks(channelID snowflake.ID) ([]Webhook, error)
	CreateWebhook(channelID snowflake.ID, webhookCreate WebhookCreate) (*IncomingWebhook, error)
	CreateWebhookMessage(webhookID snowflake.ID, webhookToken string, webhookMessageCreate WebhookMessageCreate, wait bool, threadID snowflake.ID) (*Message, error)
	UpdateWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, webhookMessageUpdate WebhookMessageUpdate, threadID snowflake.ID) (*Message, error)
	DeleteWebhookMessage(webhookID snowflake.ID, webhookToken string, messageID snowflake.ID, threadID snowflake.ID) error
	CreateMessage(channelID snowflake.ID, messageCreate MessageCreate) (*Message, error)
	UpdateMessage(channelID snowflake.ID, messageID snowflake.ID, update MessageUpdate) (*Message, error)
	DeleteMessage(channelID snowflake.ID, messageID snowflake.ID) error
}
