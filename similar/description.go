package similar

import (
	"regexp"
	"strings"
)

var numberRE = regexp.MustCompile("(#[0-9]*)")

func Sanitize(description string) string {
	description = strings.ToLower(strings.TrimSpace(description))
	if pos := strings.Index(description, "  "); pos > 0 {
		description = description[:pos]
	}
	if numberRE.MatchString(description) {
		description = numberRE.ReplaceAllLiteralString(description, "")
	}
	return description
}
