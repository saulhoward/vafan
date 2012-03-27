// Vafan - a web server for Convict Films
//
// Session
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	"code.google.com/p/gorilla/sessions"
	"errors"
	"fmt"
	"github.com/fzzbt/radix"
	"net/http"
	"net/url"
	"time"
)

var ErrUserNotLoggedIn = errors.New("session: user not logged in")
var ErrResourceRedirected = errors.New("session: resource was redirected")

// used for flashes
var sessionStore = sessions.NewCookieStore([]byte("something-very-secret"))

type session struct {
	id   string
	user *User
}

// Fetch user from cookie, set cookie, sync cookies x-domain
// main cookie flexer, called before every resource handler
func userCookie(w http.ResponseWriter, r *http.Request) (u *User, err error) {

	// login cookie
	c, err := r.Cookie("vafanLogin")
	if err != nil {
		if err == http.ErrNoCookie {
			err = nil
		} else {
			_ = logger.Err(fmt.Sprintf("Failed getting cookie: %v", err))
		}
	} else {
		_ = logger.Info(fmt.Sprintf("Login cookie found: %v", c.Value))
		u, err = getLoginUser(c.Value)
		if err == nil {
			return
		} else {
			_ = logger.Err(fmt.Sprintf("Failed getting login user: %v", err))
		}
	}

	// normal user cookie
	c, err = r.Cookie("vafanUser")
	if err != nil {
		if err == http.ErrNoCookie {
			err = nil
			// we have no user cookie
			s, env := getSite(r)
			canUserId := r.URL.Query().Get("canonical-user-id")
			userSyncSite := resourceCanonicalSites["usersSyncResource"]
			if s.Name != userSyncSite.Name && canUserId == "" {
				// we're on another site to the sync resource
				// redirect to the user sync!
				syncUrl := usersSyncResource{}.URL(r, nil)
				redirectUrl := syncUrl.String() + "?redirect-url=" + url.QueryEscape(getCurrentUrl(r).String())
				_ = logger.Info(fmt.Sprintf("Redirecting to sync url: %v", redirectUrl))
				http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
				err = ErrResourceRedirected
				return
			} else {
				// ok set a new cookie then
				if canUserId != "" {
					u = GetUser(canUserId)
				} else {
					u = NewUser()
				}
				_ = logger.Info(fmt.Sprintf("Setting a new user cookie: %v", u.Id))
				c = new(http.Cookie)
				c.Name = "vafanUser"
				c.Value = u.Id
				c.Domain = "." + env + "." + s.Host
				c.Path = "/"
				http.SetCookie(w, c)

				if canUserId != "" {
					// still got that query string? redirect again!
					curUrl := getCurrentUrl(r)
					q := curUrl.Query()
					curUrl.RawQuery = "" // remove the query string
					q.Del("canonical-user-id")
					var curUrlStr string
					if len(q) > 0 {
						curUrlStr = curUrl.String() + "?" + q.Encode() // and add it back
					} else {
						curUrlStr = curUrl.String()
					}
					http.Redirect(w, r, curUrlStr, http.StatusTemporaryRedirect)
				}
			}
		} else {
			_ = logger.Err(fmt.Sprintf("Failed getting cookie user: %v", err))
			u = NewUser()
		}
	} else {
		// we have a user cookie already
		u = GetUser(c.Value)
	}

	return
}

func newLoginSession(w http.ResponseWriter, r *http.Request, u *User) (s *session, err error) {
	sess := session{newUUID(), u}
	s = &sess
	err = nil
	sessionKey := "sessions:" + s.id
	userInfo := map[string]string{
		"Id":           u.Id,
		"Username":     u.Username,
		"EmailAddress": u.EmailAddress,
		"Role":         u.Role,
	}
	db := radix.NewClient(radix.Configuration{
		Database: 0,  // (default: 0)
		Timeout:  10, // (default: 10)
		Address:  "127.0.0.1:6379",
	})
	defer db.Close()
	reply := db.Command("hmset", sessionKey, userInfo)
	if reply.Error() != nil {
		errText := fmt.Sprintf("Failed to set Session data (Redis): %v", reply.Error())
		_ = logger.Err(errText)
		err = errors.New(errText)
		return
	}
	// set login cookie
	_, err = r.Cookie("vafanLogin")
	if err != nil {
		if err == http.ErrNoCookie {
			//no cookie, set one
			err = nil
			_ = logger.Info("Setting login cookie.")
			si, env := getSite(r)
			c := new(http.Cookie)
			c.Name = "vafanLogin"
			c.Value = s.id
			c.Path = "/"
			c.Domain = "." + env + "." + si.Host
			http.SetCookie(w, c)
		} else {
			_ = logger.Err(fmt.Sprintf("Failed getting login cookie (when trying to set): %v", err))
			return
		}
	} else {
		_ = logger.Notice("Login cookie already set!")
		err = nil
	}
	return
}

func logout(w http.ResponseWriter, r *http.Request, u *User) {
	// delete login cookie
	c, err := r.Cookie("vafanLogin")
	if err != nil {
		if err == http.ErrNoCookie {
			//no cookie, no problems
			err = nil
			return
		} else {
			_ = logger.Err(fmt.Sprintf("Failed getting login cookie (when trying to logout): %v", err))
			return
		}
	} else {
		_ = logger.Info("Attempting to delete login cookie.")
		si, env := getSite(r)
		c = new(http.Cookie)
		c.Name = "vafanLogin"
		c.Value = ""
		c.Path = "/"
		c.Domain = "." + env + "." + si.Host
		c.MaxAge = -1
		t := time.Time{}
		c.Expires = t
		http.SetCookie(w, c)
	}
	return
}

func getLoginUser(sId string) (u *User, err error) {
	err = nil
	u = NewUser()
	sessionKey := "sessions:" + sId
	// get user
	db := radix.NewClient(radix.Configuration{
		Database: 0,  // (default: 0)
		Timeout:  10, // (default: 10)
		Address:  "127.0.0.1:6379",
	})
	defer db.Close()
	reply := db.Command("hgetall", sessionKey)
	if reply.Error() != nil {
		errText := fmt.Sprintf("Failed to get Session data (Redis): %v", reply.Error())
		_ = logger.Err(errText)
		err = errors.New(errText)
		return
	}
	userInfo, err := reply.StringMap()
	if err != nil {
		errText := fmt.Sprintf("Stringmap failed (Redis): %v", reply.Error())
		_ = logger.Err(errText)
		err = errors.New(errText)
		return
	}
	u, err = getUserForUserInfo(userInfo)
	if err != nil {
		return
	}
	u.setLoggedIn()
	_ = logger.Info(fmt.Sprintf("User is logged in: %v", u.Id))
	return
}

func addFlash(w http.ResponseWriter, r *http.Request, msg string, level string) {
	_ = logger.Info(fmt.Sprintf("Flashing: %v - %v", level, msg))
	flash, _ := sessionStore.Get(r, "vafanFlashes")
	flash.AddFlash(msg, level)
	flash.Save(r, w)
}

func getFlashContent(w http.ResponseWriter, r *http.Request) (content map[string]interface{}) {
	content = make(map[string]interface{})
	flash, err := sessionStore.Get(r, "vafanFlashes")
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed getting flashes: %v", err))
		return
	}
	if f := flash.Flashes("error"); len(f) > 0 {
		content["error"] = f
	}
	if f := flash.Flashes("success"); len(f) > 0 {
		content["success"] = f
	}
	if f := flash.Flashes("warning"); len(f) > 0 {
		content["warning"] = f
	}
	if f := flash.Flashes("information"); len(f) > 0 {
		content["information"] = f
	}
	err = flash.Save(r, w)
	if err != nil {
		_ = logger.Err(fmt.Sprintf("Failed saving flashes: %v", err))
		return
	}
	return
}
