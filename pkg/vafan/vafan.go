/*
    This package provides a server framework.

    http://saulhoward.com/vafan

    @author saul@saulhoward.com
*/
package vafan

import (
    "github.com/hoisie/web.go"
    //"github.com/kless/goconfig/config"
)
type request struct {
    site site
    env env
    resource resource
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

func parsePath(p string) resource {
    l := lex(path, p)
    var r video
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

// 
func parseHost(h string) (s site, e env) {
    l := lex(host, h)
    var levels []string
    for {
        item := l.nextItem()
        if item.typ == itemText {
            levels = append(levels, item.val)
        }
        if item.typ == itemEnd ||
            item.typ == itemError ||
            item.typ == itemColon {
                break
            }
    }
    domain := levels[:]

    // Determine environment,
    // assuming dev.sitedomain or just sitedomain
    switch levels[0] {
    case "dev":
        e = dev
        domain = levels[1:]
        break
    case "staging":
        e = staging
        domain = levels[1:]
        break
    case "production":
        domain = levels[1:]
        e = production
        break
    default:
        e = production
        break
    }
    s.domain = ""
    first := true
    for _, d := range domain {
        if !first {
            s.domain += "."
        }
        first = false
        s.domain += d
    }
    //s.name = item.val
    return
}

func Route(ctx *web.Context, val string) (output string) {
    var req request
    // get the resource from the path
    req.resource = parsePath(ctx.URL.Path)
    // get the host and env from the host
    req.site, req.env = parseHost(ctx.Request.Host)
    // with the request, get the data
    output = get(req)
    return
}

