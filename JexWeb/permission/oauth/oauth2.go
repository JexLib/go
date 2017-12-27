package oauth

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/JexLib/golang/JexWeb/session"
	"github.com/labstack/echo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/linkedin"
)

const (
	KEY_TOKEN     = "jex_oauth2_token"
	KEY_NEXT_PAGE = "next"
)

var (
	// PathLogin is the path to handle OAuth 2.0 logins.
	PathLogin = "/oauth/login"
	// PathLogout is the path to handle OAuth 2.0 logouts.
	PathLogout = "/oauth/logout"
	// PathCallback is the path to handle callback from OAuth 2.0 backend
	// to exchange credentials.
	PathCallback = "/oauth/callback"
	// PathError is the path to handle error cases.
	PathError = "/oauth/error"
)

// Tokens represents a container that contains user's OAuth 2.0 access and refresh tokens.
type Tokens interface {
	Access() string
	Refresh() string
	Expired() bool
	ExpiryTime() time.Time
}

type token struct {
	oauth2.Token
}

// Access returns the access token.
func (t *token) Access() string {
	return t.AccessToken
}

// Refresh returns the refresh token.
func (t *token) Refresh() string {
	return t.RefreshToken
}

// Expired returns whether the access token is expired or not.
func (t *token) Expired() bool {
	if t == nil {
		return true
	}
	return !t.Token.Valid()
}

// ExpiryTime returns the expiry time of the user's access token.
func (t *token) ExpiryTime() time.Time {
	return t.Expiry
}

// String returns the string representation of the token.
func (t *token) String() string {
	return fmt.Sprintf("tokens: %v", t)
}

// Google returns a new Google OAuth 2.0 backend endpoint.
func Google(conf *oauth2.Config) echo.MiddlewareFunc {
	conf.Endpoint = google.Endpoint
	return NewOAuth2Provider(conf)
}

// Github returns a new Github OAuth 2.0 backend endpoint.
func Github(conf *oauth2.Config) echo.MiddlewareFunc {
	conf.Endpoint = github.Endpoint
	return NewOAuth2Provider(conf)
}

// Facebook returns a new Facebook OAuth 2.0 backend endpoint.
func Facebook(conf *oauth2.Config) echo.MiddlewareFunc {
	conf.Endpoint = facebook.Endpoint
	return NewOAuth2Provider(conf)
}

// LinkedIn returns a new LinkedIn OAuth 2.0 backend endpoint.
func LinkedIn(conf *oauth2.Config) echo.MiddlewareFunc {
	conf.Endpoint = linkedin.Endpoint
	return NewOAuth2Provider(conf)
}

// NewOAuth2Provider returns a generic OAuth 2.0 backend endpoint.
func NewOAuth2Provider(conf *oauth2.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == "GET" {
				switch c.Request().URL.Path {
				case PathLogin:
					return login(conf, c)
				case PathLogout:
					return logout(c)
				case PathCallback:
					return handleOAuth2Callback(conf, c)
				}
			}
			tk := unmarshallToken(c)
			sess := session.Default(c)
			if tk != nil {
				// check if the access token is expired
				if tk.Expired() && tk.Refresh() == "" {
					sess.Delete(KEY_TOKEN)
					sess.Save()
					tk = nil
				}
			}
			//sess.Set()
			// Inject tokens.
			//ctx.MapTo(tk, (*Tokens)(nil))
			return next(c)
		}
	}
}

func SkipperForOAuth2(c echo.Context) bool {
	return false
}

// Handler that redirects user to the login page
// if user is not logged in.
// Sample usage:
// m.Get("/login-required", oauth2.LoginRequired, func() ... {})
var LoginRequired = func() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := unmarshallToken(c)
		if token == nil || token.Expired() {
			next := url.QueryEscape(c.Request().URL.RequestURI())
			return c.Redirect(302, PathLogin+"?next="+next)
		}
		return nil
	}
}()

func login(f *oauth2.Config, c echo.Context) error {
	next := extractPath(c.QueryParam(KEY_NEXT_PAGE))
	sess := session.Default(c)
	if sess.Get(KEY_TOKEN) == nil {
		// User is not logged in.
		if next == "" {
			next = "/"
		}
		return c.Redirect(302, f.AuthCodeURL(next))

	}
	// No need to login, redirect to the next page.
	return c.Redirect(302, next)
}

func logout(c echo.Context) error {
	next := extractPath(c.QueryParam(KEY_NEXT_PAGE))
	sess := session.Default(c)
	sess.Delete(KEY_TOKEN)
	sess.Save()
	return c.Redirect(302, next)
}

func handleOAuth2Callback(f *oauth2.Config, c echo.Context) error {
	next := extractPath(c.QueryParam("state"))
	code := c.QueryParam("code")
	t, err := f.Exchange(oauth2.NoContext, code)
	if err != nil {
		// Pass the error message, or allow dev to provide its own
		// error handler.

		return c.Redirect(302, PathError)
	}
	// Store the credentials in the session.
	val, _ := json.Marshal(t)
	sess := session.Default(c)
	sess.Set(KEY_TOKEN, val)
	sess.Save()
	return c.Redirect(302, next)
}

func unmarshallToken(c echo.Context) (t *token) {
	sess := session.Default(c)
	if sess.Get(KEY_TOKEN) == nil {
		return
	}
	data := sess.Get(KEY_TOKEN).([]byte)
	var tk oauth2.Token
	json.Unmarshal(data, &tk)
	return &token{tk}
}

func extractPath(next string) string {
	n, err := url.Parse(next)
	if err != nil {
		return "/"
	}
	return n.Path
}
