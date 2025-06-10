package servicevelog

import (
	"regexp"
)

func sanitizeBasePathSpecialCase(filename string) (string, bool) {
	re := regexp.MustCompile(`[/]`)
	matched := re.MatchString(filename)
	sanitize := re.ReplaceAllString(filename, "-")
	return sanitize, matched
}
