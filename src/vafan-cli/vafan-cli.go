/*
    Vafan CLI command
    
    Provides access to various vafan library functions.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"vafan"
)

var start *bool = flag.Bool("start", false, "Start the server")
var printTweets *bool = flag.Bool("tweets", false, "Get tweets")
var javascriptFiles *bool = flag.Bool("list-javascript-files", false, "Get a list of javascript library files.")
var cssFiles *string = flag.String("list-css-files", "", "Get a list of css files for a site name.")
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
		v.UpdateExternalData()
		os.Exit(0)
	case *javascriptFiles:
		j := vafan.GetJavascriptFileList()
		fmt.Fprintln(os.Stdout, j)
		os.Exit(0)
	case *cssFiles != "":
		c := vafan.GetCSSFileList(*cssFiles)
		fmt.Fprintln(os.Stdout, c)
		os.Exit(0)
	case *printTweets:
		//vafan.PrintTweets()
		os.Exit(0)
	}
	fmt.Fprintln(os.Stdout, "No command. Quitting...")
	os.Exit(1)
}
