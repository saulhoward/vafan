// Vafan - a web server for Convict Films
//
// Session
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
    //"fmt"
    "errors"
    "net/http"
	"net/url"
    "github.com/fzzbt/radix"
    "code.google.com/p/gorilla/sessions"
)

var ErrUserNotLoggedIn = errors.New("session: user not logged in")

// used for flashes
var sessionStore = sessions.NewCookieStore([]byte("something-very-secret"))

type session struct {
    id   string
    user *User
}

// Fetch user from cookie, set cookie, sync cookies x-domain
// main cookie flexer, called before every resource handler
func userCookie(w http.ResponseWriter, r *http.Request) (u *User) {

    // login cookie
    c, err := r.Cookie("vafanLogin")
    if err != nil {
        if err == http.ErrNoCookie {
            err = nil
        } else {
            checkError(err)
        }
    } else {
        print("\nLogin cookie found.")
        u, err = getLoginUser(c.Value)
        checkError(err)
        return
    }

    // normal user cookie
    c, err = r.Cookie("vafanUser")
    if err != nil {
        if err == http.ErrNoCookie {
            err = nil
            // we have no user cookie
            s, env := getSite(r)
            canUserId := r.URL.Query().Get("canonical-user-id")
            userSyncSite := resourceCanonicalSites["usersSync"]
            if s.Name != userSyncSite.Name && canUserId == "" {
                // we're on another site to the sync resource
                // redirect to the user sync!
                sync := resources["usersSync"]
                syncUrl := getUrl(sync, r)
                redirectUrl := syncUrl.String() + "?redirect-url=" + url.QueryEscape(getCurrentUrl(r).String())
                print("\nRedirecting to sync url... " + redirectUrl)
                http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
                return
            } else {
                print("\nSetting a new cookie... ")
                // ok set a new cookie then
                if canUserId != "" {
                    u = GetUser(canUserId)
                } else {
                    u = NewUser()
                }
                c = new(http.Cookie)
                c.Name = "vafanUser"
                c.Value = u.Id
                c.Domain = "." + env + "." +  s.Host
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
                    http.Redirect(w, r, curUrlStr , http.StatusTemporaryRedirect)
                }
            }
        } else {
            checkError(err)
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
        Database: 0, // (default: 0)
        Timeout: 10, // (default: 10)
        Address: "127.0.0.1:6379",
    })
    defer db.Close()
    reply := db.Command("hmset", sessionKey, userInfo)
    if reply.Error() != nil {
        err = errors.New("set failed")
        //reply.Error()
        return
    }
    // set login cookie
    _, err = r.Cookie("vafanLogin")
    if err != nil {
        if err == http.ErrNoCookie {
            //no cookie, set one
            err = nil
            print("\nAttempting to set login cookie...")
            si, env := getSite(r)
            c := new(http.Cookie)
            c.Name = "vafanLogin"
            c.Value = s.id
            c.Path = "/"
            c.Domain = "." + env + "." +  si.Host
            http.SetCookie(w, c)
        } else {
            checkError(err)
        }
    } else {
        print("\nLogin cookie already set!")
    }
    return
}

func getLoginUser(sId string) (u *User, err error) {
    err = nil
    u = NewUser()
    sessionKey := "sessions:" + sId
    // get user
    db := radix.NewClient(radix.Configuration{
        Database: 0, // (default: 0)
        Timeout: 10, // (default: 10)
        Address: "127.0.0.1:6379",
    })
    defer db.Close()
    reply := db.Command("hgetall", sessionKey)
    if reply.Error() != nil {
        err = errors.New("get failed")
        return
    }
    userInfo, err := reply.StringMap()
    if err != nil {
        err = errors.New("stringmap failed")
        return
    }
    u, err = getUserForUserInfo(userInfo)
    if err != nil {
        return
    }
    u.setLoggedIn()
    print("\nUser is logged in as user " + u.Id)
    return
}

func addFlash(w http.ResponseWriter, r *http.Request, msg string, level string) {
    print("\nFlashing - " + level + ": " + msg)
    flash, _ := sessionStore.Get(r, "vafanFlashes")
    flash.AddFlash(msg, level)
    flash.Save(r, w)
}

func getFlashContent(w http.ResponseWriter, r *http.Request) map[string]interface{} {
    flash, err := sessionStore.Get(r, "vafanFlashes")
    checkError(err)
    content := make(map[string]interface{})
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
    checkError(err)
    return content
}
