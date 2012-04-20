/*
    Vafan Server command
    
    Starts the vafan server.
*/
package main

import (
	"os"
	"vafan"
)

func main() {
    vafan.StartServer()
    os.Exit(0)
}
