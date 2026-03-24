package repl

import (
	"strings"
	"unicode"
)

func CleanInput(text string) []string {
	var out []string
	last := -1

	for i, char := range text {
		if unicode.IsSpace(char) {
			if last != -1 {
				out = append(out, strings.ToLower(text[last:i]))
				last = -1
			}
			continue
		}

		if last == -1 {
			last = i
		}
	}

	if last != -1 {
		out = append(out, strings.ToLower(text[last:]))
	}

	return out
}
