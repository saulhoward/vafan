// Copyright 2012 Saul Howard. All rights reserved.

// Markdown text.

package vafan

type markdown string

// - renders as HTML when read as JSON
/* NOT WORKING???
func (m markdown) MarshalJSON() ([]byte, error) {
    b := []byte(m)
    html := blackfriday.MarkdownCommon(b)
    return html, nil
}
*/
