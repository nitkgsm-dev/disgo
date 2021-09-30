package core

import (
	"io"
	"runtime/debug"

	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/json"
)

var (
	DefaultEventManagerConfig = EventManagerConfig{}
)

type EventManagerConfig struct {
	EventListeners           []EventListener
	VoiceDispatchInterceptor VoiceDispatchInterceptor
	RawEventsEnabled         bool

	GatewayHandlers   map[discord.GatewayEventType]GatewayEventHandler
	HTTPServerHandler HTTPServerEventHandler

	NewMessageCollector func(channel *Channel, filter MessageFilter) (<-chan *Message, func())
}

var _ EventManager = (*eventManagerImpl)(nil)

func NewEventManager(bot *Bot, config *EventManagerConfig) EventManager {
	if config == nil {
		config = &DefaultEventManagerConfig
	}

	return &eventManagerImpl{
		bot:    bot,
		config: *config,
	}
}

// EventManager lets you listen for specific events triggered by raw gateway events
type EventManager interface {
	Bot() *Bot
	Close()
	Config() EventManagerConfig

	SetVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor)
	AddEventListeners(eventListeners ...EventListener)
	RemoveEventListeners(eventListeners ...EventListener)
	HandleGateway(gatewayEventType discord.GatewayEventType, sequenceNumber int, payload io.Reader)
	HandleHTTP(responseChannel chan<- discord.InteractionResponse, payload io.Reader)
	Dispatch(event Event)
}

// EventListener is used to create new EventListener to listen to events
type EventListener interface {
	OnEvent(event interface{})
}

// Event the basic interface each event implement
type Event interface {
	Bot() *Bot
	SequenceNumber() int
}

// GatewayEventHandler is used to handle Gateway Event(s)
type GatewayEventHandler interface {
	EventType() discord.GatewayEventType
	New() interface{}
	HandleGatewayEvent(bot *Bot, sequenceNumber int, v interface{})
}

// HTTPServerEventHandler is used to handle HTTP Event(s)
type HTTPServerEventHandler interface {
	New() interface{}
	HandleHTTPEvent(bot *Bot, responseChannel chan<- discord.InteractionResponse, v interface{})
}

// eventManagerImpl is the implementation of core.EventManager
type eventManagerImpl struct {
	bot    *Bot
	config EventManagerConfig
}

// Bot returns the core.Bot instance used by the core.EventManager
func (e *eventManagerImpl) Bot() *Bot {
	return e.bot
}

// Close closes all goroutines created by the core.EventManager
func (e *eventManagerImpl) Close() {
	e.Bot().Logger.Info("closing eventManager goroutines...")
}

func (e *eventManagerImpl) Config() EventManagerConfig {
	return e.config
}

func (e *eventManagerImpl) SetVoiceDispatchInterceptor(voiceDispatchInterceptor VoiceDispatchInterceptor) {
	e.config.VoiceDispatchInterceptor = voiceDispatchInterceptor
}

// HandleGateway calls the correct core.EventHandler
func (e *eventManagerImpl) HandleGateway(gatewayEventType discord.GatewayEventType, sequenceNumber int, reader io.Reader) {

	if handler, ok := e.config.GatewayHandlers[gatewayEventType]; ok {
		v := handler.New()
		if v != nil {
			if err := json.NewDecoder(reader).Decode(&v); err != nil {
				e.Bot().Logger.Errorf("error while unmarshalling event '%s'. error: %s", gatewayEventType, err.Error())
			}
		}
		handler.HandleGatewayEvent(e.Bot(), sequenceNumber, v)
	} else {
		e.Bot().Logger.Warnf("no handler for gateway event '%s' found", gatewayEventType)
	}
}

// HandleHTTP calls the correct core.EventHandler
func (e *eventManagerImpl) HandleHTTP(responseChannel chan<- discord.InteractionResponse, reader io.Reader) {
	v := e.config.HTTPServerHandler.New()
	if err := json.NewDecoder(reader).Decode(&v); err != nil {
		e.Bot().Logger.Error("error while unmarshalling httpserver event. error: ", err)
	}
	e.config.HTTPServerHandler.HandleHTTPEvent(e.Bot(), responseChannel, v)
}

// Dispatch dispatches a new event to the client
func (e *eventManagerImpl) Dispatch(event Event) {
	for i := range e.config.EventListeners {
		listener := e.config.EventListeners[i]
		go func() {
			defer func() {
				if r := recover(); r != nil {
					e.Bot().Logger.Error("recovered from panic in event listener: ", r)
					debug.PrintStack()
					return
				}
			}()
			listener.OnEvent(event)
		}()
	}
}

// AddEventListeners adds one or more core.EventListener(s) to the core.EventManager
func (e *eventManagerImpl) AddEventListeners(listeners ...EventListener) {
	for _, listener := range listeners {
		e.config.EventListeners = append(e.config.EventListeners, listener)
	}
}

// RemoveEventListeners removes one or more core.EventListener(s) from the core.EventManager
func (e *eventManagerImpl) RemoveEventListeners(listeners ...EventListener) {
	for _, listener := range listeners {
		for i, l := range e.config.EventListeners {
			if l == listener {
				e.config.EventListeners = append(e.config.EventListeners[:i], e.config.EventListeners[i+1:]...)
				break
			}
		}
	}
}