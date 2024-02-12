package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var WebhookDefaultName = "disgo"

var _ WebhookManager = (*webhookManagerImpl)(nil)

func NewWebhookManager(client Client, logger *slog.Logger) WebhookManager {
	return &webhookManagerImpl{
		client: client,
		logger: logger,
	}
}

type WebhookManager interface {
	GetMessenger(channel discord.Channel) (discord.WebhookMessenger[Client], error)
	GetWebhook(channel discord.Channel) (discord.Webhook, error)
}

type webhookManagerImpl struct {
	client Client
	logger *slog.Logger
}

func (w *webhookManagerImpl) GetMessenger(channel discord.Channel) (discord.WebhookMessenger[Client], error) {
	if channel.Type() == discord.ChannelTypeDM || channel.Type() == discord.ChannelTypeGroupDM {
		return nil, discord.ErrUnsupportedType
	}

	if channel.Type() == discord.ChannelTypeGuildPublicThread ||
		channel.Type() == discord.ChannelTypeGuildPrivateThread ||
		channel.Type() == discord.ChannelTypeGuildNewsThread {
		return NewThreadWebhookMessenger(w.client, channel.ID(), *channel.(*discord.GuildThread).ParentID())
	}

	return NewChannelWebhookMessenger(w.client, channel.ID())
}

func (w *webhookManagerImpl) GetWebhook(channel discord.Channel) (discord.Webhook, error) {
	m, err := w.GetMessenger(channel)
	if err != nil {
		return nil, err
	}
	return m.Webhook(), nil
}

var (
	_ discord.WebhookMessenger[Client] = (*channelWebhookMessenger)(nil)
	_ discord.WebhookMessenger[Client] = (*threadWebhookMessenger)(nil)
)

func NewChannelWebhookMessenger(client Client, channelID snowflake.ID) (discord.WebhookMessenger[Client], error) {
	webhooks, err := client.Rest().GetWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		if webhook.Type() == discord.WebhookTypeIncoming && webhook.(discord.IncomingWebhook).User.ID == client.ApplicationID() {
			return channelWebhookMessenger{channelID: channelID, webhook: webhook.(discord.IncomingWebhook)}, nil
		}
	}

	webhook, err := client.Rest().CreateWebhook(channelID, discord.WebhookCreate{
		Name: WebhookDefaultName,
	})
	if err != nil {
		return nil, err
	}
	return channelWebhookMessenger{channelID: channelID, webhook: *webhook}, nil
}

type channelWebhookMessenger struct {
	channelID snowflake.ID
	webhook   discord.IncomingWebhook
}

func (c channelWebhookMessenger) SendWebhook(message discord.MessageBuilder, client Client, username, avatarURL, threadName string) (*discord.Message, error) {
	return client.Rest().CreateWebhookMessage(c.webhook.ID(), c.webhook.Token, message.BuildWebhookCreate(username, avatarURL, threadName), false, 0)
}

func (c channelWebhookMessenger) Webhook() discord.Webhook {
	return c.webhook
}

func (c channelWebhookMessenger) Send(message discord.MessageBuilder, client Client) (*discord.Message, error) {
	return c.SendWebhook(message, client, "", "", "")
}

func (c channelWebhookMessenger) Update(target discord.Object, message discord.MessageBuilder, client Client) (*discord.Message, error) {
	return client.Rest().UpdateWebhookMessage(c.webhook.ID(), c.webhook.Token, target.ID(), message.BuildWebhookUpdate(), 0)
}

func (c channelWebhookMessenger) Delete(message discord.Object, client Client) error {
	return client.Rest().DeleteWebhookMessage(c.webhook.ID(), c.webhook.Token, message.ID(), 0)
}

func NewThreadWebhookMessenger(client Client, channelID, threadID snowflake.ID) (discord.WebhookMessenger[Client], error) {
	webhooks, err := client.Rest().GetWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		if webhook.Type() == discord.WebhookTypeIncoming && webhook.(discord.IncomingWebhook).User.ID == client.ApplicationID() {
			return threadWebhookMessenger{channelID: channelID, threadID: threadID, webhook: webhook.(discord.IncomingWebhook)}, nil
		}
	}

	webhook, err := client.Rest().CreateWebhook(channelID, discord.WebhookCreate{
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
	webhook   discord.IncomingWebhook
}

func (t threadWebhookMessenger) SendWebhook(message discord.MessageBuilder, client Client, username, avatarURL, threadName string) (*discord.Message, error) {
	return client.Rest().CreateWebhookMessage(t.webhook.ID(), t.webhook.Token, message.BuildWebhookCreate(username, avatarURL, threadName), false, t.threadID)
}

func (t threadWebhookMessenger) Webhook() discord.Webhook {
	return t.webhook
}

func (t threadWebhookMessenger) Send(message discord.MessageBuilder, client Client) (*discord.Message, error) {
	return t.SendWebhook(message, client, "", "", "")
}

func (t threadWebhookMessenger) Update(target discord.Object, message discord.MessageBuilder, client Client) (*discord.Message, error) {
	return client.Rest().UpdateWebhookMessage(t.webhook.ID(), t.webhook.Token, target.ID(), message.BuildWebhookUpdate(), t.threadID)
}

func (t threadWebhookMessenger) Delete(message discord.Object, client Client) error {
	return client.Rest().DeleteWebhookMessage(t.webhook.ID(), t.webhook.Token, message.ID(), t.threadID)
}
