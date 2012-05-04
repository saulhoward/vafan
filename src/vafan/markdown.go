// Copyright 2012 Saul Howard. All rights reserved.

// Markdown text.

package vafan

import (
	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
)

type Markdown string

func (m Markdown) HTML() (h template.HTML) {
	b := []byte(m)
	h = template.HTML(blackfriday.MarkdownCommon(b))
	return
}

// When marshalling JSON from a Markdown type, it will be of the
// form:
//     {
//          "MarkdownType:" {
//              "markdown": "content here",
//              "html":     "<p>content here</p>"
//          }
//     }
func (m Markdown) MarshalJSON() (j []byte, err error) {
	mMap := map[string]string{
		"markdown": string(m),
		"html":     string(m.HTML()),
	}
	j, err = json.Marshal(mMap)
	return
}

// When unmarshalling JSON into a Markdown type, it should be of the
// form:
//     {"MarkdownType:" {"markdown": "content here"}}
// Used, eg, when POSTing JSON.
// TODO: accept "html" property also, and convert.
func (m *Markdown) UnmarshalJSON(j []byte) (err error) {
	var mMap map[string]string
	err = json.Unmarshal(j, &mMap)
	if err != nil {
		logger.Err(fmt.Sprintf("Error when unmarshalling Markdown: %v", err))
	}
	*m = Markdown(mMap["markdown"])
	return
}
