package vafan

type resource struct {
    parts []item
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

func Host(hostname string) string {
    return hostname
}


