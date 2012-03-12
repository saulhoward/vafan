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
    //print("\n")
    //print("\nLooking for template with " + format + " " + res.name + " " + site)
	//Check for the most specific template first
    checkFormat := format
    checkResName := res.name
    checkSite := site
    for i:= 0; templateExists(file, checkFormat, checkResName, checkSite) == false; i++ {
		if i == 0 {
			checkSite = "_anySite"
		} else if i == 1 {
			checkResName = "_anyResource"
            checkSite = site
		} else if i == 2 {
            checkResName = res.name
			checkSite = "_anySite"
		} else if i == 3 {
			checkResName = "_anyResource"
			checkSite = "_anySite"
		} else if i == 4 {
			checkFormat = "_anyFormat"
		} else if i > 4 {
			// error checking here pls
			os.Exit(1)
		}
	}
    return filepath.Join(baseDir, "templates", checkFormat, checkResName, checkSite, file)
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

func templateExists(file string, format string, resName string, site string) bool {
	path := filepath.Join(baseDir, "templates", format, resName, site, file)
	_, err := os.Stat(path)
	if err != nil {
        /* print("\nNope:  ") */
        /* print(path) */
		return false
	}
    /* print("\nFound: ") */
    /* print(path) */
	return true
}
