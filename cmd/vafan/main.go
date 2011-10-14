package main

import (
    "saulhoward.com/vafan"
    "github.com/hoisie/web.go"
    /* "github.com/hoisie/mustache.go" */
    // "launchpad.net/mgo"
    /* "path" */
    /* "os" */
)
//func parse(filename string) string {
    //output := mustache.RenderFile(filename, map[string]string{"host": host})
//}

func main() {
    web.Get("/(.*)", vafan.Route)
    web.Run("0.0.0.0:9999")
}
