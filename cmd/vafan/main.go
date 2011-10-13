package main

import (
    "saulhoward.com/vafan"
    "github.com/hoisie/web.go"
    /* "github.com/hoisie/mustache.go" */
    // "launchpad.net/mgo"
    /* "path" */
    /* "os" */
)

func routeRequest(ctx *web.Context, val string) string {
    r := vafan.Path(ctx.URL.Path)
    h := vafan.Host(ctx.Request.Host)
    s := ""
    for _, p := range r {
        s += p
        s += " "
    }
    s += " "
    s += h
    return s
    /* filename := path.Join(path.Join(os.Getenv("PWD"), "templates"), "index.html.mustache")
    return mustache.RenderFile(filename, map[string]string{"host":host.Name}) */
}

//func parse(filename string) string {
    //output := mustache.RenderFile(filename, map[string]string{"host": host})
//}

func main() {
    web.Get("/(.*)", routeRequest)
    web.Run("0.0.0.0:9999")
}
