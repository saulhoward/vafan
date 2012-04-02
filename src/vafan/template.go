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
	"fmt"
    "strings"
    "reflect"
	"path/filepath"
	"html/template"
    "github.com/russross/blackfriday"
)

// Strip the 'Resource' from the name
func resourceDirName(res Resource) string {
    return strings.Replace(resourceName(res), "Resource", "", 1)
}

func getTemplatePath(file string, format string, res Resource, s *site) string {
    var baseDir, err = conf.String("default", "base-dir")
    if err != nil {
        // fatal error?
        fmt.Fprintf(os.Stderr, "Failed reading 'base-dir' from configuration: %v", err)
        os.Exit(1)
    }

    //_ = logger.Info(fmt.Sprintf("Looking for template format: %v and resource: %v", format, resourceDirName(res)))

    //Check for the most specific template first
    checkFormat := format
    checkResName := resourceDirName(res)
    checkSite := s.Name
    for i:= 0; templateExists(baseDir, file, checkFormat, checkResName, checkSite) == false; i++ {
		if i == 0 {
			checkSite = "_anySite"
		} else if i == 1 {
			checkResName = "_anyResource"
            checkSite = s.Name
		} else if i == 2 {
            checkResName = resourceDirName(res)
			checkSite = "_anySite"
		} else if i == 3 {
			checkResName = "_anyResource"
			checkSite = "_anySite"
		} else if i == 4 {
			checkFormat = "_anyFormat"
		} else if i > 4 {
            _ = logger.Err("Failed to find a matching template! This is bad!")
            // if you ever get here, come back and rewrite this horrible
            // func
            os.Exit(1)
		}
	}
    return filepath.Join(baseDir, "templates", checkFormat, checkResName, checkSite, file)
}

func getPageTemplate(format string, res Resource, s *site) (t *template.Template, err error) {
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
        paths = append(paths, getTemplatePath(file, format, res, s))
    }
	t, err = template.New("page.html").
        Funcs(template.FuncMap{"eq": reflect.DeepEqual, "markdown": markdownToHtml }).
        ParseFiles(paths...)
    if err != nil {
        _ = logger.Err(fmt.Sprintf("Failed to create template: %v", err))
        return
    }
	return
}

// Uses black friday library to convert markdown to html, this is
// assumed safe
func markdownToHtml(md markdown) template.HTML {
    bmd := []byte(md)
    bhtml := blackfriday.MarkdownCommon(bmd)
    return template.HTML(bhtml)
}

func templateExists(baseDir string, file string, format string, resName string, siteName string) bool {
	path := filepath.Join(baseDir, "templates", format, resName, siteName, file)
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
