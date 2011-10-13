package vafan

type resource struct {
    parts []string
}

func ParseResource(path string) {
    l := lex("path", path)
    items := []item
    r := resource
    for {
        item := l.nextItem()
        r.parts = append(r.parts, item)
        if item.typ == itemEnd || item.typ == itemError {
              break
        }
    }
}

func (r *resource)Root() {
    if r.parts

}
