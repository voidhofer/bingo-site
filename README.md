
# BinGO Site

[![Go Report Card](https://goreportcard.com/badge/github.com/voidhofer/bingo-site)](https://goreportcard.com/report/github.com/voidhofer/bingo-site) [![GoDoc](https://godoc.org/github.com/voidhofer/bingo-site?status.svg)](https://godoc.org/github.com/voidhofer/bingo-site) 

## What is BinGO Site?
A simple but powerful golang website framework (kindof).

This site is capable of serving multiple domains with different routing for each domain.
The program can run in both HTTP and HTTPS mode and it support multiple TLS keys for multiple domains.

## Features
- Multi-domain support (with different routing rules for each domain)
- Multi-domain HTTPS support (using multiple tls configurations for domains)
- Template support
- Template plugin support (use existing plugins or write your own easily)
- MongoDB storage
- Fast and stable httprouter by Julien Smidth
- CSRF protection by Joseph Spurrier
- Google ReCaptcha support by Haisum
- Example templates built with Bootstrap (source included for local hosting)

## Warning!
Since this is a hobby project it is disencouraged to use it in production.
The software comes AS IT IS. There is no support for this software and there will not be no frequent updates, so it might get vulnerable/slow/deprecated as it gets outdated!

## Authors
[voidhofer](github.com/voidhofer) - [voidhofer.github.io](https://voidhofer.github.io)

## Installing with go get
Just simply go get the project and run 'go build' in its folder.
```bash
go get github.com/voidhofer/bingo-site
cd $GOPATH/src/github.com/voidhofer/bingo-site
go build
```

## Installing with git
You can also download or clone the project from github, move it to your $GOPATH and build it.
```bash
cd $GOPATH
cd src/github.com/
mkdir voidhofer
cd voidhofer
git clone https://github.com/voidhofer/bingo-site.git
cd bingo-site
go build
```

## Preparing run folder
After you built the package you can move the compiled file to wherever you like just make sure to copy static folder with it. Static folder contains your config file, tempaltes, js and css files, tls keys, images and favicons. After your web app folder is set up you can run the executable.

## Running on linux:
You should add execution permission to the executable:
```bash
chmod +x ./bingo-site
./bingo-site
```

## Running on Windows:
You might need to run as administrator priviledges.

## Running on OS X:
You should add execution permission to the executable just like linux users.
When you run the executable a popup will be presented. Allow the program to bind to HTTP/HTTPS ports.
```bash
chmod +x ./bingo-site
./bingo-site
```

## Configuration
```
"Database": {
	"MongoDB": {
		"URL": "127.0.0.1",                          // SET MongoDB host (Usually 127.0.0.1 - localhost)
		"Database": "bingo"                          // SET database name
	}
},
"Recaptcha": {
	"Enabled": false,                                    // ENABLE recaptcha
	"Secret": "",                                        // SET secret (the one you get from google)
	"SiteKey": ""                                        // SET sitekey (the one you get from google)
},
"Server": {
	"Hostname": "",                                      // SET hostname (not necessary)
	"UseHTTP": true,                                     // ENABLE HTTP
	"UseHTTPS": false,                                   // ENABLE HTTPS
	"HTTPPort": 80,                                      // SET HTTP port
	"HTTPSPort": 443,                                    // SET HTTPS port
	"CertFile": "static/tls/domain.crt",                 // SET TLS Cert File (only for single domain use!)
	"KeyFile": "static/tls/domain.key"                   // SET TLS Key File (only for single domain use!)
},
"Session": {
	"SecretKey": "VKh37w&=2dqZ&CS3NJaEEf@X?32W3qpr",     // SET session key
	"Name": "SESSI0N",                                   // SET session name
	"Options": {
		"Path": "/",                                 // SET Cookie scope
		"Domain": "",                                // SET Cookie domain
		"MaxAge": 28800,                             // SET Max-age
		"Secure": false,                             // ENABLE SECURE
		"HttpOnly": false                            // ENABLE HTTP
	}
},
"Template": {
	"Root": "base",                                      // SET main template
	"Children": [            
		"partial/menu",				     // SET partial templates (add everything
		"partial/footer"			     // you want to use on multiple pages)
	]
},
"View": {
	"BaseURI": "/",                                      // SET BaseURI
	"Extension": "tmpl",                                 // SET template extension
	"Folder": "static/template",                         // SET template folder
	"Caching": false                                     // ENABLE caching
}
```

## Multiple domain with multiple TLS config
Add domain in app/route/route.go like this:
```go
// Setting domains
var domain1 = "example1.com"
var domain2 = "example2.com"
var domain3 = "example3.com"
```

Then create an httprouter for each domain:
```go
// Creating httprouters for each domain
r1 := httprouter.New()
r2 := httprouter.New()
r3 := httprouter.New()
```

Add routing rules for all domains:
```go
// Serve /about for example1.com
r1.GET("/about", hr.Handler(alice.
	New().
	ThenFunc(controller.AboutGET)))
// Serve /about for example2.com
r2.GET("/about", hr.Handler(alice.
	New().
	ThenFunc(controller.AboutGET2)))
// Serve /about for example3.com
r3.GET("/about", hr.Handler(alice.
	New().				    // You can use the same 
	ThenFunc(controller.AboutGET2)))    // controller for multiple routing rules
```

Add the httproutes to hostswitch map:
```go
// set hostswitch value for routes
hs[domain1] = r1
hs[domain2] = r2
hs[domain3] = r3
```

At last in app/shared/server/server.go you can add multiple tls configurations (if HTTPS is enabled):
```go
// Setting tls conf for domain1
cert, err1 := tls.LoadX509KeyPair("static/tls/domain.crt", "static/tls/domain.key")
if err1 != nil {
	log.Fatal(err1)
}
// Appending key and crt for cfg.Certificates (domain1)
cfg.Certificates = append(cfg.Certificates, cert)
// Setting tls conf for domain2
cert2, err2 := tls.LoadX509KeyPair("static/tls/domain2.crt", "static/tls/domain2.key")
if err2 != nil {
	log.Fatal(err2)
}
// Appending key and crt for cfg.Certificates (domain2)
cfg.Certificates = append(cfg.Certificates, cert2)
// ...
```

## Built with:
[Gorilla Context](github.com/gorilla/context)

[Gorilla Sessions](github.com/gorilla/sessions)

[Julien Schmidt's HTTP router](github.com/julienschmidt/httprouter)

[Justinas's Alice - middleware chaining](github.com/justinas/alice)

[Joseph Spurrier's CSRFbanana](github.com/josephspurrier/csrfbanana)

[Haisum's implementation of Google recaptcha](github.com/haisum/recaptcha)


## Acknowledgments: 
This project was inspired by Joseph Spurrier's blog post about creating Go web applications.
I used some of his packages and change them a little bit so they fit my needs. I would suggest you all to read his article and take a look at his github because he did an excelent job explaining his concept of golang web development.

[Joseph Spurrier's blog post - Go Web App Example](http://www.josephspurrier.com/go-web-app-example/)
