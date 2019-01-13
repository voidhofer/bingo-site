package main

import (
	// base core pkg
	"encoding/json"
	"log"
	"os"
	"runtime"

	// local imports
	"github.com/voidhofer/bingo_site/app/route"
	"github.com/voidhofer/bingo_site/app/shared/jsonconfig"
	"github.com/voidhofer/bingo_site/app/shared/recaptcha"
	"github.com/voidhofer/bingo_site/app/shared/server"
	"github.com/voidhofer/bingo_site/app/shared/session"
	"github.com/voidhofer/bingo_site/app/shared/view"
	"github.com/voidhofer/bingo_site/app/shared/view/plugin"
)

func init() {
	// set log flags
	log.SetFlags(log.Lshortfile)
	// use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
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
type configuration struct {
	Recaptcha recaptcha.Info  `json:"Recaptcha"`
	Server    server.Server   `json:"Server"`
	Session   session.Session `json:"Session"`
	Template  view.Template   `json:"Template"`
	View      view.View       `json:"View"`
}

// bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
