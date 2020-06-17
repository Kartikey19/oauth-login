package website

import (
	"net/http"
	config "recro_demo/config"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/facebook"
	"github.com/dghubble/gologin/v2/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
)

// Website interface for accessing all website related functionality
type Website struct {
	Env    *config.Env
	Router *chi.Mux
}

// GetRouter : Return a bisnode v1 api router
func (web *Website) GetRouter() *chi.Mux {
	web.Router = chi.NewRouter()
	web.registerRoutes()
	return web.Router
}

// GetRoutes returns a new ServeMux with app routes.
func (web *Website) registerRoutes() {
	web.Router.HandleFunc("/", web.welcomeHandler)
	web.Router.Handle("/profile", requireLogin(http.HandlerFunc(web.profileHandler)))
	web.Router.HandleFunc("/logout", web.logoutHandler)
	web.registerGithubRoutes()
	web.registerFacebookRoutes()
	web.registerTwitterRoutes()
}

func (web *Website) registerGithubRoutes() {
	// 1. Register LoginHandler and CallbackHandler
	// oauth2GithubConfig := &oauth2.Config{
	// 	ClientID:     web.Env.Config.GithubClientID,
	// 	ClientSecret: web.Env.Config.GithubClientSecret,
	// 	RedirectURL:  githubCallBackURL,
	// 	Endpoint:     githubOAuth2.Endpoint,
	// }
	// // state param cookies require HTTPS by default; disable for localhost development
	// // stateConfig := gologin.DebugOnlyCookieConfig
	web.Router.HandleFunc("/github/login", web.redirectHandler)
	web.Router.HandleFunc("/github/callback", web.callbackHandler)
}

func (web *Website) registerFacebookRoutes() {
	// 1. Register Login and Callback handlers
	oauth2FacebookConfig := &oauth2.Config{
		ClientID:     web.Env.Config.FacebookClientID,
		ClientSecret: web.Env.Config.FacebookClientSecret,
		RedirectURL:  faceBookCallBackURL,
		Endpoint:     facebookOAuth2.Endpoint,
		Scopes:       []string{"email"},
	}
	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	web.Router.Handle("/facebook/login", facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2FacebookConfig, nil)))
	web.Router.Handle("/facebook/callback", facebook.StateHandler(stateConfig, facebook.CallbackHandler(oauth2FacebookConfig, issueSession(), nil)))
}

func (web *Website) registerTwitterRoutes() {
	// 1. Register Twitter login and callback handlers
	oauth1TwitterConfig := &oauth1.Config{
		ConsumerKey:    web.Env.Config.TwitterConsumerKey,
		ConsumerSecret: web.Env.Config.TwitterConsumerSecret,
		CallbackURL:    twitterCallBackURL,
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}
	web.Router.Handle("/twitter/login", twitter.LoginHandler(oauth1TwitterConfig, nil))
	web.Router.Handle("/twitter/callback", twitter.CallbackHandler(oauth1TwitterConfig, issueSession(), nil))
}
