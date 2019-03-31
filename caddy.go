package turnstile

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

// Turnstile is a Caddy middleware which records incoming traffic to a
// downstream Telegram bot.
type Turnstile struct {
	collector Collector
	next      httpserver.Handler
}

func init() {
	caddy.RegisterPlugin("turnstile", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	_ = c.Next() // skip directive name
	var collectorName string
	if ok := c.NextArg(); ok {
		collectorName = c.Val()
	} else {
		return c.ArgErr()
	}

	var collector Collector
	var err error
	if collectorFactory, ok := collectors[collectorName]; ok {
		collector, err = collectorFactory(&c.Dispenser)
		if err != nil {
			return err
		}
	} else {
		return c.Errf(`turnstile: no such collector "%s"`, collectorName)
	}

	cfg := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return Turnstile{
			collector: collector,
			next:      next,
		}
	}
	cfg.AddMiddleware(mid)
	return nil
}

func (h Turnstile) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	var buf bytes.Buffer
	rdr := io.TeeReader(r.Body, &buf)

	var update Update
	err := json.NewDecoder(rdr).Decode(&update)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	r.Body.Close()

	go func(t time.Time, u Update) {
		if event := ExtractEvent(t, u); event != nil {
			err := h.collector.Collect(*event)
			if err != nil {
				log.Printf("[ERROR] turnstile: error collecting event: %s", err)
			}
		}
	}(time.Now(), update)

	r.Body = ioutil.NopCloser(&buf)
	return h.next.ServeHTTP(w, r)
}
