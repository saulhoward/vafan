// Vafan - a web server for Convict Films
//
// User
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/ziutek/mymysql/godrv"
	"hash"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var ErrWrongPassword = errors.New("user: password fail")

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
	userRes := usersResource{u}
	return userRes.URL(r, nil).String()
}

const defaultUserRole = "user"

// brand new user, freshly minted id
func NewUser() *User {
	u := User{newUUID(), "", "", defaultUserRole, "", "", "", false}
	return &u
}

// -- DB

func connectSQLDB() *sql.DB {
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
		_ = logger.Err(fmt.Sprintf("Failed to create random salt: %v", err))
		return "defaultSalt"
	}
	return string(b)
}

func hashPassword(password string, salt string) string {
	var h hash.Hash = hmac.New(sha1.New, []byte(salt))
	h.Write([]byte(password))
	return string(h.Sum(nil))
}

// -- User stuff

func getUserByUsername(username string) (u *User) {
	u = NewUser()
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select id, username, emailAddress, role, passwordHash, salt from users where username=?`)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return
	}
	err = selectUser.QueryRow(username).Scan(&u.Id, &u.Username, &u.EmailAddress, &u.Role, &u.passwordHash, &u.salt)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
		return
	}
	return
}

func getUserByEmailAddress(emailAddress string) (u *User) {
	u = NewUser()
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select id, username, emailAddress, role, passwordHash, salt from users where emailAddress=?`)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return
	}
	err = selectUser.QueryRow(emailAddress).Scan(&u.Id, &u.Username, &u.EmailAddress, &u.Role, &u.passwordHash, &u.salt)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
		return
	}
	return
}

func getUserById(id string) (u *User) {
	u = NewUser()
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select id, username, emailAddress, role, passwordHash, salt from users where id=?`)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return
	}
	err = selectUser.QueryRow(id).Scan(&u.Id, &u.Username, &u.EmailAddress, &u.Role, &u.passwordHash, &u.salt)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
		return
	}
	return
}

// just wants an id, the simplest form of user
func GetUser(id string) *User {
	u := User{id, "", "", defaultUserRole, "", "", "", false}
	return &u
}

// needs a map of user properties, id must be set
func getUserForUserInfo(userInfo map[string]string) (u *User, err error) {
	if userInfo["Id"] == "" {
		err = errors.New("User: ID must be set")
		return
	}
	newU := User{userInfo["Id"], userInfo["Username"], userInfo["EmailAddress"], userInfo["Role"], "", "", "", false}
	return &newU, err
}

// Use UUID v4 as user IDs
func newUUID() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to create random uuid: %v", err))
		return "00000000-0000-0000-0000-000000000000"
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func MakeUserAdmin(name string) (err error) {
	u := getUserByUsername(name)
	err = u.changeRole("superadmin")
	return
}

func (u *User) save(password string) error {
	u.salt = createSalt()
	u.passwordHash = hashPassword(password, u.salt)
	db := connectSQLDB()
	defer db.Close()
	query := `insert into users values (?, ?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(query)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return err
	}
	_, err = stmt.Exec(u.Id, u.Username, u.EmailAddress, u.passwordHash, u.salt, u.Role)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to create user (MySQL): %v", err))
		return err
	}
	return nil
}

func (u *User) changeRole(role string) (err error) {
	db := connectSQLDB()
	defer db.Close()
	query := `update users set role=? where id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return
	}
	_, err = stmt.Exec(role, u.Id)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to change user role (MySQL): %v", err))
		return
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
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select id from users where id=?`)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return false
	}
	var id int
	err = selectUser.QueryRow(u.Id).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			_ = logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
			return false
		}
	}
	return true
}

func (u *User) isUsernameNew() bool {
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select username from users where username=?`)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return false
	}
	var username string
	err = selectUser.QueryRow(u.Username).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return true
		} else {
			_ = logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
			return false
		}
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
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select emailAddress from users where emailAddress=?`)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return false
	}
	var emailAddress string
	err = selectUser.QueryRow(u.EmailAddress).Scan(&emailAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			return true
		} else {
			_ = logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
			return false
		}
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

func userIsSame(u1 *User, u2 *User) bool {
	if u1.Id == u2.Id {
		return true
	}
	return false
}
