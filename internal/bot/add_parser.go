package bot

import (
	"regexp"
	"strings"

	"github.com/IlyasYOY/telethings/internal/thingsurl"
)

// modifierRe matches key:value or key:"quoted value" modifiers.
var modifierRe = regexp.MustCompile(`(when|deadline|tags|notes):("(?:[^"\\]|\\.)*"|\S+)`)

// parseAddCommand parses the text after "/add" and returns the Things3 add URL string.
// Supported modifiers (order-independent, after the title):
//
//	when:<value>
//	deadline:<value>
//	tags:<csv>
//	notes:<word>
//	notes:"quoted text with spaces"
func parseAddCommand(authToken, text string) string {
	title, modifiers := splitTitleAndModifiers(text)
	if title == "" {
		return ""
	}

	u := thingsurl.New(authToken).Add(title)

	if v, ok := modifiers["when"]; ok {
		u = u.WithWhen(v)
	}
	if v, ok := modifiers["deadline"]; ok {
		u = u.WithDeadline(v)
	}
	if v, ok := modifiers["tags"]; ok {
		u = u.WithTags(strings.Split(v, ",")...)
	}
	if v, ok := modifiers["notes"]; ok {
		u = u.WithNotes(v)
	}

	return u.String()
}

// splitTitleAndModifiers separates the plain title from key:value modifiers.
// Title is the substring before the first modifier match.
func splitTitleAndModifiers(text string) (title string, modifiers map[string]string) {
	modifiers = make(map[string]string)
	loc := modifierRe.FindStringIndex(text)
	if loc == nil {
		return strings.TrimSpace(text), modifiers
	}

	title = strings.TrimSpace(text[:loc[0]])

	matches := modifierRe.FindAllStringSubmatch(text, -1)
	for _, m := range matches {
		key := m[1]
		val := m[2]
		// strip surrounding quotes
		if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
			val = val[1 : len(val)-1]
		}
		if _, exists := modifiers[key]; !exists {
			modifiers[key] = val
		}
	}

	return title, modifiers
}
