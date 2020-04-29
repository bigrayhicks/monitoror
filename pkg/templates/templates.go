// Based on https://github.com/docker/cli/blob/master/templates/templates.go

package templates

import (
	"strings"
	"text/template"

	"github.com/labstack/gommon/color"
)

var basicFunctions = template.FuncMap{
	"split": strings.Split,
	"join":  strings.Join,
	"lower": strings.ToLower,
	"upper": strings.ToUpper,

	// For terminal only
	"blue":         color.Blue,
	"green":        color.Green,
	"red":          color.Red,
	"yellow":       color.Yellow,
	"grey":         color.Grey,
	"inverseColor": color.Inverse,
}

// Parse creates a new anonymous template with the basic functions
// and parses the given format.
func Parse(format string) (*template.Template, error) {
	return NewParse("", format)
}

// New creates a new empty template with the provided tag and built-in
// template functions.
func New(tag string) *template.Template {
	return template.New(tag).Funcs(basicFunctions)
}

// NewParse creates a new tagged template with the basic functions
// and parses the given format.
func NewParse(tag, format string) (*template.Template, error) {
	return New(tag).Parse(format)
}
