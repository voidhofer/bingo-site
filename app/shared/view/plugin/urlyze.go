package plugin

import (
	"html/template"
	"strings"
	"regexp"
	"log"
)

// Math returns a template.FuncMap
// * MATH is a plugin to do math in templates
func URLyze() template.FuncMap {
	f := make(template.FuncMap)

	f["URLYZE"] = func(text string) string {
		var result string
		text = strings.ToLower(text)
		text = strings.Replace(text, "á", "a", -1)
		text = strings.Replace(text, "é", "e", -1)
		text = strings.Replace(text, "í", "i", -1)
		text = strings.Replace(text, "ó", "o", -1)
		text = strings.Replace(text, "ö", "o", -1)
		text = strings.Replace(text, "ő", "o", -1)
		text = strings.Replace(text, "ú", "u", -1)
		text = strings.Replace(text, "ü", "u", -1)
		text = strings.Replace(text, "ű", "u", -1)
		text = strings.Replace(text, " ", "_", -1)
		text = strings.Replace(text, "-", "_", -1)
		reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	    if err != nil {
	        log.Fatal(err)
	    }
	    result = reg.ReplaceAllString(text, "")
		return result
	}

	return f
}
