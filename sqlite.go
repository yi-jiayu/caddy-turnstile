package turnstile

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyfile"
)

type SQLiteCollector struct {
	db *sql.DB
}

func insertEvent(stmt *sql.Stmt, event Event) (sql.Result, error) {
	languageCode := sql.NullString{
		String: event.LanguageCode,
		Valid:  event.LanguageCode != "",
	}
	chatID := sql.NullInt64{
		Int64: event.ChatID,
		Valid: event.ChatID != 0,
	}
	chatType := sql.NullString{
		String: event.ChatType,
		Valid:  event.ChatType != "",
	}
	return stmt.Exec(event.Time, event.Type, event.UserID, languageCode, chatID, chatType)
}

func (c *SQLiteCollector) Collect(event Event) error {
	tx, err := c.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`insert into events values (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = insertEvent(stmt, event)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// Close closes the underlying SQLite database connection.
func (c *SQLiteCollector) Close() error {
	return c.db.Close()
}

func (c *SQLiteCollector) initDB() error {
	stmt := `create table if not exists events
(
  time          datetime not null,
  type          text     not null,
  user_id       int      not null,
  language_code text,
  chat_id       int,
  chat_type     text
)`
	_, err := c.db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func NewSQLiteCollector(f string) (*SQLiteCollector, error) {
	db, err := sql.Open("sqlite3", f)
	if err != nil {
		return nil, err
	}
	c := &SQLiteCollector{
		db: db,
	}
	caddy.RegisterEventHook("turnstile_sqlite_close", func(eventType caddy.EventName, eventInfo interface{}) error {
		if eventType == caddy.ShutdownEvent {
			log.Printf("[INFO] turnstile: sqlite: closing database")
			err := c.Close()
			if err != nil {
				log.Printf("[ERROR] turnstile: sqlite: failed to close database: %s", err)
			}
		}
		return nil
	})
	err = c.initDB()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func SQLiteCollectorFactory(d *caddyfile.Dispenser) (Collector, error) {
	var f string
	if d.NextArg() {
		f = d.Val()
	} else {
		return nil, d.ArgErr()
	}
	log.Printf("[INFO] turnstile: using sqlite collector (database file: %s)", f)
	return NewSQLiteCollector(f)
}
