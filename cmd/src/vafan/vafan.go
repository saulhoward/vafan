package main

import (
    "flag"
    "vafan"
)

var start *bool = flag.Bool("start", false, "Start the server")

func main() {
    flag.Parse()
    if *start {
        print("Starting vafan...\n")
        vafan.StartServer()
    } else {
        print("No command. Quitting...\n")
    }
}
