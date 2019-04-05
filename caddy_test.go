package turnstile

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
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
	Body []byte
}

func (m *MockMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}
	m.Body = b
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
	turnstile := New(collector, next)
	_, err = turnstile.ServeHTTP(nil, r)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("collected expected event", func(t *testing.T) {
		turnstile.wg.Wait()
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
		var gotUpdate Update
		err := json.Unmarshal(next.Body, &gotUpdate)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(gotUpdate, wantUpdate) {
			t.Errorf("next.Update = %v, want %v", gotUpdate, wantUpdate)
		}
	})

	t.Run("request body is still available to the next middleware even when json is invalid", func(t *testing.T) {
		r := &http.Request{Body: ioutil.NopCloser(strings.NewReader("invalid json"))}
		collector := &MemoryCollector{}
		next := &MockMiddleware{}
		turnstile := New(collector, next)
		_, err = turnstile.ServeHTTP(nil, r)
		if err != nil {
			t.Fatal(err)
		}
		wantBody := "invalid json"
		if gotBody := string(next.Body); wantBody != gotBody {
			t.Errorf("next.Body = %v, want %v", gotBody, wantBody)
		}
	})
}
