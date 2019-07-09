package turnstile

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

// Turnstile is a Caddy middleware which records incoming traffic to a
// downstream Telegram bot.
type Turnstile struct {
	collector Collector
	next      httpserver.Handler
	wg        *sync.WaitGroup
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
		return New(collector, next)
	}
	cfg.AddMiddleware(mid)
	return nil
}

func (h Turnstile) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.Method != http.MethodPost {
		return h.next.ServeHTTP(w, r)
	}

	var buf bytes.Buffer
	rdr := io.TeeReader(r.Body, &buf)

	var update Update
	err := json.NewDecoder(rdr).Decode(&update)
	// restore the response body
	// note: this will only contain whatever was read by the call to Decode
	r.Body = ioutil.NopCloser(&buf)
	if err != nil {
		log.Printf("[WARNING] turnstile: error decoding json: %s", err)
		return h.next.ServeHTTP(w, r)
	}
	r.Body.Close()

	h.wg.Add(1)
	go func(t time.Time, u Update) {
		defer h.wg.Done()
		if event := ExtractEvent(t, u); event != nil {
			err := h.collector.Collect(*event)
			if err != nil {
				log.Printf("[ERROR] turnstile: error collecting event: %s", err)
			}
		}
	}(time.Now(), update)

	return h.next.ServeHTTP(w, r)
}

func New(c Collector, next httpserver.Handler) Turnstile {
	return Turnstile{
		collector: c,
		next:      next,
		wg:        new(sync.WaitGroup),
	}
}
