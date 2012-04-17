// Copyright 2012 Saul Howard. All rights reserved.

// CSS support.

package vafan

import (
	"fmt"
	"html/template"
)

const minifiedCSSFile = "/js/vafan.min.js"

var cssFiles = []string{
	"/css/bootstrap/bootstrap.css",
	"/css/bootstrap/bootstrap-responsive.css",
	"/css/font-awesome/font-awesome.css",
	"/css/datepicker/datepicker.css",
	"/css/style.css",
}

func getCSSFiles(siteName string) []string {
	return append(cssFiles, getSiteSpecificCSSFile(siteName))
}

func getSiteSpecificMinifiedCSSFile(siteName string) string {
	return "/css/" + siteName + ".min.css"
}

func getSiteSpecificCSSFile(siteName string) string {
	return "/css/" + siteName + ".css"
}

const ie6CSSFile = "/css/ie6.1.b.css"

func getCSSHTML(site *site, env string) template.HTML {
	var html string
	html = html + "<!--[if ! lte IE 6]><!-->" + "\n"
	switch env {
	case "dev":
		for _, c := range getCSSFiles(site.Name) {
			html = html + getCSSTagHTML(c) + "\n"
		}
	default:
		html = html + getCSSTagHTML(getSiteSpecificMinifiedCSSFile(site.Name))
	}
	html = html + "<!--<![endif]--><!--[if lte IE 6]>" + "\n"
	html = html + getCSSTagHTML(ie6CSSFile) + "\n"
	html = html + "<![endif]-->" + "\n"
	return template.HTML(html)
}

func getCSSTagHTML(css string) string {
	return fmt.Sprintf("<link rel=\"stylesheet\" href=\"%v\">", css)
}

func GetCSSFileList(siteName string) string {
	var list string
	for _, c := range getCSSFiles(siteName) {
		list = list + vafanConf.baseDir + "/static" + c + "\n"
	}
	return list
}
