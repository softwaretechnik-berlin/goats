package goinsp

import (
	"regexp"
	"strings"
)

var versionPathSegmentPattern = regexp.MustCompile(`v\d+`)

type ImportPath string

func (p ImportPath) LastNonVersionSegment() string {
	path := string(p)
	i := strings.LastIndex(path, "/")
	if segment := path[i+1:]; !versionPathSegmentPattern.MatchString(segment) {
		return segment
	}
	path = path[:i]
	return path[strings.LastIndex(path, "/")+1:]
}

type TypeName string

func (n TypeName) String() string {
	return string(n)
}
