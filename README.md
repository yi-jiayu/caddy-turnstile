# turnstile
Caddy plugin for tracking incoming Telegram bot metrics

## Use case

I have a webhook-based Telegram bot running behind Caddy as a reverse proxy. Instead of writing code
to track usage at the application level, this plugin runs as a HTTP middleware that transparently
reads the content of incoming webhooks from Telegram to collect usage statistics.

## Backends

Currently, the only supported backend is SQLite. Collected events are written to an SQLite database.

## Dependencies

- git
- go 1.12 (required by Caddy v0.11.5)
- A working C toolchain (for cgo)

## Installation

Installing a Caddy plugins requires modifying the Caddy source and rebuilding Caddy. In addition,
SQLite requires `cgo`, so some modifications have to be made to the Caddy build process as well.

1. `go get` Caddy: `go get -u -v github.com/mholt/caddy`
2. Checkout Caddy `v0.11.5`: `cd $(go env GOPATH)/src/github.com/mholt/caddy && git checkout
   v0.11.5`
3. `go get` turnstile: `go get -u -v github.com/yi-jiayu/turnstile`
4. Apply `install_turnstile.patch` to Caddy: (from Caddy directory) `git apply $(go env
   GOPATH)/src/github.com/yi-jiayu/turnstile/install_turnstile.patch`
5. Build Caddy: (from Caddy directory) `cd caddy && go run build.go`
6. Add the following line to your site configuration in your Caddyfile: `turnstile sqlite <path>`,
   where `<path>` is the path for the SQLite database containing events.

