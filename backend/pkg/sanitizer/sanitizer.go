package sanitizer

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

var policy = bluemonday.StrictPolicy()

// HTML strips all HTML tags from the input, preventing XSS.
func HTML(input string) string {
	return strings.TrimSpace(policy.Sanitize(input))
}

// Text trims whitespace and strips HTML.
func Text(input string) string {
	return HTML(input)
}
