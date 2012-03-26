// Vafan - a web server for Convict Films
//
// Videos
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
    "regexp"
    "errors"
    "launchpad.net/mgo/bson"
    "launchpad.net/mgo"
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

type youtubeVideo struct {
    Id string
}

type vimeoVideo struct {
    Id string
}

func getVideoByName(name string) (v *video, err error) {
    session, err := mgo.Dial("127.0.0.1")
    if err != nil {
        panic(err)
    }
    defer session.Close()
    c := session.DB("vafan").C("videos")
    v = new(video)
    err = c.Find(bson.M{"name": name}).One(v)
    if err != nil {
        if err == mgo.NotFound {
            err = ErrVideoNotFound
            return
        }
        checkError(err)
    }
    return
}

func (v *video) save() (err error) {
    session, err := mgo.Dial("127.0.0.1")
    if err != nil {
        panic(err)
    }
    defer session.Close()
    c := session.DB("vafan").C("videos")
    err = c.Insert(v)
    if err != nil {
        panic(err)
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
