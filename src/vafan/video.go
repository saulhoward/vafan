// Vafan - a web server for Convict Films
//
// Videos
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
    "fmt"
	"errors"
	"launchpad.net/mgo"
	"launchpad.net/mgo/bson"
	"regexp"
)

var ErrVideoNotFound = errors.New("video: doesn't exist")

type video struct {
	Id          string
	Name        string
	Title       string
	Description markdown
	Sites       []*site // the sites that display this vid
	Youtube     youtubeVideo
	Vimeo       vimeoVideo
}

// -- markdown type 

type markdown string

// - renders as HTML when read as JSON
/* NOT WORKING???
func (m markdown) MarshalJSON() ([]byte, error) {
    b := []byte(m)
    html := blackfriday.MarkdownCommon(b)
    return html, nil
}
*/

// -- 

// External video types, youtube, vimeo
type externalVideo interface {
	FetchDetails() (err error)
}

type vimeoVideo struct {
	Id string
}

// --

func GetVideoByName(name string) (v *video, err error) {
	v = new(video)
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to dial db (Mongo): %v", err))
		return
	}
	defer session.Close()
	c := session.DB("vafan").C("videos")
	err = c.Find(bson.M{"name": name}).One(v)
	if err != nil {
		if err == mgo.NotFound {
			err = ErrVideoNotFound
			return
		}
		_ = logger.Err(fmt.Sprintf("Failed to get video (Mongo): %v", err))
		return
	}
	return
}

func (v *video) save() (err error) {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to dial db (Mongo): %v", err))
		return
	}
	defer session.Close()
	c := session.DB("vafan").C("videos")
	err = c.Insert(v)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to insert video (Mongo): %v", err))
		return
	}
	return
}

// name must be unicode alphanumericals and dashes only
func (v *video) isNameLegal() bool {
	var illegalCharsRe = regexp.MustCompile(`[^\-\p{L}\p{M}\p{N}]+`)
	if illegalCharsRe.MatchString(v.Name) {
		return false
	}
	return true
}
