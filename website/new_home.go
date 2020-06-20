package website

import (
	"fmt"
	"net/http"

	gocialite "gopkg.in/danilopolani/gocialite.v1"
)

// Was trying to build my own Handler for callback
// Used gologin instead
var gocial = gocialite.NewDispatcher()

// Redirect to correct oAuth URL
func (web *Website) redirectHandler(w http.ResponseWriter, r *http.Request) {
	authURL, err := gocial.New().
		Driver("github").                // Set provider
		Scopes([]string{"public_repo"}). // Set optional scope(s)
		Redirect(                        //
			web.Env.Config.GithubClientID,           // Client ID
			web.Env.Config.FacebookClientSecret,     // Client Secret
			"http://localhost:8080/github/callback", // Redirect URL
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	// Redirect with authURL
	http.Redirect(w, r, authURL, 302) // Redirect with 302 HTTP code
}

func (web *Website) callbackHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve query params for code and state
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// Handle callback and check for errors
	user, token, err := gocial.Handle(state, code)
	if err != nil {
		w.Write([]byte("Error: " + err.Error()))
		return
	}

	// Print in terminal user information
	fmt.Printf("%#v", token)
	fmt.Printf("%#v", user)

	// If no errors, show provider name
	w.Write([]byte("Hi, " + user.FullName))
}
