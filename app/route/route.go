package route

import (
	// base core imports
	"fmt"
	"net/http"
	"strings"

	// local imports
	"github.com/voidhofer/bingo_site/app/controller"
	hr "github.com/voidhofer/bingo_site/app/shared/middleware/httprouterwrapper"
	"github.com/voidhofer/bingo_site/app/shared/middleware/logrequest"
	"github.com/voidhofer/bingo_site/app/shared/session"

	// external imports
	"github.com/gorilla/context"
	"github.com/josephspurrier/csrfbanana"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Add more domains for multi domain serving with different sites
var domain1 = "127.0.0.1"

// Load returns the routes and middleware
func Load() http.Handler {
	return middleware(routes())
}

// LoadHTTPS returns the HTTP routes and middleware
func LoadHTTPS() http.Handler {
	return middleware(routes())
}

// LoadHTTP returns the HTTPS routes and middleware
func LoadHTTP() http.Handler {
	return middleware(routes())

	// Uncomment this and comment out the line above to always redirect to HTTPS
	//return http.HandlerFunc(redirectToHTTPS)
}

// Optional method to make it easy to redirect from HTTP to HTTPS
func redirectToHTTPS(w http.ResponseWriter, req *http.Request) {
	newHost := strings.Replace(req.Host, ":80", ":443", 1)
	http.Redirect(w, req, "https://"+newHost, http.StatusMovedPermanently)
}

// *****************************************************************************
// Routes
// *****************************************************************************

// Declare HostSwitch type
type HostSwitch map[string]http.Handler

// ServeHTTP method for HostSwitch
func (hs HostSwitch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Check if a http.Handler is registered for the given host.
	// If yes, use it to handle the request.
	if handler := hs[getHost(r)]; handler != nil {
		handler.ServeHTTP(w, r)
	} else {
		fmt.Println("Cannot find domains: ")
		fmt.Println(getHost(r))
		http.Error(w, "Forbidden", 403)
	}
}

// getHost returns the requested hostname without PORT number and "www." prefix
func getHost(r *http.Request) string {
	if r.URL.IsAbs() {
		host := r.Host
		if n := strings.Index(host, "www."); n != -1 {
			host = host[4:]
		}
		// Slice off any port information.
		if i := strings.Index(host, ":"); i != -1 {
			host = host[:i]
		}
		return host
	} else {
		host := r.Host
		if n := strings.Index(host, "www."); n != -1 {
			host = host[4:]
		}
		// Slice off any port information.
		if i := strings.Index(host, ":"); i != -1 {
			if i2 := strings.Index(host, "/"); i2 != -1 {
				host = host[:i] + host[i2:]
			} else {
				host = host[:i]
			}
		}
		return host
	}
	return r.URL.Host
}

// routes lists the available urls
func routes() *HostSwitch {
	// new routers
	r := httprouter.New()
	hs := make(HostSwitch)
	// COPYGURU.HU
	// serve static files, no directory browsing
	r.GET("/static/*filepath", hr.Handler(alice.
		New().
		ThenFunc(controller.Static)))
	// home
	r.GET("/", hr.Handler(alice.
		New().
		ThenFunc(controller.IndexGET)))
	// 404
	r.NotFound = alice.
		New().
		ThenFunc(controller.Error404)
	// login/register page
	r.GET("/login", hr.Handler(alice.
		New().
		ThenFunc(controller.LoginGET)))
	// login post
	r.POST("/login", hr.Handler(alice.
		New().
		ThenFunc(controller.LoginPOST)))
	// register post
	r.POST("/register", hr.Handler(alice.
		New().
		ThenFunc(controller.RegisterPOST)))
	// about
	r.GET("/about", hr.Handler(alice.
		New().
		ThenFunc(controller.AboutGET)))
	// set hostswitch value for routes
	hs[domain1] = r
	// respond
	return &hs
}

// *****************************************************************************
// Middleware
// *****************************************************************************
func middleware(h http.Handler) http.Handler {
	// Prevents CSRF and Double Submits
	cs := csrfbanana.New(h, session.Store, session.Name)
	cs.FailureHandler(http.HandlerFunc(controller.InvalidToken))
	cs.ClearAfterUsage(true)
	cs.ExcludeRegexPaths([]string{"/static(.*)"})
	csrfbanana.TokenLength = 32
	csrfbanana.TokenName = "token"
	csrfbanana.SingleToken = false
	h = cs
	// Log every request
	h = logrequest.Handler(h)
	// Clear handler for Gorilla Context
	h = context.ClearHandler(h)
	return h
}
