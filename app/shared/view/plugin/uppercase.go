package plugin

import (
	"html/template"
	"strings"
)

// Uppercase returns a template.FuncMap
// * returns uppercase version of string
func Uppercase() template.FuncMap {
	f := make(template.FuncMap)

	f["UPPERCASE"] = func(text string) string {
		var result string
		result = strings.ToUpper(text)
		return result
	}

	return f
}
