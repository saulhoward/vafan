// Vafan - a web server for Convict Films
//
// Templating Functions
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	"os"
	//"log"
	"path/filepath"
	"html/template"
)

func getTemplatePath(file string, format string, res *resource, site string) string {
	//Check for the most specific template first
    checkFormat := format
    checkRes := res
    checkSite := site
    for i:= 0; templateExists(file, checkFormat, checkRes, checkSite) == false; i++ {
		if i == 0 {
			checkSite = "_anySite"
		} else if i == 1 {
			checkRes.name = "_anyResource"
            checkSite = site
		} else if i == 2 {
			checkRes.name = "_anyResource"
			checkSite = "_anySite"
		} else if i == 3 {
			checkFormat = "_anyFormat"
		} else if i > 2 {
			// error checking here pls
			os.Exit(1)
		}
	}
    res = checkRes
    format = checkFormat
    site = checkSite
    return filepath.Join(baseDir, "templates", format, res.name, site, file)
}

func getPageTemplate(format string, res *resource, site string) *template.Template {
    // Templates that make up a page
    // See http://www.w3.org/WAI/PF/aria/roles#landmark_roles 
    var tmplFiles = [...]string{
        "page.html",
        "banner.html",
        "navigation.html",
        "search.html",
        "main.html",
        "complementary.html",
        "contentinfo.html",
    }
    var paths = make([]string, 0)
    for _, file := range tmplFiles {
        paths = append(paths, getTemplatePath(file, format, res, site))
    }
	t, err := template.New("page.html").ParseFiles(paths...)
    checkError(err)
	return t
}

func templateExists(file string, format string, res *resource, site string) bool {
	path := filepath.Join(baseDir, "templates", format, res.name, site, file)
	_, err := os.Stat(path)
	if err != nil {
        print("\nNope:  ")
        print(path)
		return false
	}
    print("\nFound: ")
    print(path)
	return true
}
