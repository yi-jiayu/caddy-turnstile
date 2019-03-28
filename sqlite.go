package turnstile

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyfile"
)

type SQLiteCollector struct {
	db *sql.DB
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
	// TODO: store empty fields as nulls?
	_, err = stmt.Exec(event.Time, event.Type, event.UserID, event.LanguageCode, event.ChatID, event.ChatType)
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
