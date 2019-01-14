package controller

import (
	"net/http"

	"github.com/voidhofer/bingo-site/app/shared/view"
)

// MemberAreaGET displays protected area page
func MemberAreaGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "member/area"
	v.Render(w)
}
