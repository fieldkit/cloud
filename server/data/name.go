package data

import (
	"regexp"
	"strings"
)

var (
	nameSpaces *regexp.Regexp
)

func init() {
	nameSpaces = regexp.MustCompile(`\s+`)
}

func Name(name string) string {
	return strings.TrimSpace(nameSpaces.ReplaceAllString(name, " "))
}
