package turnstile

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func scanEvents(rows *sql.Rows) ([]Event, error) {
	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.Time, &e.Type, &e.UserID, &e.LanguageCode, &e.ChatID, &e.ChatType)
		if err != nil {
			return events, err
		}
		events = append(events, e)
	}
	return events, nil
}

func TestSQLiteCollector(t *testing.T) {
	c, err := NewSQLiteCollector(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	events := []Event{
		{
			Type:         "message",
			UserID:       1,
			LanguageCode: "en",
			ChatID:       1,
			ChatType:     "private",
		},
		{
			Type:         "message",
			UserID:       2,
			LanguageCode: "en",
		},
	}
	for _, event := range events {
		err = c.Collect(event)
		if err != nil {
			t.Fatal(err)
		}
	}

	rows, err := c.db.Query(`select * from events`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	got, err := scanEvents(rows)
	if err != nil {
		t.Fatal(err)
	}
	want := events
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtractEvent() = %v, want %v", got, want)
	}
}

func TestSQLiteCollectorFactory(t *testing.T) {

}
