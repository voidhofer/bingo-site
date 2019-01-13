package controller

import (
	"net/http"

	"github.com/voidhofer/bingo_site/app/shared/view"
)

// AboutGET displays the About page
func AboutGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "about/about"
	v.Render(w)
}

// AboutGET displays the About page
func AboutGET2(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "about/about2"
	v.Render(w)
}
