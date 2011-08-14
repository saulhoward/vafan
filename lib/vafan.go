package main

import (
    "github.com/fiber/web.go"
    "github.com/hoisie/mustache.go"
    "path"
    "os"
)

func route_requests(ctx *web.Context, val string) string {
    host := ctx.Request.Host
    filename := path.Join(path.Join(os.Getenv("PWD"), "templates"), "index.html.mustache")
    output := mustache.RenderFile(filename, map[string]string{"host": host})
    return output
}

func main() {
    web.Get("/(.*)", route_requests)
    web.Run("127.0.0.1:9999")
}
