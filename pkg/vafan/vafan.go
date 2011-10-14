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
    name string
}

// env identifies the environment type
type env int
const (
	dev env = iota
	staging
	production
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
    for {
        item := l.nextItem()
        if item.typ == itemText {
            if item.val == "dev" {
                e = dev
            }
            if item.val == "production" {
                e = production
            }
            if item.val == "staging" {
                e = staging
            }
            s.name = item.val
        }
        if item.typ == itemEnd ||
            item.typ == itemError ||
            item.typ == itemColon {
                break
            }
    }
    return s, e
}

func Route(ctx *web.Context, val string) string {
    r := parsePath(ctx.URL.Path)
    s, e := parseHost(ctx.Request.Host)
    out := ""
    for _, p := range r.parts {
        out += p
        out += " "
    }
    out += " "
    out += s.name
    if (e == dev) {
        out += " ... dev "
    }
    return out
    /* filename := path.Join(path.Join(os.Getenv("PWD"), "templates"), "index.html.mustache")
    return mustache.RenderFile(filename, map[string]string{"host":host.Name}) */
}

