// Copyright 2012 Saul Howard. All rights reserved.

// Javascript support.

package vafan

import (
	"fmt"
	"html/template"
)

const minifiedJavascriptFile = "/js/vafan.min.js"

var javascriptFiles = [...]string{
	"/js/libs/underscore.js",
	"/js/libs/backbone.js",
	"/js/libs/bootstrap.js",
	"/js/libs/Three.js",
	"/js/libs/froogaloop.js",
	"/js/libs/bootstrap-datepicker.js",
	"/js/vafan/models/video.js",
	"/js/vafan/views/app.js",
	"/js/vafan/views/fonts.js",
	"/js/vafan/views/threeDeeDvd.js",
	"/js/vafan/views/twitter.js",
	"/js/vafan/views/video.js",
	"/js/vafan/views/brightonWokTrailer.js",
	"/js/global.js",
}

func getJavascriptLibraryHTML(site *site, env string) template.HTML {
	var html string
	switch env {
	case "dev":
		for _, j := range javascriptFiles {
			html = html + "\n" + getJavascriptTagHTML(j)
		}
	default:
		html = getJavascriptTagHTML(minifiedJavascriptFile)
	}
	return template.HTML(html)
}

func getJavascriptTagHTML(js string) string {
	return fmt.Sprintf("<script src=\"%v\"></script>", js)
}

func GetJavascriptFileList() string {
	var list string
	for _, j := range javascriptFiles {
		list = list + vafanConf.baseDir + "/static" + j + "\n"
	}
	return list
}
