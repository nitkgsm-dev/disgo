package discord

type WebhookMessenger interface {
	Messenger
	Webhook() Webhook
	SendWebhook(message MessageBuilder, username, avatarURL, threadName string) (*Message, error)
}
