package discord

import (
	"fmt"
	"io"

	"github.com/disgoorg/disgo/internal/nillabe"
	"github.com/disgoorg/snowflake/v2"
)

func NewMessageBuilder() MessageBuilder {
	return &messageBuilderImpl{}
}

func NewMessageBuilderFromMessage(message Message) MessageBuilder {
	attachments := make([]AttachmentUpdate, len(message.Attachments))
	for i, attachment := range message.Attachments {
		attachments[i] = AttachmentKeep{ID: attachment.ID}
	}
	builder := messageBuilderImpl{
		TTS:              message.TTS,
		MessageReference: message.MessageReference,
		Content:          &message.Content,
		Embeds:           &message.Embeds,
		Components:       &message.Components,
		Attachments:      &attachments,
		Flags:            &message.Flags,
	}
	return &builder
}

type MessageBuilder interface {
	SetContent(content string) MessageBuilder
	SetContentf(content string, a ...any) MessageBuilder
	SetTTS(tts bool) MessageBuilder
	SetEmbeds(embeds ...Embed) MessageBuilder
	SetEmbed(i int, embed Embed) MessageBuilder
	AddEmbeds(embeds ...Embed) MessageBuilder
	ClearEmbeds() MessageBuilder
	RemoveEmbed(i int) MessageBuilder
	SetContainerComponents(containerComponents ...ContainerComponent) MessageBuilder
	SetContainerComponent(i int, container ContainerComponent) MessageBuilder
	AddActionRow(components ...InteractiveComponent) MessageBuilder
	AddContainerComponents(containers ...ContainerComponent) MessageBuilder
	RemoveContainerComponent(i int) MessageBuilder
	ClearContainerComponents() MessageBuilder
	AddStickers(stickerIds ...snowflake.ID) MessageBuilder
	SetStickers(stickerIds ...snowflake.ID) MessageBuilder
	ClearStickers() MessageBuilder
	SetFiles(files ...*File) MessageBuilder
	SetFile(i int, file *File) MessageBuilder
	AddFiles(files ...*File) MessageBuilder
	AddFile(name string, description string, reader io.Reader, flags ...FileFlags) MessageBuilder
	ClearFiles() MessageBuilder
	RemoveFile(i int) MessageBuilder
	SetAllowedMentions(allowedMentions *AllowedMentions) MessageBuilder
	ClearAllowedMentions() MessageBuilder
	SetMessageReference(messageReference *MessageReference) MessageBuilder
	SetMessageReferenceByID(messageID snowflake.ID) MessageBuilder
	SetFlags(flags MessageFlags) MessageBuilder
	AddFlags(flags ...MessageFlags) MessageBuilder
	RemoveFlags(flags ...MessageFlags) MessageBuilder
	ClearFlags() MessageBuilder
	SetEphemeral(ephemeral bool) MessageBuilder
	SetSuppressEmbeds(suppressEmbeds bool) MessageBuilder

	ClearContent() MessageBuilder
	RetainAttachments(attachments ...Attachment) MessageBuilder
	RetainAttachmentsByID(attachmentIDs ...snowflake.ID) MessageBuilder

	BuildCreate() MessageCreate
	BuildUpdate() MessageUpdate
	BuildWebhookCreate(username string, avatarURL string, threadName string) WebhookMessageCreate
	BuildWebhookUpdate() WebhookMessageUpdate

	messageBuilder()
}

type messageBuilderImpl struct {
	Nonce            string                `json:"nonce,omitempty"`
	TTS              bool                  `json:"tts,omitempty"`
	StickerIDs       []snowflake.ID        `json:"sticker_ids,omitempty"`
	MessageReference *MessageReference     `json:"message_reference,omitempty"`
	Content          *string               `json:"content,omitempty"`
	Embeds           *[]Embed              `json:"embeds,omitempty"`
	Components       *[]ContainerComponent `json:"components,omitempty"`
	Attachments      *[]AttachmentUpdate   `json:"attachments,omitempty"`
	Files            []*File               `json:"-"`
	AllowedMentions  *AllowedMentions      `json:"allowed_mentions,omitempty"`
	Flags            *MessageFlags         `json:"flags,omitempty"`
}

func (m *messageBuilderImpl) SetContent(content string) MessageBuilder {
	m.Content = &content
	return m
}

func (m *messageBuilderImpl) SetContentf(content string, args ...any) MessageBuilder {
	s := fmt.Sprintf(content, args...)
	m.Content = &s
	return m
}

func (m *messageBuilderImpl) SetTTS(tts bool) MessageBuilder {
	m.TTS = tts
	return m
}

func (m *messageBuilderImpl) SetEmbeds(embeds ...Embed) MessageBuilder {
	m.Embeds = &embeds
	return m
}

func (m *messageBuilderImpl) SetEmbed(i int, embed Embed) MessageBuilder {
	if m.Embeds == nil {
		m.Embeds = &[]Embed{}
	}
	if len(*m.Embeds) > i {
		(*m.Embeds)[i] = embed
	}
	return m
}

func (m *messageBuilderImpl) AddEmbeds(embeds ...Embed) MessageBuilder {
	if m.Embeds == nil {
		m.Embeds = &[]Embed{}
	}
	*m.Embeds = append(*m.Embeds, embeds...)
	return m
}

func (m *messageBuilderImpl) ClearEmbeds() MessageBuilder {
	return m.SetEmbeds()
}

func (m *messageBuilderImpl) RemoveEmbed(i int) MessageBuilder {
	if m.Embeds == nil {
		m.Embeds = &[]Embed{}
	}
	*m.Embeds = append((*m.Embeds)[:i], (*m.Embeds)[i+1:]...)
	return m
}

func (m *messageBuilderImpl) SetContainerComponents(containerComponents ...ContainerComponent) MessageBuilder {
	m.Components = &containerComponents
	return m
}

func (m *messageBuilderImpl) SetContainerComponent(i int, container ContainerComponent) MessageBuilder {
	if m.Components == nil {
		m.Components = &[]ContainerComponent{}
	}
	if len(*m.Components) > i {
		(*m.Components)[i] = container
	}
	return m
}

func (m *messageBuilderImpl) AddActionRow(components ...InteractiveComponent) MessageBuilder {
	if m.Components == nil {
		m.Components = &[]ContainerComponent{}
	}
	*m.Components = append(*m.Components, ActionRowComponent(components))
	return m
}

func (m *messageBuilderImpl) AddContainerComponents(containers ...ContainerComponent) MessageBuilder {
	if m.Components == nil {
		m.Components = &[]ContainerComponent{}
	}
	*m.Components = append(*m.Components, containers...)
	return m
}

func (m *messageBuilderImpl) RemoveContainerComponent(i int) MessageBuilder {
	if m.Components == nil {
		m.Components = &[]ContainerComponent{}
	}
	*m.Components = append((*m.Components)[:i], (*m.Components)[i+1:]...)
	return m
}

func (m *messageBuilderImpl) ClearContainerComponents() MessageBuilder {
	return m.SetContainerComponents()
}

func (m *messageBuilderImpl) AddStickers(stickerIds ...snowflake.ID) MessageBuilder {
	m.StickerIDs = append(m.StickerIDs, stickerIds...)
	return m
}

func (m *messageBuilderImpl) SetStickers(stickerIds ...snowflake.ID) MessageBuilder {
	m.StickerIDs = stickerIds
	return m
}

func (m *messageBuilderImpl) ClearStickers() MessageBuilder {
	return m.SetStickers()
}

func (m *messageBuilderImpl) SetFiles(files ...*File) MessageBuilder {
	m.Files = files
	return m
}

func (m *messageBuilderImpl) SetFile(i int, file *File) MessageBuilder {
	if len(m.Files) > i {
		m.Files[i] = file
	}
	return m
}

func (m *messageBuilderImpl) AddFiles(files ...*File) MessageBuilder {
	m.Files = append(m.Files, files...)
	return m
}

func (m *messageBuilderImpl) AddFile(name string, description string, reader io.Reader, flags ...FileFlags) MessageBuilder {
	m.Files = append(m.Files, NewFile(name, description, reader, flags...))
	return m
}

func (m *messageBuilderImpl) ClearFiles() MessageBuilder {
	return m.SetFiles()
}

func (m *messageBuilderImpl) RemoveFile(i int) MessageBuilder {
	if len(m.Files) > i {
		m.Files = append(m.Files[:i], m.Files[i+1:]...)
	}
	return m
}

func (m *messageBuilderImpl) SetAllowedMentions(allowedMentions *AllowedMentions) MessageBuilder {
	m.AllowedMentions = allowedMentions
	return m
}

func (m *messageBuilderImpl) ClearAllowedMentions() MessageBuilder {
	return m.SetAllowedMentions(nil)
}

func (m *messageBuilderImpl) SetMessageReference(messageReference *MessageReference) MessageBuilder {
	m.MessageReference = messageReference
	return m
}

func (m *messageBuilderImpl) SetMessageReferenceByID(messageID snowflake.ID) MessageBuilder {
	m.MessageReference = &MessageReference{MessageID: &messageID}
	return m
}

func (m *messageBuilderImpl) SetFlags(flags MessageFlags) MessageBuilder {
	m.Flags = &flags
	return m
}

func (m *messageBuilderImpl) AddFlags(flags ...MessageFlags) MessageBuilder {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	*m.Flags = m.Flags.Add(flags...)
	return m
}

func (m *messageBuilderImpl) RemoveFlags(flags ...MessageFlags) MessageBuilder {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	*m.Flags = m.Flags.Remove(flags...)
	return m

}

func (m *messageBuilderImpl) ClearFlags() MessageBuilder {
	return m.SetFlags(MessageFlagsNone)
}

func (m *messageBuilderImpl) SetEphemeral(ephemeral bool) MessageBuilder {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	if ephemeral {
		*m.Flags = m.Flags.Add(MessageFlagEphemeral)
	} else {
		*m.Flags = m.Flags.Remove(MessageFlagEphemeral)
	}
	return m
}

func (m *messageBuilderImpl) SetSuppressEmbeds(suppressEmbeds bool) MessageBuilder {
	if m.Flags == nil {
		m.Flags = new(MessageFlags)
	}
	if suppressEmbeds {
		*m.Flags = m.Flags.Add(MessageFlagSuppressEmbeds)
	} else {
		*m.Flags = m.Flags.Remove(MessageFlagSuppressEmbeds)
	}
	return m
}

func (m *messageBuilderImpl) ClearContent() MessageBuilder {
	return m.SetContent("")
}

func (m *messageBuilderImpl) RetainAttachments(attachments ...Attachment) MessageBuilder {
	if m.Attachments == nil {
		m.Attachments = new([]AttachmentUpdate)
	}
	for _, attachment := range attachments {
		*m.Attachments = append(*m.Attachments, AttachmentKeep{ID: attachment.ID})
	}
	return m
}

func (m *messageBuilderImpl) RetainAttachmentsByID(attachmentIDs ...snowflake.ID) MessageBuilder {
	if m.Attachments == nil {
		m.Attachments = new([]AttachmentUpdate)
	}
	for _, attachmentID := range attachmentIDs {
		*m.Attachments = append(*m.Attachments, AttachmentKeep{ID: attachmentID})
	}
	return m
}

func (m *messageBuilderImpl) BuildCreate() MessageCreate {
	var attachments []AttachmentCreate
	attachments = parseAttachments(m.Files)
	return MessageCreate{
		Nonce:            m.Nonce,
		Content:          nillabe.NonNil(m.Content),
		TTS:              m.TTS,
		Embeds:           nillabe.NonNil(m.Embeds),
		Components:       nillabe.NonNil(m.Components),
		StickerIDs:       m.StickerIDs,
		Files:            m.Files,
		Attachments:      attachments,
		AllowedMentions:  m.AllowedMentions,
		MessageReference: m.MessageReference,
		Flags:            nillabe.NonNil(m.Flags),
	}
}

func (m *messageBuilderImpl) BuildUpdate() MessageUpdate {
	return MessageUpdate{
		Content:         m.Content,
		Embeds:          m.Embeds,
		Components:      m.Components,
		Attachments:     m.Attachments,
		Files:           m.Files,
		AllowedMentions: m.AllowedMentions,
		Flags:           m.Flags,
	}
}

func (m *messageBuilderImpl) BuildWebhookCreate(username string, avatarURL string, threadName string) WebhookMessageCreate {
	return WebhookMessageCreate{
		Content:         nillabe.NonNil(m.Content),
		Username:        username,
		AvatarURL:       avatarURL,
		TTS:             m.TTS,
		Embeds:          nillabe.NonNil(m.Embeds),
		Components:      nillabe.NonNil(m.Components),
		Files:           m.Files,
		AllowedMentions: m.AllowedMentions,
		Flags:           nillabe.NonNil(m.Flags),
		ThreadName:      threadName,
	}
}

func (m *messageBuilderImpl) BuildWebhookUpdate() WebhookMessageUpdate {
	return WebhookMessageUpdate{
		Content:         m.Content,
		Embeds:          m.Embeds,
		Components:      m.Components,
		Attachments:     m.Attachments,
		Files:           m.Files,
		AllowedMentions: m.AllowedMentions,
	}
}

func (m *messageBuilderImpl) messageBuilder() {}
