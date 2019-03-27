package turnstile

import (
	"reflect"
	"testing"
	"time"
)

func TestExtractEvent(t *testing.T) {
	tests := []struct {
		name   string
		update Update
		want   *Event
	}{
		{
			name:   "private message",
			update: Update{Message: &Message{From: User{ID: 1, LanguageCode: "en"}, Chat: Chat{ID: 1, Type: "private"}}},
			want: &Event{
				Type:         "message",
				UserID:       1,
				LanguageCode: "en",
				ChatID:       1,
				ChatType:     "private",
			},
		},
		{
			name:   "group message",
			update: Update{Message: &Message{From: User{ID: 1, LanguageCode: "en"}, Chat: Chat{ID: -1, Type: "group"}}},
			want: &Event{
				Type:         "message",
				UserID:       1,
				LanguageCode: "en",
				ChatID:       -1,
				ChatType:     "group",
			},
		},
		{
			name:   "inline query",
			update: Update{InlineQuery: &InlineQuery{ID: "1", From: User{ID: 1, LanguageCode: "en"}}},
			want: &Event{
				Type:         "inline_query",
				UserID:       1,
				LanguageCode: "en",
			},
		},
		{
			name:   "chosen inline result",
			update: Update{ChosenInlineResult: &ChosenInlineResult{ResultID: "1", From: User{ID: 1, LanguageCode: "en"}}},
			want: &Event{
				Type:         "chosen_inline_result",
				UserID:       1,
				LanguageCode: "en",
			},
		},
		{
			name: "callback query with message",
			update: Update{CallbackQuery: &CallbackQuery{
				ID:      "1",
				From:    User{ID: 1, LanguageCode: "en"},
				Message: Message{From: User{ID: 2, LanguageCode: ""}, Chat: Chat{ID: 1, Type: "private"}},
			}},
			want: &Event{
				Type:         "callback_query",
				UserID:       1,
				LanguageCode: "en",
				ChatID:       1,
				ChatType:     "private",
			},
		},
		{
			name: "callback query without message",
			update: Update{CallbackQuery: &CallbackQuery{
				ID:   "1",
				From: User{ID: 1, LanguageCode: "en"},
			}},
			want: &Event{
				Type:         "callback_query",
				UserID:       1,
				LanguageCode: "en",
			},
		},
		{
			name:   "update without event",
			update: Update{},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExtractEvent(time.Time{}, tt.update); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
