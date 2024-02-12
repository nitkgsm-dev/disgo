package discord

type WebhookMessenger[T any] interface {
	Messenger[T]
	Webhook() Webhook
	SendWebhook(message MessageBuilder, client T, username, avatarURL, threadName string) (*Message, error)
}
