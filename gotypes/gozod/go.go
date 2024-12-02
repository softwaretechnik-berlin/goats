package gozod

import (
	"strings"
)

func tagHasFlag(tag string, flag string) bool {
	_, tag, _ = strings.Cut(tag, ",")
	for tag != "" {
		var currentFlag string
		currentFlag, tag, _ = strings.Cut(tag, ",")
		if currentFlag == flag {
			return true
		}
	}
	return false
}
