// http://saulhoward.com/vafan
// @author saul@saulhoward.com

package vafan

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
	itemEnd
	itemGroup
)

const end = -1

type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name  string    // used only for error reports.
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

// lex creates a new scanner for the input string.
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
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

func lexText(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], slash) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexSlash // Next state.
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
	return lexText // Now back to text
}

