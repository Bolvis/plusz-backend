package util

import (
	"strconv"
	"strings"
)

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for i := 1; i <= tmpCount; i++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(i), 1)
	}
	return old
}
