package website

import (
	"net/http"
	config "recro_demo/config"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	v2 "github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/facebook"
	"github.com/dghubble/gologin/v2/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
	githubOAuth2 "golang.org/x/oauth2/github"
)

// Website interface for accessing all website related functionality
type Website struct {
	Env          *config.Env
	Router       *chi.Mux
	GithubConfig *oauth2.Config
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
	web.Router.Route("/user", func(r chi.Router) {
		r.Use(requireLogin) // middleware layer
		r.Get("/all", web.fetchAllUsers)
		r.Get("/{id:[1-9]+}", web.fetchUser)
	})

	web.Router.HandleFunc("/logout", web.logoutHandler)
	web.registerGithubRoutes()
	web.registerFacebookRoutes()
	web.registerTwitterRoutes()
}

func (web *Website) registerGithubRoutes() {
	// 1. Register LoginHandler and CallbackHandler
	oauth2GithubConfig := &oauth2.Config{
		ClientID:     web.Env.Config.GithubClientID,
		ClientSecret: web.Env.Config.GithubClientSecret,
		RedirectURL:  githubCallBackURL,
		Endpoint:     githubOAuth2.Endpoint,
	}
	// // state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	web.Router.Handle("/github/login", github.StateHandler(stateConfig, github.LoginHandler(oauth2GithubConfig, nil)))
	web.Router.Handle("/github/callback", github.StateHandler(stateConfig, github.CallbackHandler(oauth2GithubConfig, web.issueSession(sourceGithub), nil)))
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
	stateConfig := v2.DebugOnlyCookieConfig
	web.Router.Handle("/facebook/login", facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2FacebookConfig, nil)))
	web.Router.Handle("/facebook/callback", facebook.StateHandler(stateConfig, facebook.CallbackHandler(oauth2FacebookConfig, web.issueSession(sourceFaceBook), nil)))
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
	web.Router.Handle("/twitter/callback", twitter.CallbackHandler(oauth1TwitterConfig, web.issueSession(sourceTwitter), nil))
}
