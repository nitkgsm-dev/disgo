package discord

import (
	"github.com/disgoorg/snowflake/v2"
)

var WebhookDefaultName = "disgo"

var (
	_ WebhookMessenger = (*channelWebhookMessenger)(nil)
	_ WebhookMessenger = (*threadWebhookMessenger)(nil)
)

type WebhookMessenger interface {
	Messenger
	Webhook() Webhook
	SendWebhook(message MessageBuilder, client ClientInterface, username, avatarURL, threadName string) (*Message, error)
}

func NewChannelWebhookMessenger(client ClientInterface, channelID snowflake.ID) (WebhookMessenger, error) {
	webhooks, err := client.Rest().GetWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		if webhook.Type() == WebhookTypeIncoming && webhook.(*IncomingWebhook).User.ID == client.ApplicationID() {
			return channelWebhookMessenger{channelID: channelID, webhook: webhook.(IncomingWebhook)}, nil
		}
	}

	webhook, err := client.Rest().CreateWebhook(channelID, WebhookCreate{
		Name: WebhookDefaultName,
	})
	if err != nil {
		return nil, err
	}
	return channelWebhookMessenger{channelID: channelID, webhook: *webhook}, nil
}

type channelWebhookMessenger struct {
	channelID snowflake.ID
	webhook   IncomingWebhook
}

func (c channelWebhookMessenger) SendWebhook(message MessageBuilder, client ClientInterface, username, avatarURL, threadName string) (*Message, error) {
	return client.Rest().CreateWebhookMessage(c.webhook.id, c.webhook.Token, message.buildWebhookCreate(username, avatarURL, threadName), false, 0)
}

func (c channelWebhookMessenger) Webhook() Webhook {
	return c.webhook
}

func (c channelWebhookMessenger) Send(message MessageBuilder, client ClientInterface) (*Message, error) {
	return c.SendWebhook(message, client, "", "", "")
}

func (c channelWebhookMessenger) Update(target Object, message MessageBuilder, client ClientInterface) (*Message, error) {
	return client.Rest().UpdateWebhookMessage(c.webhook.id, c.webhook.Token, target.ID(), message.buildWebhookUpdate(), 0)
}

func (c channelWebhookMessenger) Delete(message Object, client ClientInterface) error {
	return client.Rest().DeleteWebhookMessage(c.webhook.id, c.webhook.Token, message.ID(), 0)
}

func NewThreadWebhookMessenger(client ClientInterface, channelID, threadID snowflake.ID) (WebhookMessenger, error) {
	webhooks, err := client.Rest().GetWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		if webhook.Type() == WebhookTypeIncoming && webhook.(*IncomingWebhook).User.ID == client.ApplicationID() {
			return threadWebhookMessenger{channelID: channelID, threadID: threadID, webhook: webhook.(IncomingWebhook)}, nil
		}
	}

	webhook, err := client.Rest().CreateWebhook(channelID, WebhookCreate{
		Name: WebhookDefaultName,
	})
	if err != nil {
		return nil, err
	}
	return threadWebhookMessenger{channelID: channelID, threadID: threadID, webhook: *webhook}, nil
}

type threadWebhookMessenger struct {
	channelID snowflake.ID
	threadID  snowflake.ID
	webhook   IncomingWebhook
}

func (t threadWebhookMessenger) SendWebhook(message MessageBuilder, client ClientInterface, username, avatarURL, threadName string) (*Message, error) {
	return client.Rest().CreateWebhookMessage(t.webhook.id, t.webhook.Token, message.buildWebhookCreate(username, avatarURL, threadName), false, t.threadID)
}

func (t threadWebhookMessenger) Webhook() Webhook {
	return t.webhook
}

func (t threadWebhookMessenger) Send(message MessageBuilder, client ClientInterface) (*Message, error) {
	return t.SendWebhook(message, client, "", "", "")
}

func (t threadWebhookMessenger) Update(target Object, message MessageBuilder, client ClientInterface) (*Message, error) {
	return client.Rest().UpdateWebhookMessage(t.webhook.id, t.webhook.Token, target.ID(), message.buildWebhookUpdate(), t.threadID)
}

func (t threadWebhookMessenger) Delete(message Object, client ClientInterface) error {
	return client.Rest().DeleteWebhookMessage(t.webhook.id, t.webhook.Token, message.ID(), t.threadID)
}
