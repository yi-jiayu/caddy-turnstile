package turnstile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("turnstile", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	cfg := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return Turnstile{Next: next}
	}
	cfg.AddMiddleware(mid)
	return nil
}

// Turnstile is a Caddy middleware which records incoming traffic to a
// downstream Telegram bot.
type Turnstile struct {
	Next httpserver.Handler
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

	if event := ExtractEvent(time.Now(), update); event != nil {
		fmt.Println(event)
	}

	r.Body = ioutil.NopCloser(&buf)
	return h.Next.ServeHTTP(w, r)
}
