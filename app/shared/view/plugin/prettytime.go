package plugin

import (
	"html/template"
	"time"
)

// PrettyTime returns a template.FuncMap
// * PRETTYTIME outputs a nice time format
func PrettyTime() template.FuncMap {
	f := make(template.FuncMap)

	f["PRETTYTIME"] = func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	}

	return f
}
