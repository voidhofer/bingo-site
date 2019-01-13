package plugin

import (
	"html/template"
)

// Math returns a template.FuncMap
// * MATH is a plugin to do math in templates
func Math() template.FuncMap {
	f := make(template.FuncMap)

	f["MATH"] = func(i1 int, mod int, i2 int) int {
		var result int
		switch mod {
			case 0:
				result = (i1 + i2)
			case 1:
				result = (i1 - i2)
			case 2:
				result = (i1 * i2)
			case 3:
				result = int(i1 / i2)
		}
		return result
	}

	return f
}
