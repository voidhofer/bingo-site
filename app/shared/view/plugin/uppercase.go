package plugin

import (
	"html/template"
	"strings"
)

// Math returns a template.FuncMap
// * MATH is a plugin to do math in templates
func Uppercase() template.FuncMap {
	f := make(template.FuncMap)

	f["UPPERCASE"] = func(text string) string {
		var result string
		result = strings.ToUpper(text)
		return result
	}

	return f
}
