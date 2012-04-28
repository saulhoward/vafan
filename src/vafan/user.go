// Vafan - a web server for Convict Films
//
// User
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	"code.google.com/p/gorilla/mux"
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
	"net/url"
	"regexp"
	"strings"
)

var ErrWrongPassword = errors.New("user: password fail")

type user struct {
	ID           string `json:"id"` // UUID v4 with dashes
	Username     string `json:"username"`
	EmailAddress string `json:"emailAddress"`
	Role         string `json:"role"`
	URL          string `json:"url"`
	IsLoggedIn   bool   `json:"isLoggedIn"`
	passwordHash string
	salt         string
}

// HTTP Resource methods

func (u user) GetURL(req *http.Request, s *site) *url.URL {
	return makeURL(u, req, s, []string{"id", u.ID})
}

func (u user) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *user) {
	res := Resource{
		title:       "User",
		description: "User page",
	}
	res.content = make(resourceContent)

	// check if user has permission? whose user page is this?
	vars := mux.Vars(r)
	u = *getUserById(vars["id"])
	if userIsSame(reqU, &u) || reqU.Role == "superadmin" {
		res.content["user"] = u
		res.write(w, r, &u, reqU)
		return
	}
	forbidden{}.ServeHTTP(w, r, reqU)
	return
}

// Other methods

const defaultUserRole = "user"

// brand new user, freshly minted id
func NewUser() *user {
	u := user{ID: newUUID(), Role: defaultUserRole}
	return &u
}

// -- DB

func connectSQLDB() *sql.DB {
	db, err := sql.Open("mymysql", fmt.Sprintf("vafan/%v/%v", vafanConf.mysql.user, vafanConf.mysql.password))
	if err != nil {
		panic("Error connecting to mysql db: " + err.Error())
	}
	return db
}

func createSalt() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to create random salt: %v", err))
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

func getUserByUsername(username string) (u *user) {
	u = NewUser()
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select id, username, emailAddress, role, passwordHash, salt from users where username=?`)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return
	}
	err = selectUser.QueryRow(username).Scan(&u.ID, &u.Username, &u.EmailAddress, &u.Role, &u.passwordHash, &u.salt)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
		return
	}
	return
}

func getUserByEmailAddress(emailAddress string) (u *user) {
	u = NewUser()
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select id, username, emailAddress, role, passwordHash, salt from users where emailAddress=?`)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return
	}
	err = selectUser.QueryRow(emailAddress).Scan(&u.ID, &u.Username, &u.EmailAddress, &u.Role, &u.passwordHash, &u.salt)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
		return
	}
	return
}

func getUserById(id string) (u *user) {
	u = NewUser()
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select id, username, emailAddress, role, passwordHash, salt from users where id=?`)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return
	}
	err = selectUser.QueryRow(id).Scan(&u.ID, &u.Username, &u.EmailAddress, &u.Role, &u.passwordHash, &u.salt)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
		return
	}
	return
}

// just wants an id, the simplest form of user
func GetUser(id string) *user {
	u := user{ID: id, Role: defaultUserRole}
	return &u
}

// needs a map of user properties, id must be set
func getUserForUserInfo(userInfo map[string]string) (u *user, err error) {
	if userInfo["Id"] == "" {
		err = errors.New("User: ID must be set")
		return
	}
    newU := user{
        ID: userInfo["Id"],
        Username: userInfo["Username"],
        EmailAddress: userInfo["EmailAddress"],
        Role: userInfo["Role"],
    }
    return &newU, err
}

// Use UUID v4 as user IDs
func newUUID() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to create random uuid: %v", err))
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

func (u *user) save(password string) error {
	u.salt = createSalt()
	u.passwordHash = hashPassword(password, u.salt)
	db := connectSQLDB()
	defer db.Close()
	query := `insert into users values (?, ?, ?, ?, ?, ?)`
	stmt, err := db.Prepare(query)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return err
	}
	_, err = stmt.Exec(u.ID, u.Username, u.EmailAddress, u.passwordHash, u.salt, u.Role)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to create user (MySQL): %v", err))
		return err
	}
	return nil
}

func (u *user) changeRole(role string) (err error) {
	db := connectSQLDB()
	defer db.Close()
	query := `update users set role=? where id=?`
	stmt, err := db.Prepare(query)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return
	}
	_, err = stmt.Exec(role, u.ID)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to change user role (MySQL): %v", err))
		return
	}
	return nil
}

func (u *user) isLegal(password string) bool {
	if !u.isUsernameLegal() ||
		!u.isEmailAddressLegal() ||
		!u.isPasswordLegal(password) {
		return false
	}
	return true
}

func (u *user) isUsernameLegal() bool {
	var re = regexp.MustCompile(`[^\d+|\w+]`)
	if re.MatchString(u.Username) {
		return false
	}
	return true
}

func (u *user) isRegistered() bool {
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select id from users where id=?`)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return false
	}
	var id int
	err = selectUser.QueryRow(u.ID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
			return false
		}
	}
	return true
}

// TODO: Saul saul and SaUl should be the same here.
func (u *user) isUsernameNew() bool {
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select username from users where username=?`)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return false
	}
	var username string
	err = selectUser.QueryRow(u.Username).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return true
		} else {
			logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
			return false
		}
	}
	return false
}

func (u *user) isEmailAddressLegal() bool {
	if strings.Contains(u.EmailAddress, "@") {
		return true
	}
	return false
}

func (u *user) isEmailAddressNew() bool {
	db := connectSQLDB()
	defer db.Close()
	selectUser, err := db.Prepare(`select emailAddress from users where emailAddress=?`)
	if err != nil {
		logger.Err(fmt.Sprintf("Failed to prepare db (MySQL): %v", err))
		return false
	}
	var emailAddress string
	err = selectUser.QueryRow(u.EmailAddress).Scan(&emailAddress)
	if err != nil {
		if err == sql.ErrNoRows {
			return true
		} else {
			logger.Err(fmt.Sprintf("Failed to select user (MySQL): %v", err))
			return false
		}
	}
	return false
}

func (u *user) isPasswordLegal(password string) bool {
	if len(password) < 6 {
		return false
	}
	return true
}

func (u *user) isNew() bool {
	if u.isRegistered() || !u.isUsernameNew() || !u.isEmailAddressNew() {
		return false
	}
	return true
}

func (u *user) setLoggedIn() {
	u.IsLoggedIn = true
}

// Login
func login(usernameOrEmailAddress string, password string) (u *user, err error) {
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

func userIsSame(u1 *user, u2 *user) bool {
	if u1.ID == u2.ID {
		return true
	}
	return false
}
