package utils

import (
	"regexp"
	"strings"
)

func HideMobile(mobile string) string {
	re := regexp.MustCompile("([0-9]{3})[0-9]+([0-9]{4})")
	return re.ReplaceAllString(mobile, "$1****$2")
}

func Normalize(str string) string {
	return strings.TrimSpace(strings.ToLower(str))
}
