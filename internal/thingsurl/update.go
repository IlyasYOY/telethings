package thingsurl

import (
	"net/url"
	"strings"
)

type updateThingsURL struct {
	*thingsURL

	id                    string
	title                 string
	notes                 *string
	prependNotes          string
	appendNotes           string
	when                  WhenValue
	deadline              *string
	tags                  []string
	addTags               []string
	checklistItems        []string
	prependChecklistItems []string
	appendChecklistItems  []string
	listID                string
	list                  string
	headingID             string
	heading               string
	completed             *bool
	canceled              *bool
	reveal                bool
	duplicate             bool
	creationDate          string
	completionDate        string
}

func (tu *thingsURL) Update(id string) *updateThingsURL {
	return &updateThingsURL{
		thingsURL: tu,
		id:        id,
	}
}

func (utu *updateThingsURL) WithTitle(title string) *updateThingsURL {
	utu.title = title
	return utu
}

// WithNotes replaces the existing notes.
func (utu *updateThingsURL) WithNotes(notes string) *updateThingsURL {
	utu.notes = &notes
	return utu
}

// ClearNotes clears the notes field.
func (utu *updateThingsURL) ClearNotes() *updateThingsURL {
	empty := ""
	utu.notes = &empty
	return utu
}

func (utu *updateThingsURL) WithPrependNotes(text string) *updateThingsURL {
	utu.prependNotes = text
	return utu
}

func (utu *updateThingsURL) WithAppendNotes(text string) *updateThingsURL {
	utu.appendNotes = text
	return utu
}

func (utu *updateThingsURL) WithWhen(when WhenValue) *updateThingsURL {
	utu.when = when
	return utu
}

// WithDeadline sets the deadline. Pass an empty string to clear it.
func (utu *updateThingsURL) WithDeadline(deadline string) *updateThingsURL {
	utu.deadline = &deadline
	return utu
}

func (utu *updateThingsURL) WithTags(tags ...string) *updateThingsURL {
	utu.tags = tags
	return utu
}

func (utu *updateThingsURL) WithAddTags(tags ...string) *updateThingsURL {
	utu.addTags = tags
	return utu
}

func (utu *updateThingsURL) WithChecklistItems(items ...string) *updateThingsURL {
	utu.checklistItems = items
	return utu
}

func (utu *updateThingsURL) WithPrependChecklistItems(items ...string) *updateThingsURL {
	utu.prependChecklistItems = items
	return utu
}

func (utu *updateThingsURL) WithAppendChecklistItems(items ...string) *updateThingsURL {
	utu.appendChecklistItems = items
	return utu
}

func (utu *updateThingsURL) WithListID(id string) *updateThingsURL {
	utu.listID = id
	return utu
}

func (utu *updateThingsURL) WithList(list string) *updateThingsURL {
	utu.list = list
	return utu
}

func (utu *updateThingsURL) WithHeadingID(id string) *updateThingsURL {
	utu.headingID = id
	return utu
}

func (utu *updateThingsURL) WithHeading(heading string) *updateThingsURL {
	utu.heading = heading
	return utu
}

func (utu *updateThingsURL) Completed() *updateThingsURL {
	v := true
	utu.completed = &v
	return utu
}

func (utu *updateThingsURL) Uncompleted() *updateThingsURL {
	v := false
	utu.completed = &v
	return utu
}

func (utu *updateThingsURL) Canceled() *updateThingsURL {
	v := true
	utu.canceled = &v
	return utu
}

func (utu *updateThingsURL) Uncanceled() *updateThingsURL {
	v := false
	utu.canceled = &v
	return utu
}

func (utu *updateThingsURL) Reveal() *updateThingsURL {
	utu.reveal = true
	return utu
}

func (utu *updateThingsURL) Duplicate() *updateThingsURL {
	utu.duplicate = true
	return utu
}

func (utu *updateThingsURL) WithCreationDate(date string) *updateThingsURL {
	utu.creationDate = date
	return utu
}

func (utu *updateThingsURL) WithCompletionDate(date string) *updateThingsURL {
	utu.completionDate = date
	return utu
}

func (utu *updateThingsURL) String() string {
	params := url.Values{}
	params.Set("auth-token", utu.authToken)
	params.Set("id", utu.id)
	if utu.title != "" {
		params.Set("title", utu.title)
	}
	if utu.notes != nil {
		params.Set("notes", *utu.notes)
	}
	if utu.prependNotes != "" {
		params.Set("prepend-notes", utu.prependNotes)
	}
	if utu.appendNotes != "" {
		params.Set("append-notes", utu.appendNotes)
	}
	if utu.when != "" {
		params.Set("when", utu.when)
	}
	if utu.deadline != nil {
		params.Set("deadline", *utu.deadline)
	}
	if len(utu.tags) > 0 {
		params.Set("tags", strings.Join(utu.tags, ","))
	}
	if len(utu.addTags) > 0 {
		params.Set("add-tags", strings.Join(utu.addTags, ","))
	}
	if len(utu.checklistItems) > 0 {
		params.Set("checklist-items", strings.Join(utu.checklistItems, "\n"))
	}
	if len(utu.prependChecklistItems) > 0 {
		params.Set("prepend-checklist-items", strings.Join(utu.prependChecklistItems, "\n"))
	}
	if len(utu.appendChecklistItems) > 0 {
		params.Set("append-checklist-items", strings.Join(utu.appendChecklistItems, "\n"))
	}
	if utu.listID != "" {
		params.Set("list-id", utu.listID)
	} else if utu.list != "" {
		params.Set("list", utu.list)
	}
	if utu.headingID != "" {
		params.Set("heading-id", utu.headingID)
	} else if utu.heading != "" {
		params.Set("heading", utu.heading)
	}
	if utu.completed != nil {
		params.Set("completed", boolToString(*utu.completed))
	}
	if utu.canceled != nil {
		params.Set("canceled", boolToString(*utu.canceled))
	}
	if utu.reveal {
		params.Set("reveal", "true")
	}
	if utu.duplicate {
		params.Set("duplicate", "true")
	}
	if utu.creationDate != "" {
		params.Set("creation-date", utu.creationDate)
	}
	if utu.completionDate != "" {
		params.Set("completion-date", utu.completionDate)
	}
	return "things:///update?" + encodeParams(params)
}
