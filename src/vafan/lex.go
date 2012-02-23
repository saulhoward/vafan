// http://saulhoward.com/vafan
// @author saul@saulhoward.com

package main

// This file provides the lexer for parsing hostnames and paths.

import (
	"fmt"
	"strings"
	"utf8"
)

// lex items
type item struct {
	typ itemType // Type, such as itemSlash
	val string   // Value such as "videos"
}

// pretty print items
func (i item) String() string {
	switch {
	case i.typ == itemEnd:
		return "END"
	case i.typ == itemError:
		return i.val
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items
type itemType int

const (
	itemError itemType = iota // error occured value is text of error
	itemSlash                 // seperator for path '/'
	itemText
	itemColon
	itemDot
	itemEnd
	itemPort
	itemGroup
)

// lexerTarget identifies the things being lexed
type lexerTarget int
const (
	path lexerTarget = iota
	host
)

const end = -1

type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
    target  lexerTarget    // the thign being lexed
    input string    // the string being scanned.
    state stateFn   // the next lexing function to enter
    start int       // start position of this item.
    pos   int       // current position in the input.
    width int       // width of last rune read from input.
    items chan item // channel of scanned items.
}

// next returns the next rune in the input.
func (l *lexer) next() (rune int) {
	if l.pos >= len(l.input) {
		l.width = 0
		return end
	}
	rune, l.width =
		utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item {
	for {
		select {
		case item := <-l.items:
			return item
		default:
			l.state = l.state(l)
		}
	}
	panic("not reached")
}


// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
    if strings.IndexRune(valid, l.next()) >= 0 {
        return true
    }
    l.backup()
    return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
    for strings.IndexRune(valid, l.next()) >= 0 {
    }
    l.backup()
}

// backup steps back one rune. Can only be called once per call of next.
func (l *lexer) backup() {
    l.pos -= l.width
}

// lex creates a new scanner for the input string.
func lex(target lexerTarget, input string) *lexer {
	l := &lexer{
		target:  target,
		input: input,
		state: lexText,
		items: make(chan item, 2), // Two items sufficient.
	}
	return l
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// state fucntions

const slash = "/"
const colon = ":"
const dot   = "."

func lexText(l *lexer) stateFn {
    switch l.target {
    case path:
        for {
            if strings.HasPrefix(l.input[l.pos:], slash) {
                if l.pos > l.start {
                    l.emit(itemText)
                }
                return lexSlash // Next state.
            }
            if l.next() == end {
                break
            }
        }
    case host:
        for {
            if strings.HasPrefix(l.input[l.pos:], dot) {
                if l.pos > l.start {
                    l.emit(itemText)
                }
                return lexDot // Next state.
            }
            if strings.HasPrefix(l.input[l.pos:], colon) {
                if l.pos > l.start {
                    l.emit(itemText)
                }
                return lexColon // Next state.
            }
            if l.next() == end {
                break
            }
        }
    }
    // Correctly reached End.
    if l.pos > l.start {
        l.emit(itemText)
    }
    l.emit(itemEnd) // Useful to make End a token.
    return nil      // Stop the run loop.
}

func lexSlash(l *lexer) stateFn {
    l.pos += len(slash) // move past it
    l.emit(itemSlash)
    return lexText // Now back to text
}

func lexColon(l *lexer) stateFn {
	l.pos += len(colon) // move past it
	l.emit(itemColon)
    if l.target == host {
        return lexPort // Move onto port no.
    }
	return lexText // Or back to text
}

func lexDot(l *lexer) stateFn {
	l.pos += len(dot) // move past it
	l.emit(itemDot)
	return lexText // Now back to text
}

func lexPort(l *lexer) stateFn {
    digits := "0123456789"
    l.acceptRun(digits)
    if l.next() == end {
        // Correctly reached End.
        if l.pos > l.start {
            l.emit(itemPort)
        }
        l.emit(itemEnd) // Useful to make End a token.

    } else {
        l.emit(itemError) // Useful to make End a token.
    }
    return nil
}


