package turnstile

import (
	"time"
)

// Types of updates
const (
	EventTypeMessage            = "message"
	EventTypeInlineQuery        = "inline_query"
	EventTypeChosenInlineResult = "chosen_inline_result"
	EventTypeCallbackQuery      = "callback_query"
)

// Event represents a single interaction with aaa bot.
type Event struct {
	Time         time.Time
	Type         string
	UserID       int64
	LanguageCode string
	ChatID       int64
	ChatType     string
}

// NewMessageEvent returns an event representing an incoming message.
func NewMessageEvent(t time.Time, message Message) Event {
	return Event{
		Time:         t,
		Type:         EventTypeMessage,
		UserID:       message.From.ID,
		LanguageCode: message.From.LanguageCode,
		ChatID:       message.Chat.ID,
		ChatType:     message.Chat.Type,
	}
}

// NewInlineQueryEvent returns an event representing an incoming inline query.
func NewInlineQueryEvent(t time.Time, query InlineQuery) Event {
	return Event{
		Time:         t,
		Type:         EventTypeInlineQuery,
		UserID:       query.From.ID,
		LanguageCode: query.From.LanguageCode,
	}
}

// NewChosenInlineResultEvent returns an event representing a chosen inline result.
func NewChosenInlineResultEvent(t time.Time, result ChosenInlineResult) Event {
	return Event{
		Time:         t,
		Type:         EventTypeChosenInlineResult,
		UserID:       result.From.ID,
		LanguageCode: result.From.LanguageCode,
	}
}

// NewCallbackQueryEvent returns an event representing an incoming callback query.
func NewCallbackQueryEvent(t time.Time, query CallbackQuery) Event {
	return Event{
		Time:         t,
		Type:         EventTypeCallbackQuery,
		UserID:       query.From.ID,
		LanguageCode: query.From.LanguageCode,
		ChatID:       query.Message.Chat.ID,
		ChatType:     query.Message.Chat.Type,
	}
}

// ExtractEvent returns a pointer to an event derived from an incoming update,
// or nil if the update does not contain an event.
func ExtractEvent(t time.Time, update Update) *Event {
	var event Event
	switch {
	case update.Message != nil:
		event = NewMessageEvent(t, *update.Message)
	case update.InlineQuery != nil:
		event = NewInlineQueryEvent(t, *update.InlineQuery)
	case update.ChosenInlineResult != nil:
		event = NewChosenInlineResultEvent(t, *update.ChosenInlineResult)
	case update.CallbackQuery != nil:
		event = NewCallbackQueryEvent(t, *update.CallbackQuery)
	default:
		return nil
	}
	return &event
}
