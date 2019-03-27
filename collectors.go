package turnstile

import (
	"github.com/mholt/caddy/caddyfile"
)

var collectors = map[string]CollectorFactory{
	"sqlite": SQLiteCollectorFactory,
}

// Collector is the interface implemented by event consumers.
type Collector interface {
	Collect(Event) error
}

type CollectorFactory func(d *caddyfile.Dispenser) (Collector, error)
