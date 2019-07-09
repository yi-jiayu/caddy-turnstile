package main

import (
	"github.com/caddyserver/caddy/caddy/caddymain"

	_ "github.com/yi-jiayu/turnstile"
)

func main() {
	caddymain.Run()
}
