# turnstile
Caddy plugin for tracking incoming Telegram bot metrics

## Use case

I have a webhook-based Telegram bot running behind Caddy as a reverse proxy. Instead of writing code
to track usage at the application level, this plugin runs as a HTTP middleware that transparently
reads the content of incoming webhooks from Telegram to collect usage statistics.

## Backends

Currently, the only supported backend is SQLite. Collected events are written to an SQLite database.

## Installation

Installing a Caddy plugins requires modifying the Caddy source and rebuilding Caddy. In addition,
SQLite requires `cgo`, so some modifications have to be made to the Caddy build process as well.

1. Apply `install_turnstile.patch` to Caddy 1.11.5
2. Build Caddy
3. Add the following line to your site configuration in your Caddyfile: `turnstile sqlite <path>`,
   where `<path>` is the path for the SQLite database containing events.

