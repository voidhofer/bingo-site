// bingosite is a simple but powerful webapp written in go.
// This package is the base of BinGO-Site.
// It inlcudes all the needed packages to build the webserver.

package bingosite

import (
	// base core pkg
	"encoding/json"
	"log"
	"os"
	"runtime"

	// local imports
	"github.com/voidhofer/bingo-site/app/route"
	"github.com/voidhofer/bingo-site/app/shared/jsonconfig"
	"github.com/voidhofer/bingo-site/app/shared/recaptcha"
	"github.com/voidhofer/bingo-site/app/shared/server"
	"github.com/voidhofer/bingo-site/app/shared/session"
	"github.com/voidhofer/bingo-site/app/shared/view"
	"github.com/voidhofer/bingo-site/app/shared/view/plugin"
)

// init sets log flags and allows more CPU core usage
func init() {
	// set log flags
	log.SetFlags(log.Lshortfile)
	// use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// bingosite is the starter function of bingosite package
// It loads the configuration file which is a json file called config.json located in static/config directory.
// The information from the config file gets injected into session, recaptcha, view and server packages.
func bingosite() {
	// load config file
	jsonconfig.Load("static"+string(os.PathSeparator)+"config"+string(os.PathSeparator)+"config.json", config)
	// configure session cookie store
	session.Configure(config.Session)
	// configure Google reCaptcha
	recaptcha.Configure(config.Recaptcha)
	// setup views
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		plugin.Math(),
		plugin.URLyze(),
		plugin.Uppercase(),
		recaptcha.Plugin())
	// start listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

// declare config variable
var config = &configuration{}

// configuration struct
// Holds all the needed information for bingosite function to configure packages.
type configuration struct {
	Recaptcha recaptcha.Info  `json:"Recaptcha"`
	Server    server.Server   `json:"Server"`
	Session   session.Session `json:"Session"`
	Template  view.Template   `json:"Template"`
	View      view.View       `json:"View"`
}

// ParseJSON returns configuration object filled with data retrieved from input.
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
