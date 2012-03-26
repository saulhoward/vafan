// Vafan - a web server for Convict Films
//
// Resource handlers
//
// @url    http://saulhoward.com/vafan
// @author saul@saulhoward.com
//
package vafan

import (
	//"fmt"
	//"log"
    "reflect"
    "regexp"
	"net/url"
	"net/http"
	"code.google.com/p/gorilla/mux"
	"code.google.com/p/gorilla/schema"
)

// decodes form values
var decoder = schema.NewDecoder()

// Our resources must return URL & Content, and they must
// serveHTTP
type Resource interface {
	URL(req *http.Request, s *site) *url.URL
	Content(req *http.Request, s *site) resourceContent
    ServeHTTP(w http.ResponseWriter, r *http.Request, u *User)
}

// generic data map for resource content
type resourceData map[string]interface{}

type resourceContent struct {
    title       string
    description string
    content     resourceData
}

// somewhat crufty helper map
var emptyContent = map[string]interface{}{}

// crufty...
var resourceCanonicalSites = map[string]*site{
    "usersRegistrarResource": defaultSite,
    "usersAuthResource":      defaultSite,
    "usersSyncResource":      defaultSite,
}

// Gets a URL for a resource.
// used as a helper by someResource.URL(r, s) function
func getUrl(res Resource, req *http.Request, s *site, urlData []string) *url.URL {
    curSite, env := getSite(req)
    canonicalSite, err := getCanonicalSite(res)
    if s == nil {
        s = curSite
    }
    if err == nil && canonicalSite.Name != s.Name {
        s = canonicalSite
    }
	format := getFormat(req)
    format = "." + format
    if format == ".html" {
        format = ""
    }
    host := env + "." + s.Host + ":8888"
    urlPairs := []string{"format", format, "host", host}
    if urlData != nil {
        for _, p := range urlData {
            urlPairs = append(urlPairs, p)
        }
    }
    url, err := router.GetRoute(resourceName(res)).Host(hostRe).URL(urlPairs...)
    checkError(err)
    return url
}

// get the resource's type name by reflection
func resourceName(r Resource) string {
    n := reflect.TypeOf(r).String()
    re := regexp.MustCompile(`\.([a-zA-Z]+)$`)
    m := re.FindStringSubmatch(n)
    return m[1]
}

// --

// Now the actual resources, these will go in separate files, I think

// -- Index resource

type indexResource struct {
}

func (res indexResource) URL(req *http.Request, s *site) *url.URL {
    return getUrl(res, req, s, nil)
}

func (res indexResource) Content(req *http.Request, s *site) (c resourceContent) {
    c.title = s.Tagline
    c.description = "Home page"
    c.content = emptyContent
    return
}

func (res indexResource) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *User) {
    writeResource(w, r, res, reqU)
    return
}

// -- Registrar resource

type usersRegistrarResource struct {
    data resourceData
}

func (res usersRegistrarResource) URL(req *http.Request, s *site) *url.URL {
    // limit registration to default site
    return getUrl(res, req, defaultSite, nil)
}

func (res usersRegistrarResource) Content(req *http.Request, s *site) (c resourceContent) {
    c.title = "Register"
    c.description = "Register here to access Convict Films"
    if res.data == nil {
        res.data = emptyContent
    }
    c.content = res.data
    return
}

func (res usersRegistrarResource) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *User) {
	switch r.Method {
	case "POST":
		// This is a post to create a new user
		r.ParseForm()
        u := new(User)
		decoder.Decode(u, r.Form)

		// check for errors in post
        if !u.isLegal(r.Form.Get("Password")) || r.Form.Get("Password") != r.Form.Get("RepeatPassword") || !u.isNew() {
            // found errors in post
			errors := map[string]interface{}{}
			if !u.isUsernameLegal() {
				errors["Username"] = "Must contain only letters and numbers, with no spaces."
            } else if !u.isUsernameNew() {
				errors["Username"] = "Username already taken, sorry."
            }
			if !u.isEmailAddressLegal() {
				errors["EmailAddress"] = "Must be a valid email address."
            } else if !u.isEmailAddressNew() {
				errors["EmailAddress"] = "This email address is already associated with another user."
            }
			if !u.isPasswordLegal(r.Form.Get("Password")) {
				errors["Password"] = "Password must be more than 6 characters."
            } else if r.Form.Get("Password") != r.Form.Get("RepeatPassword") {
				errors["Password"] = "Password must match repeat password."
			}
			res.data["errors"] = errors
			writeResource(w, r, res, u)
			return
		}

        // legal user, try to save
        err := u.save(r.Form.Get("Password"))
        var url *url.URL
        if err != nil {
            url = usersRegistrarResource{}.URL(r, nil)
            addFlash(w, r, "Failed to save new user", "error")
        } else {
            url = usersAuthResource{}.URL(r, nil)
            addFlash(w, r, "Registered a new user, please log in.", "success")
        }

        http.Redirect(w, r, url.String(), http.StatusSeeOther)
		return
	case "GET":
        if reqU.isNew() {
            writeResource(w, r, res, reqU)
        } else {
            url := usersAuthResource{}.URL(r, nil)
            addFlash(w, r, "Your user ID already has an account, please log in.", "warning")
            http.Redirect(w, r, url.String(), http.StatusSeeOther)
        }
		return
	}
}

// -- Auth resource

type usersAuthResource struct {
    data resourceData
}

func (res usersAuthResource) URL(req *http.Request, s *site) *url.URL {
    // limit authentication to default site
    return getUrl(res, req, defaultSite, nil)
}

func (res usersAuthResource) Content(req *http.Request, s *site) (c resourceContent) {
    c.title = "Login"
    c.description = "Login here to access Convict Films"
    if res.data == nil {
        res.data = emptyContent
    }
    c.content = res.data
    return
}

func (res usersAuthResource) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *User) {
	switch r.Method {
	case "POST":
		// This is a post to login or logout
        var url *url.URL
		r.ParseForm()
        switch {
        case r.Form.Get("login") != "":
            // try to login
            // TODO: THE CURRENT USER MUST BE LOGGED OUT

            // login user
            loginUser, err := login(r.Form.Get("UsernameOrEmailAddress"), r.Form.Get("Password"))
            if err != nil {
                url = usersAuthResource{}.URL(r, nil)
                addFlash(w, r, "Failed to login", "error")
            } else {
                // set the login session
                _, err := newLoginSession(w, r, loginUser)
                if err != nil {
                    checkError(err)
                }
                url = indexResource{}.URL(r, nil)
                addFlash(w, r, "Login!", "success")
            }
            http.Redirect(w, r, url.String(), http.StatusSeeOther)
        case r.Form.Get("logout") != "":
            // try to logout
            logout(w, r, reqU)
            addFlash(w, r, "Logged out.", "success")
            url = indexResource{}.URL(r, nil)
            http.Redirect(w, r, url.String(), http.StatusSeeOther)
        }
	case "GET":
        writeResource(w, r, res, reqU)
    }
    return
}

// -- Sync resource

type usersSyncResource struct {
}

func (res usersSyncResource) URL(req *http.Request, s *site) *url.URL {
    // limit sync to default site
    return getUrl(res, req, defaultSite, nil)
}

func (res usersSyncResource) Content(req *http.Request, s *site) (c resourceContent) {
    c.title = "User Sync"
    c.description = "Performs a user sync redirect"
    c.content = emptyContent
    return
}

// send people back to the redirect-url param, with a canonical user id
func (res usersSyncResource) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *User) {
    ruStr := r.URL.Query().Get("redirect-url")
    if ruStr == "" {
        ruStr = "/"
    }
    ru, err := url.Parse(ruStr)
    checkError(err)
    q := ru.Query()
    ru.RawQuery = "" // remove the query string
    q.Set("canonical-user-id", url.QueryEscape(reqU.Id))
    fullUrl := ru.String() + "?" + q.Encode() // and add it back
    print("\nReturning to url... " + fullUrl)
    http.Redirect(w, r, fullUrl, http.StatusTemporaryRedirect)
    return
}

// -- User resource

type usersResource struct {
    user *User
}

func (res usersResource) URL(req *http.Request, s *site) *url.URL {
    return getUrl(res, req, s, []string{"id", res.user.Id})
}

func (res usersResource) Content(req *http.Request, s *site) (c resourceContent) {
    c.title = "User"
    c.description = "User page"
    c.content = map[string]interface{}{"user": res.user}
    return
}

func (res usersResource) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *User) {
    // check if user has permission? whose user page is this?
	vars := mux.Vars(r)
    res.user = getUserById(vars["id"])
    if userIsSame(reqU, res.user) || reqU.Role == "admin" {
        writeResource(w, r, res, reqU)
        return
    }
    forbiddenResource{}.ServeHTTP(w, r, reqU)
    return
}

// -- 404 resource

type notFoundResource struct {
}

func (res notFoundResource) URL(req *http.Request, s *site) *url.URL {
    return getUrl(res, req, s, nil)
}

func (res notFoundResource) Content(req *http.Request, s *site) (c resourceContent) {
    c.title = "404 - Not Found"
    c.description = "The requested resource does not exist."
    content := map[string]interface{}{}
    content["message"] = "404 Not found"
    content["body"] = "Sorry, this resource could not be found."
    c.content = content
    return
}

func (res notFoundResource) ServeHTTP(w http.ResponseWriter, r *http.Request, u *User) {
    w.WriteHeader(http.StatusNotFound)
    writeResource(w, r, res, u)
    return
}

// -- Forbidden resource (403)

type forbiddenResource struct {
}

func (res forbiddenResource) URL(req *http.Request, s *site) *url.URL {
    return getUrl(res, req, s, nil)
}

func (res forbiddenResource) Content(req *http.Request, s *site) (c resourceContent) {
    c.title = "403 - Forbidden"
    c.description = "You are forbidden to access this resource"
    content := map[string]interface{}{}
    content["message"] = "403 Forbidden"
    content["body"] = "Sorry, you may not access this resource."
    c.content = content
    return
}

func (res forbiddenResource) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *User) {
    w.WriteHeader(http.StatusForbidden)
    writeResource(w, r, res, reqU)
    return
}

// --  Video resource -- a particular video

type videoResource struct {
    video *video
}

func (res videoResource) URL(req *http.Request, s *site) *url.URL {
    return getUrl(res, req, s, []string{"name", res.video.Name})
}

func (res videoResource) Content(req *http.Request, s *site) (c resourceContent) {
    c.title = "Video"
    c.description = "Video page"
    c.content = map[string]interface{}{"video": res.video}
    return
}

func (res videoResource) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *User) {
    switch r.Method {
    case "GET":
        vars := mux.Vars(r)
        var err error
        res.video, err = getVideoByName(vars["name"])
        if err != nil {
            if err == ErrVideoNotFound {
                notFoundResource{}.ServeHTTP(w, r, reqU)
                return
            }
            checkError(err)
        }
        writeResource(w, r, res, reqU)
        return
    }
}

// -- Videos resource -- all videos

type videosResource struct {
    video *video
}

func (res videosResource) URL(req *http.Request, s *site) *url.URL {
    return getUrl(res, req, s, nil)
}

func (res videosResource) Content(req *http.Request, s *site) (c resourceContent) {
    c.title = "Video"
    c.description = "Video page"
    c.content = emptyContent
    return
}

func (res videosResource) ServeHTTP(w http.ResponseWriter, r *http.Request, reqU *User) {
    switch r.Method {
    case "GET":
        writeResource(w, r, res, reqU)
        return
    case "POST":

    }
    return
}

