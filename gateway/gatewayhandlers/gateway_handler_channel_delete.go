package gatewayhandlers

import (
	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/events"
)

// gatewayHandlerChannelDelete handles core.GatewayEventChannelDelete
type gatewayHandlerChannelDelete struct{}

// EventType returns the core.GatewayGatewayEventType
func (h *gatewayHandlerChannelDelete) EventType() discord.GatewayEventType {
	return discord.GatewayEventTypeChannelDelete
}

// New constructs a new payload receiver for the raw gateway event
func (h *gatewayHandlerChannelDelete) New() interface{} {
	return &discord.UnmarshalChannel{}
}

// HandleGatewayEvent handles the specific raw gateway event
func (h *gatewayHandlerChannelDelete) HandleGatewayEvent(bot *core.Bot, sequenceNumber int, v interface{}) {
	channel := v.(*discord.UnmarshalChannel).Channel

	genericChannelEvent := &events.GenericChannelEvent{
		GenericEvent: events.NewGenericEvent(bot, sequenceNumber),
		ChannelID:    channel.ID(),
		Channel:      bot.EntityBuilder.CreateChannel(channel, core.CacheStrategyNo),
	}

	if ch, ok := channel.(discord.GuildChannel); ok {
		bot.EventManager.Dispatch(&events.GuildChannelDeleteEvent{
			GenericGuildChannelEvent: &events.GenericGuildChannelEvent{
				GenericChannelEvent: genericChannelEvent,
				GuildID:             ch.GuildID(),
			},
		})
	} else {
		bot.EventManager.Dispatch(&events.DMChannelDeleteEvent{
			GenericDMChannelEvent: &events.GenericDMChannelEvent{
				GenericChannelEvent: genericChannelEvent,
			},
		})
	}
}
