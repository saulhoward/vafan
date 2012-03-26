// Vafan - a web server for Convict Films
//
// Videos
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
    "errors"
    "launchpad.net/mgo/bson"
    "launchpad.net/mgo"
)

var ErrVideoNotFound = errors.New("video: doesn't exist")

type video struct {
    Id      string
    Name    string
    Title   string
    Description   string
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
    /* c := session.DB("vafan").C("videos") */
    /* err = c.Insert(&video{"Ale", "+55 53 8116 9639"}, */
    /* &Person{"Cla", "+55 53 8402 8510"}) */
    /* if err != nil { */
        /* panic(err) */
    /* } */

    return
}
