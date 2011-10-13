type item struct {
    typ itemType // Type, such as itemSlash
    val string // Value such as "videos"
}

// itemType identifies the type of lex items
type itemType int

const (
    itemError itemType = iota // error occured.
    // value ios text of error the cursor, spelled '.'
    itemSlash // seperator for path '/'
    itemEnd
    itemGroup

)

type stateFn func(*lexer) stateFn



func run() {
    for state := startState; state != nil; {
        state = state(lexer)
    }
}


type lexer struct {
    name string
    input string
    start int
    pos int
    width int
    items chan item
}
