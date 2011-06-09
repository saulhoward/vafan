package main

import (
    "github.com/hoisie/web.go"
    "github.com/kless/goconfig/config"
)

func route_requests(ctx *web.Context, val string) string {
    host := ctx.Request.Host
    return host
}

func main() {
    web.Get("/(.*)", route_requests)
    web.Run("0.0.0.0:9999")
}
