package service

import (
	"regexp"
)

func sanitizeBasePathSpecialCase(filename string) string {
	re := regexp.MustCompile(`[/]`)
	return re.ReplaceAllString(filename, "-")
}
