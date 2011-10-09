package main

import (
    "github.com/fiber/web.go"
    "github.com/hoisie/mustache.go"
    //"launchpad.net/mgo"
    "path"
    "os"
)

type Host struct {
    Name string
    Title string
}

type Resource struct {
    Name string
}

func routeRequest(ctx *web.Context, val string) string {

    host := getHost(ctx.Request.Host)
    /*
    URLPath := ctx.URL.Path
    var resource Resource
    switch {
    case URLPath == "/":
        resource = Resource{"index"}
    default:
        resource = Resource{"404"}
    }
    */

    filename := path.Join(path.Join(os.Getenv("PWD"), "templates"), "index.html.mustache")
    return mustache.RenderFile(filename, map[string]string{"host":host.Name})
}

func getHost(hostname string) *Host {
    // TODO: get the host details from the DB (or config?) and fill this in properly
    h := Host{hostname, hostname}
    return &h
}

//func parse(filename string) string {
    //output := mustache.RenderFile(filename, map[string]string{"host": host})
//}

func main() {
    web.Get("/(.*)", routeRequest)
    web.Run("0.0.0.0:9999")
}
