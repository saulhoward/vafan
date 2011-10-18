// This package provides a server framework.
// http://saulhoward.com/vafan
// @author saul@saulhoward.com

package vafan

import (
    "github.com/hoisie/web.go"
)

type resource struct {
    parts []string
}

type site struct {
    domain string
    name string
}

// env identifies the environment type
type env int
const (
	production env = iota
	staging
    dev
)

func parsePath(p string) (r resource) {
    l := lex(path, p)
    for {
        item := l.nextItem()
        if item.typ == itemText {
            r.parts = append(r.parts, item.val)
        }
        if item.typ == itemEnd || item.typ == itemError {
            break
        }
    }
    return r
}

func parseHost(h string) (s site, e env) {
    l := lex(host, h)
    var texts []string
    for {
        item := l.nextItem()
        if item.typ == itemText {
            texts = append(texts, item.val)
        }
        if item.typ == itemEnd ||
            item.typ == itemError ||
            item.typ == itemColon {
                break
            }
    }
    domain := texts[:]

    // Determine environment,
    // assuming dev.sitedomain or just sitedomain
    switch texts[0] {
    case "dev":
        e = dev
        domain = texts[1:]
        break
    case "staging":
        e = staging
        domain = texts[1:]
        break
    case "production":
        domain = texts[1:]
        e = production
        break
    default:
        e = production
        break
    }
    s.domain = ""
    // Determine site, look up domain in config
    for _, d := range domain {
        s.domain += d
    }
    //s.name = item.val
    return
}

func Route(ctx *web.Context, val string) (out string) {
    r := parsePath(ctx.URL.Path)
    s, e := parseHost(ctx.Request.Host)

    out = ""
    for _, p := range r.parts {
        out += p
        out += " "
    }
    out += " "
    out += s.name
    out += "/ "
    out += s.domain
    out += "\\ "
    if (e == dev) {
        out += " ... dev "
    }
    return
    /* filename := path.Join(path.Join(os.Getenv("PWD"), "templates"), "index.html.mustache")
    return mustache.RenderFile(filename, map[string]string{"host":host.Name}) */
}

