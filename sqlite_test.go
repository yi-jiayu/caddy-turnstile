package turnstile

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type eventWithNulls struct {
	Time         time.Time
	Type         string
	UserID       int64
	LanguageCode sql.NullString
	ChatID       sql.NullInt64
	ChatType     sql.NullString
}

func scanEvents(rows *sql.Rows) ([]eventWithNulls, error) {
	var events []eventWithNulls
	for rows.Next() {
		var e eventWithNulls
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
			Type:   "message",
			UserID: 2,
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
	want := []eventWithNulls{
		{
			Type:   "message",
			UserID: 1,
			LanguageCode: sql.NullString{
				String: "en",
				Valid:  true,
			},
			ChatID: sql.NullInt64{
				Int64: 1,
				Valid: true,
			},
			ChatType: sql.NullString{
				String: "private",
				Valid:  true,
			},
		},
		{
			Type:   "message",
			UserID: 2,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ExtractEvent() = %v, want %v", got, want)
	}
}
