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

func markdownPictureMatcher(contents string) []string {
	re := regexp.MustCompile(`!\[[^\]]*\]\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(contents, -1)

	pictures := []string{}
	for _, match := range matches {
		pictures = append(pictures, match[1])
	}

	return pictures
}
