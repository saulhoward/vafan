package main

import (
    "saulhoward.com/vafan"
    //"github.com/hoisie/web.go"
    /* "github.com/hoisie/mustache.go" */
    //"launchpad.net/mgo"
    /* "path" */
    /* "os" */
)

/* type Host struct { */
    /* Name string */
    /* Title string */
/* } */

/* type Resource struct { */
    /* Name string */
/* } */

func routeRequest(ctx *web.Context, val string) string {

    //host := getHost(ctx.Request.Host)
    resource := vafan.Resource(ctx.URL.Path)
    return resource.Name

    /* filename := path.Join(path.Join(os.Getenv("PWD"), "templates"), "index.html.mustache")
    return mustache.RenderFile(filename, map[string]string{"host":host.Name}) */
}

/* func getHost(hostname string) *Host { */
    /* // TODO: get the host details from the DB (or config?) and fill this in properly */
    /* h := Host{hostname, hostname} */
    /* return &h */
/* } */

//func parse(filename string) string {
    //output := mustache.RenderFile(filename, map[string]string{"host": host})
//}

func main() {
    web.Get("/(.*)", routeRequest)
    web.Run("0.0.0.0:9999")
}
