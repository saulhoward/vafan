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

// TODO: Sort out this mess!
/*


CHANGE THIS to be a specific function for fetching templates for a full
page, one that will only include

    page, header, footer etc.


THIS WOULD MEAN that if you want a 'one-off' template, they have to be
completely generic, that is, they can't have the dimensions

    eg, video, dvd, music

I GUESS that's OK then? It makes it simpler anyway...

*/
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
        path := getTemplatePath(file, format, res, site)
        print("\nAppend:  ")
        print(path)
        paths = append(paths, path)
    }
    var err error
	t := template.New("page.html")
    for _, path := range paths {
        t, err = t.ParseFiles(path)
        checkError(err)
    }
	//t, err := template.New("page.html").ParseFiles("/home/saul/code/vafan/templates/html/_anyResource/_anySite/page.html")
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
