package main

import (
    "os"
    "fmt"
    "flag"
    "vafan"
)

var start *bool = flag.Bool("start", false, "Start the server")
var makeUserAdmin *string = flag.String("make-admin", "", "Make this user an admin")
var getVideoDetails *string = flag.String("get-video-details", "", "Get external video details")

func main() {
    flag.Parse()
    switch {
    case *start:
        fmt.Fprintln(os.Stdout, "Starting vafan...")
        vafan.StartServer()
        os.Exit(0)
    case *makeUserAdmin != "":
        fmt.Fprintln(os.Stdout, "Making user admin...")
        vafan.MakeUserAdmin(*makeUserAdmin)
        os.Exit(0)
    case *getVideoDetails != "":
        fmt.Fprintln(os.Stdout, "Getting video details...")
        v, _ := vafan.GetVideoByName(*getVideoDetails)
        v.Youtube.FetchDetails()
        os.Exit(0)
    }
    fmt.Fprintln(os.Stdout, "No command. Quitting...")
    os.Exit(1)
}
