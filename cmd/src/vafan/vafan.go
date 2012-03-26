package main

import (
    "flag"
    "vafan"
)

var start *bool = flag.Bool("start", false, "Start the server")
var makeUserAdmin *string = flag.String("make-admin", "", "Make this user an admin")

func main() {
    flag.Parse()
    switch {
    case *start:
        print("Starting vafan...\n")
        vafan.StartServer()
        return
    case *makeUserAdmin != "":
        print("Making user admin...\n")
        vafan.MakeUserAdmin(*makeUserAdmin)
        return
    }
    print("No command. Quitting...\n")
    return
}
