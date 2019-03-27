package turnstile

import (
	"testing"
)

func TestSQLiteCollector(t *testing.T) {
	c, err := NewSQLiteCollector(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	event := Event{
		Type:         "message",
		UserID:       1,
		LanguageCode: "en",
	}
	err = c.Collect(event)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: assert on contents of sqlite database
}
