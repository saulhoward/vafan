// Vafan - a web server for Convict Films
//
// User
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
    "fmt"
    "io"
    "log"
	"errors"
	"regexp"
	"strings"
    "hash"
    "net/http"
    "crypto/hmac"
    "crypto/sha1"
    "crypto/rand"
    "database/sql"
    _ "github.com/ziutek/mymysql/godrv"
)

var ErrWrongPassword = errors.New("user: password fail")

// -- DB

func connectDb() *sql.DB {
    db, err := sql.Open("mymysql", "vafan/root/password")
    if err != nil {
        panic("Error connecting to mysql db: " + err.Error())
    }
    return db
}

func createSalt() string {
    b := make([]byte, 16)
    _, err := io.ReadFull(rand.Reader, b)
    if err != nil {
        log.Fatal(err)
    }
    return string(b)
}

func hashPassword(password string, salt string) string {
    var h hash.Hash = hmac.New(sha1.New, []byte(salt))
    h.Write([]byte(password))
    return string(h.Sum(nil))
}

// -- User stuff

func getUserByUsername(username string) *User {
    db := connectDb()
    defer db.Close()
    selectUser, err := db.Prepare(`select id, username, emailAddress, role, passwordHash, salt from users where username=?`)
    if err != nil {
        panic(err)
    }
    u := NewUser()
    err = selectUser.QueryRow(username).Scan(&u.Id, &u.Username, &u.EmailAddress, &u.Role, &u.passwordHash, &u.salt)
    if err != nil {
        panic(err)
    }
    return u
}

func getUserByEmailAddress(emailAddress string) *User {
    db := connectDb()
    defer db.Close()
    selectUser, err := db.Prepare(`select id, username, emailAddress, role, passwordHash, salt from users where emailAddress=?`)
    if err != nil {
        panic(err)
    }
    u := NewUser()
    err = selectUser.QueryRow(emailAddress).Scan(&u.Id, &u.Username, &u.EmailAddress, &u.Role, &u.passwordHash, &u.salt)
    if err != nil {
        panic(err)
    }
    return u
}

func getUserById(id string) *User {
    db := connectDb()
    defer db.Close()
    selectUser, err := db.Prepare(`select id, username, emailAddress, role, passwordHash, salt from users where id=?`)
    if err != nil {
        panic(err)
    }
    u := NewUser()
    err = selectUser.QueryRow(id).Scan(&u.Id, &u.Username, &u.EmailAddress, &u.Role, &u.passwordHash, &u.salt)
    if err != nil {
        panic(err)
    }
    return u
}

type User struct {
    Id           string // UUID v4 with dashes
    Username     string
    EmailAddress string
    Role         string
    URL          string
    passwordHash string
    salt         string
    IsLoggedIn   bool
}

func (u *User) getURL(r *http.Request) string {
    return resources["users"].URL(r, []string{"id", u.Id}).String()
}

const defaultUserRole = "user"

// brand new user, freshly minted id
func NewUser() *User {
    u := User{newUUID(), "", "", defaultUserRole, "", "", "", false }
    return &u
}

// just wants an id, the simplest form of user
func GetUser(id string) *User {
    u := User{id, "", "", defaultUserRole, "", "", "", false}
    return &u
}

// needs a map of user properties, id must be set
func getUserForUserInfo(userInfo map[string]string) (u *User, err error) {
    if userInfo["Id"] == "" {
        err = errors.New("Id must be set")
        return
    }
    newU := User{userInfo["Id"], userInfo["Username"], userInfo["EmailAddress"], userInfo["Role"], "", "", "",  false}
    return &newU, err
}

// Use UUID v4 as user IDs
func newUUID() string {
    b := make([]byte, 16)
    _, err := io.ReadFull(rand.Reader, b)
    if err != nil {
        log.Fatal(err)
    }
    b[6] = (b[6] & 0x0F) | 0x40
    b[8] = (b[8] &^ 0x40) | 0x80
    return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func (u *User) save(password string) error {
    u.salt = createSalt()
    u.passwordHash = hashPassword(password, u.salt)
    db := connectDb()
    defer db.Close()
    query := `insert into users values (?, ?, ?, ?, ?, ?)`
    stmt, err := db.Prepare(query)
    if err != nil {
        panic(err)
    }
    _, err = stmt.Exec(u.Id, u.Username, u.EmailAddress, u.passwordHash, u.salt, u.Role)
    if err != nil {
        return err
    }
    return nil
}

func (u *User) isLegal(password string) bool {
    if !u.isUsernameLegal() ||
    !u.isEmailAddressLegal() ||
    !u.isPasswordLegal(password) {
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

func (u *User) isRegistered() bool {
    db := connectDb()
    defer db.Close()
    selectUser, err := db.Prepare(`select id from users where id=?`)
    if err != nil {
        panic(err)
    }
    var id int
    err = selectUser.QueryRow(u.Id).Scan(&id)
    if err == sql.ErrNoRows {
        return false
    }
    return true
}

func (u *User) isUsernameNew() bool {
    db := connectDb()
    defer db.Close()
    selectUser, err := db.Prepare(`select username from users where username=?`)
    if err != nil {
        panic(err)
    }
    var username string
    err = selectUser.QueryRow(u.Username).Scan(&username)
    if err == sql.ErrNoRows {
        return true
    }
    return false
}

func (u *User) isEmailAddressLegal() bool {
	if strings.Contains(u.EmailAddress, "@") {
		return true
	}
	return false
}

func (u *User) isEmailAddressNew() bool {
    db := connectDb()
    defer db.Close()
    selectUser, err := db.Prepare(`select emailAddress from users where emailAddress=?`)
    if err != nil {
        panic(err)
    }
    var emailAddress string
    err = selectUser.QueryRow(u.EmailAddress).Scan(&emailAddress)
    if err == sql.ErrNoRows {
        return true
    }
    return false
}

func (u *User) isPasswordLegal(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}

func (u *User) isNew() bool {
    if u.isRegistered() || !u.isUsernameNew() || !u.isEmailAddressNew() {
		return false
	}
	return true
}

func (u *User) setLoggedIn() {
    u.IsLoggedIn = true
}

// Login
func login(usernameOrEmailAddress string, password string) (u *User, err error) {
    // confirm username and or email exists, get user
    u = NewUser()
    err = nil
    u.Username = usernameOrEmailAddress
    if !u.isUsernameNew() {
        u = getUserByUsername(usernameOrEmailAddress)
    } else {
        u.EmailAddress = usernameOrEmailAddress
        if !u.isEmailAddressNew() {
            u = getUserByEmailAddress(usernameOrEmailAddress)
        }
    }
    // confirm that the user's password is correct
    if hashPassword(password, u.salt) == u.passwordHash {
        return
    }
    err = ErrWrongPassword
    return
}

