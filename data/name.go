package data

import (
	"regexp"
	"strings"
)

var (
	nameSpaceRegexp, nameNotAlnumRegexp *regexp.Regexp
)

func init() {
	nameSpaceRegexp = regexp.MustCompile(`\s+`)
	nameNotAlnumRegexp = regexp.MustCompile(`[[:^alnum:]]+`)
}

func Name(input string) (name string, slug string) {
	name = nameSpaceRegexp.ReplaceAllString(input, " ")
	name = strings.Trim(name, " ")
	slug = nameNotAlnumRegexp.ReplaceAllString(name, "-")
	slug = strings.Trim(slug, "-")
	slug = strings.ToLower(slug)
	return
}
