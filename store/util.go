package store

import (
	"strings"
)

func EscapeNewlines(s string) string {
	var escapeReplacer = strings.NewReplacer(
		"\\", "\\\\",
		"\n", "\\n",
		"\r", "\\r",
	)

	return escapeReplacer.Replace(s)
}


func UnescapeNewlines(s string) string {
	var UnescapeReplacer = strings.NewReplacer(
		"\\n", "\n",
		"\\r", "\r",
		"\\\\", "\\",
	)

	return UnescapeReplacer.Replace(s)
}