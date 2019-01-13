package session

import (
	// base core imports
	"net/http"
	// external imports
	"github.com/gorilla/sessions"
)

// set cookie store variables
var (
	// cookie store
	Store *sessions.CookieStore
	// session name
	Name string
)

// Session struct
// Stores session options, name and secret key values
type Session struct {
	Options   sessions.Options `json:"Options"`
	Name      string           `json:"Name"`
	SecretKey string           `json:"SecretKey"`
}

// Configure session
// Setting session values
func Configure(s Session) {
	Store = sessions.NewCookieStore([]byte(s.SecretKey))
	Store.Options = &s.Options
	Name = s.Name
}

// Instance function
// spawn new session instance
func Instance(r *http.Request) *sessions.Session {
	session, _ := Store.Get(r, Name)
	return session
}

// Empty function
// delete all session values
func Empty(sess *sessions.Session) {
	for s := range sess.Values {
		delete(sess.Values, s)
	}
}
