package turnstile

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

type MemoryCollector struct {
	Event Event
}

func (c *MemoryCollector) Collect(e Event) error {
	c.Event = e
	return nil
}

type MockMiddleware struct {
	Update Update
}

func (m *MockMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	err := json.NewDecoder(r.Body).Decode(&m.Update)
	if err != nil {
		return 0, err
	}
	return 0, r.Body.Close()
}

func TestTurnstile_ServeHTTP(t *testing.T) {
	update := Update{Message: &Message{From: User{ID: 1, LanguageCode: "en"}, Chat: Chat{ID: 1, Type: "private"}}}
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(update)
	if err != nil {
		t.Fatal(err)
	}
	r := &http.Request{Body: ioutil.NopCloser(&b)}
	collector := &MemoryCollector{}
	next := &MockMiddleware{}
	turnstile := Turnstile{
		collector: collector,
		next:      next,
	}
	_, err = turnstile.ServeHTTP(nil, r)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("collected expected event", func(t *testing.T) {
		wantEvent := Event{
			Type:         "message",
			UserID:       1,
			LanguageCode: "en",
			ChatID:       1,
			ChatType:     "private",
		}
		gotEvent := collector.Event
		gotEvent.Time = time.Time{}
		if gotEvent != wantEvent {
			t.Errorf("collector.Event = %v, want %v", gotEvent, wantEvent)
		}
	})

	t.Run("request body is still available to next middleware", func(t *testing.T) {
		wantUpdate := update
		if gotUpdate := next.Update; !reflect.DeepEqual(gotUpdate, wantUpdate) {
			t.Errorf("next.Update = %v, want %v", gotUpdate, wantUpdate)
		}
	})
}
