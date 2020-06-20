package website

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"recro_demo/jsonwrap"
	"recro_demo/postgres"

	"github.com/dghubble/gologin/facebook"
	"github.com/dghubble/gologin/github"
	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/sessions"
)

const (
	sessionName    = "recro-demo-app"
	sessionSecret  = "a123b53609191910"
	sessionUserKey = "userID"
	sourceGithub   = "github"
	sourceFaceBook = "facebook"
	sourceTwitter  = "twitter"
)

// sessionStore encodes and decodes session data stored in signed cookies
var sessionStore = sessions.NewCookieStore([]byte(sessionSecret), nil)

// issueSession issues a cookie session after successful Github login
func (web *Website) issueSession(source string) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var ID string
		switch source {
		case sourceGithub:
			user, err := github.UserFromContext(ctx)

			ID = string(*user.ID)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// userMap := structs.Map(user)
			// web.saveData(source, userMap, w)
		case sourceFaceBook:
			user, err := facebook.UserFromContext(ctx)
			ID = user.ID
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// userMap := structs.Map(user)
			// web.saveData(source, userMap, w)
		case sourceTwitter:
			user, err := twitter.UserFromContext(ctx)
			ID = string(user.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// userMap := structs.Map(user)
			// web.saveData(source, userMap, w)
		}
		// 2. Implement a success handler to issue some form of session

		session := sessionStore.New(sessionName)
		session.Values[sessionUserKey] = ID
		session.Save(w)
		http.Redirect(w, req, "/profile", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

// welcomeHandler shows a welcome message and login button.
func (web *Website) welcomeHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	if isAuthenticated(req) {
		http.Redirect(w, req, "/profile", http.StatusFound)
		return
	}
	page, _ := ioutil.ReadFile("templates/home.html")
	fmt.Fprintf(w, string(page))
}

// profileHandler shows protected user content.
func (web *Website) profileHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, `<p>You are logged in!</p><form action="/logout" method="post"><input type="submit" value="Logout"></form>`)
}

// logoutHandler destroys the session on POSTs and redirects to home.
func (web *Website) logoutHandler(w http.ResponseWriter, req *http.Request) {
	// if req.Method == "POST" {

	// }
	sessionStore.Destroy(w, sessionName)
	http.Redirect(w, req, "/", http.StatusFound)
}

// requireLogin redirects unauthenticated users to the login route.
func requireLogin(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		if !isAuthenticated(req) {
			http.Redirect(w, req, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

// isAuthenticated returns true if the user has a signed session cookie.
func isAuthenticated(req *http.Request) bool {
	if _, err := sessionStore.Get(req, sessionName); err == nil {
		return true
	}
	return false
}

func (web *Website) saveData(
	source string, userData map[string]interface{}, w http.ResponseWriter) {
	jsonResp := make(map[string]interface{})
	// check if email in userData
	if _, ok := userData["email"]; !ok {
		resp, _ := jsonwrap.MakeJSONResponse(fmt.Sprintf("Error in fetching email from %s.", source), jsonResp, false)
		jsonwrap.SendJSONHttpResponse(w, resp, http.StatusInternalServerError)
		return
	}
	// check if phone exists
	if _, ok := userData["phone"]; !ok {
		resp, _ := jsonwrap.MakeJSONResponse(fmt.Sprintf("Error in fetching phone from %s.", source), jsonResp, false)
		jsonwrap.SendJSONHttpResponse(w, resp, http.StatusInternalServerError)
		return
	}

	// check if name exists
	if _, ok := userData["phone"]; !ok {
		resp, _ := jsonwrap.MakeJSONResponse(fmt.Sprintf("Error in fetching name from %s.", source), jsonResp, false)
		jsonwrap.SendJSONHttpResponse(w, resp, http.StatusInternalServerError)
		return
	}
	// check if email already exists
	dbUser := web.Env.DB.CheckUserExists(fmt.Sprintf("%v", userData["email"]))
	if dbUser.ID == 0 {
		// create new user
		web.createNewUser(jsonResp, source, w)
	} else {
		web.updateUser(jsonResp, source, dbUser, w)
	}

}

func (web *Website) createUserStruct(map[string]interface{}) *postgres.User {
	var dbUser *postgres.User
	return dbUser
}

func (web *Website) createNewUser(
	user map[string]interface{},
	source string, w http.ResponseWriter) {
	// creates new user
	dbUser := web.createUserStruct(user)
	metaMap := make(map[string]interface{})
	metaMap[source] = user

	dbUser.Meta = metaMap
	id := web.Env.DB.CreateUser(dbUser)
	if id == -1 {
		fmt.Println("New User: Valid user not created")
	}
}

func (web *Website) updateUser(
	user map[string]interface{},
	source string,
	dbUser *postgres.User, w http.ResponseWriter) {
	// update the current user
	jsonData := dbUser.Meta
	if jsonData == nil {
		jsonData = make(map[string]interface{})
	}
	if _, ok := jsonData[source]; !ok {
		jsonData[source] = user
		err := web.Env.DB.UpdateUserMeta(dbUser.ID, jsonData)
		if err != nil {
			fmt.Println("Repeated User:Meta not updated")
			return
		}
	}
}
