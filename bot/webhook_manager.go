package bot

import (
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

var (
	WebhookDefaultName                 = "disgo"
	WebhookDefaultAvatar *discord.Icon = nil
)

var _ WebhookManager = (*webhookManagerImpl)(nil)

func NewWebhookManager(client Client, logger *slog.Logger) WebhookManager {
	return &webhookManagerImpl{
		client: client,
		logger: logger,
	}
}

type WebhookManager interface {
	GetMessenger(channel discord.Channel) (discord.WebhookMessenger, error)
	GetWebhook(channel discord.Channel) (discord.Webhook, error)
}

type webhookManagerImpl struct {
	client Client
	logger *slog.Logger
}

func (w *webhookManagerImpl) GetMessenger(channel discord.Channel) (discord.WebhookMessenger, error) {
	if channel.Type() == discord.ChannelTypeDM || channel.Type() == discord.ChannelTypeGroupDM {
		return nil, discord.ErrUnsupportedType
	}

	if channel.Type() == discord.ChannelTypeGuildPublicThread ||
		channel.Type() == discord.ChannelTypeGuildPrivateThread ||
		channel.Type() == discord.ChannelTypeGuildNewsThread {
		c, ok := channel.(discord.GuildThread)
		if ic, iok := channel.(discord.InteractionChannel); !ok && iok {
			c = ic.MessageChannel.(discord.GuildThread)
		} else {
			return nil, discord.ErrUnsupportedType
		}
		return NewThreadWebhookMessenger(w.client, *c.ParentID(), c.ID())
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
	_ discord.WebhookMessenger = (*channelWebhookMessenger)(nil)
	_ discord.WebhookMessenger = (*threadWebhookMessenger)(nil)
)

func NewChannelWebhookMessenger(client Client, channelID snowflake.ID) (discord.WebhookMessenger, error) {
	webhooks, err := client.Rest().GetWebhooks(channelID)
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		if webhook.Type() == discord.WebhookTypeIncoming && webhook.(discord.IncomingWebhook).User.ID == client.ApplicationID() {
			return channelWebhookMessenger{channelID: channelID, webhook: webhook.(discord.IncomingWebhook), client: client}, nil
		}
	}

	webhook, err := client.Rest().CreateWebhook(channelID, discord.WebhookCreate{
		Name:   WebhookDefaultName,
		Avatar: WebhookDefaultAvatar,
	})
	if err != nil {
		return nil, err
	}
	return channelWebhookMessenger{channelID: channelID, webhook: *webhook, client: client}, nil
}

type channelWebhookMessenger struct {
	channelID snowflake.ID
	webhook   discord.IncomingWebhook
	client    Client
}

func (c channelWebhookMessenger) SendWebhook(message discord.MessageBuilder, username, avatarURL, threadName string) (*discord.Message, error) {
	return c.client.Rest().CreateWebhookMessage(c.webhook.ID(), c.webhook.Token, message.BuildWebhookCreate(username, avatarURL, threadName), false, 0)
}

func (c channelWebhookMessenger) Webhook() discord.Webhook {
	return c.webhook
}

func (c channelWebhookMessenger) Send(message discord.MessageBuilder) (*discord.Message, error) {
	return c.SendWebhook(message, "", "", "")
}

func (c channelWebhookMessenger) Update(target snowflake.ID, message discord.MessageBuilder) (*discord.Message, error) {
	return c.client.Rest().UpdateWebhookMessage(c.webhook.ID(), c.webhook.Token, target, message.BuildWebhookUpdate(), 0)
}

func (c channelWebhookMessenger) Delete(message snowflake.ID) error {
	return c.client.Rest().DeleteWebhookMessage(c.webhook.ID(), c.webhook.Token, message, 0)
}

func NewThreadWebhookMessenger(client Client, parentID, threadID snowflake.ID) (discord.WebhookMessenger, error) {
	webhooks, err := client.Rest().GetWebhooks(parentID)
	if err != nil {
		return nil, err
	}

	for _, webhook := range webhooks {
		if webhook.Type() == discord.WebhookTypeIncoming && webhook.(discord.IncomingWebhook).User.ID == client.ApplicationID() {
			return threadWebhookMessenger{channelID: parentID, threadID: threadID, webhook: webhook.(discord.IncomingWebhook), client: client}, nil
		}
	}

	webhook, err := client.Rest().CreateWebhook(parentID, discord.WebhookCreate{
		Name:   WebhookDefaultName,
		Avatar: WebhookDefaultAvatar,
	})
	if err != nil {
		return nil, err
	}
	return threadWebhookMessenger{channelID: parentID, threadID: threadID, webhook: *webhook, client: client}, nil
}

type threadWebhookMessenger struct {
	channelID snowflake.ID
	threadID  snowflake.ID
	webhook   discord.IncomingWebhook
	client    Client
}

func (t threadWebhookMessenger) SendWebhook(message discord.MessageBuilder, username, avatarURL, threadName string) (*discord.Message, error) {
	return t.client.Rest().CreateWebhookMessage(t.webhook.ID(), t.webhook.Token, message.BuildWebhookCreate(username, avatarURL, threadName), false, t.threadID)
}

func (t threadWebhookMessenger) Webhook() discord.Webhook {
	return t.webhook
}

func (t threadWebhookMessenger) Send(message discord.MessageBuilder) (*discord.Message, error) {
	return t.SendWebhook(message, "", "", "")
}

func (t threadWebhookMessenger) Update(target snowflake.ID, message discord.MessageBuilder) (*discord.Message, error) {
	return t.client.Rest().UpdateWebhookMessage(t.webhook.ID(), t.webhook.Token, target, message.BuildWebhookUpdate(), t.threadID)
}

func (t threadWebhookMessenger) Delete(message snowflake.ID) error {
	return t.client.Rest().DeleteWebhookMessage(t.webhook.ID(), t.webhook.Token, message, t.threadID)
}
