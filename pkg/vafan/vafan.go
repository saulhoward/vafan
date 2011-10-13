// This package provides a server framework.
// http://saulhoward.com/vafan
// @author saul@saulhoward.com

package vafan

type resource struct {
    parts []item
}

type host struct {
    name string
}

func Path(path string) []string {

    l := lex("path", path)

    /* r := resource{} */

    parts := make([]string, 1)

    for {
        item := l.nextItem()
        if item.typ == itemText {
            /* parts = append(r.parts, item) */
            parts = append(parts, item.val)
        }
        if item.typ == itemEnd || item.typ == itemError {
              break
        }
    }
    return parts
}

func Host(host string) string {
    l := lex("host", host)
    var h string
    for {
        item := l.nextItem()
        if item.typ == itemText {
            h = item.val
        }
        if item.typ == itemEnd || item.typ == itemError || item.typ == itemColon {
              break
        }
    }
    return h
}
