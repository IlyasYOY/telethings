package bot

import (
	"regexp"
	"strings"
)

// modifierRe matches key:value or key:"quoted value" modifiers.
var modifierRe = regexp.MustCompile(`(when|deadline|tags|notes):("(?:[^"\\]|\\.)*"|\S+)`)

type addCommandInput struct {
	title    string
	when     string
	deadline string
	tags     []string
	notes    string
}

func parseAddCommandInput(text string) *addCommandInput {
	title, modifiers := splitTitleAndModifiers(text)
	if title == "" {
		return nil
	}
	in := &addCommandInput{title: title}
	if v, ok := modifiers["when"]; ok {
		in.when = v
	}
	if v, ok := modifiers["deadline"]; ok {
		in.deadline = v
	}
	if v, ok := modifiers["tags"]; ok {
		in.tags = strings.Split(v, ",")
	}
	if v, ok := modifiers["notes"]; ok {
		in.notes = v
	}
	return in
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
