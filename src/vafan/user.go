// Vafan - a web server for Convict Films
//
// User
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
    "crypto/rand"
    "fmt"
    "io"
    "log"
	"regexp"
	"strings"
    "database/sql"
    _ "github.com/ziutek/mymysql/godrv"
)

// -- DB

func connectDb() *sql.DB {
    db, err := sql.Open("mymysql", "vafan/user/pass")
    if err != nil {
        panic("Error connecting to mysql db: " + err.Error())
    }
    return db
}

// -- User stuff

type User struct {
    Id           string
    Username     string
    EmailAddress string
    Password     string
    Role         string
}

const defaultUserRole = "user"

func NewUser() *User {
    u := User{uuid(), "", "", "", defaultUserRole}
    return &u
}

// Use UUID v4 as user IDs
func uuid() string {
    b := make([]byte, 16)
    _, err := io.ReadFull(rand.Reader, b)
    if err != nil {
        log.Fatal(err)
    }
    b[6] = (b[6] & 0x0F) | 0x40
    b[8] = (b[8] &^ 0x40) | 0x80
    return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (u *User) save() {
    db := connectDb()
    defer db.Close()
    query := `insert into users values (?, ?, ?, ?, ?)`
    stmt, err := db.Prepare(query)
    if err != nil {
        panic(err)
    }
    _, err = stmt.Exec(u.Id, u.Username, u.EmailAddress, u.Password, u.Role)
    if err != nil {
        panic(err)
    }
    return
}

func (u *User) isLegal() bool {
    if !u.isUsernameLegal() ||
    !u.isEmailAddressLegal() ||
    !u.isPasswordLegal() {
        return false
    }
    return true
}

func (u *User) isUsernameLegal() bool {
	var re = regexp.MustCompile(`[^\d+|\w+]`)
	if re.MatchString(u.Username) {
		return false
	}
	return true
}

func (u *User) isEmailAddressLegal() bool {
	if strings.Contains(u.EmailAddress, "@") {
		return true
	}
	return false
}

func (u *User) isPasswordLegal() bool {
	if len(u.Password) < 6 {
		return false
	}
	return true
}
