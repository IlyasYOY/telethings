package thingsurl

import (
	"net/url"
	"strings"
)

type updateProjectThingsURL struct {
	*thingsURL

	id             string
	title          string
	notes          *string
	prependNotes   string
	appendNotes    string
	when           WhenValue
	deadline       *string
	tags           []string
	addTags        []string
	areaID         string
	area           string
	completed      *bool
	canceled       *bool
	reveal         bool
	duplicate      bool
	creationDate   string
	completionDate string
}

func (tu *thingsURL) UpdateProject(id string) *updateProjectThingsURL {
	return &updateProjectThingsURL{
		thingsURL: tu,
		id:        id,
	}
}

func (upu *updateProjectThingsURL) WithTitle(title string) *updateProjectThingsURL {
	upu.title = title
	return upu
}

// WithNotes replaces the existing notes.
func (upu *updateProjectThingsURL) WithNotes(notes string) *updateProjectThingsURL {
	upu.notes = &notes
	return upu
}

// ClearNotes clears the notes field.
func (upu *updateProjectThingsURL) ClearNotes() *updateProjectThingsURL {
	empty := ""
	upu.notes = &empty
	return upu
}

func (upu *updateProjectThingsURL) WithPrependNotes(text string) *updateProjectThingsURL {
	upu.prependNotes = text
	return upu
}

func (upu *updateProjectThingsURL) WithAppendNotes(text string) *updateProjectThingsURL {
	upu.appendNotes = text
	return upu
}

func (upu *updateProjectThingsURL) WithWhen(when WhenValue) *updateProjectThingsURL {
	upu.when = when
	return upu
}

// WithDeadline sets the deadline. Pass an empty string to clear it.
func (upu *updateProjectThingsURL) WithDeadline(deadline string) *updateProjectThingsURL {
	upu.deadline = &deadline
	return upu
}

func (upu *updateProjectThingsURL) WithTags(tags ...string) *updateProjectThingsURL {
	upu.tags = tags
	return upu
}

func (upu *updateProjectThingsURL) WithAddTags(tags ...string) *updateProjectThingsURL {
	upu.addTags = tags
	return upu
}

func (upu *updateProjectThingsURL) WithAreaID(id string) *updateProjectThingsURL {
	upu.areaID = id
	return upu
}

func (upu *updateProjectThingsURL) WithArea(area string) *updateProjectThingsURL {
	upu.area = area
	return upu
}

func (upu *updateProjectThingsURL) Completed() *updateProjectThingsURL {
	v := true
	upu.completed = &v
	return upu
}

func (upu *updateProjectThingsURL) Uncompleted() *updateProjectThingsURL {
	v := false
	upu.completed = &v
	return upu
}

func (upu *updateProjectThingsURL) Canceled() *updateProjectThingsURL {
	v := true
	upu.canceled = &v
	return upu
}

func (upu *updateProjectThingsURL) Uncanceled() *updateProjectThingsURL {
	v := false
	upu.canceled = &v
	return upu
}

func (upu *updateProjectThingsURL) Reveal() *updateProjectThingsURL {
	upu.reveal = true
	return upu
}

func (upu *updateProjectThingsURL) Duplicate() *updateProjectThingsURL {
	upu.duplicate = true
	return upu
}

func (upu *updateProjectThingsURL) WithCreationDate(date string) *updateProjectThingsURL {
	upu.creationDate = date
	return upu
}

func (upu *updateProjectThingsURL) WithCompletionDate(date string) *updateProjectThingsURL {
	upu.completionDate = date
	return upu
}

func (upu *updateProjectThingsURL) String() string {
	params := url.Values{}
	params.Set("auth-token", upu.authToken)
	params.Set("id", upu.id)
	if upu.title != "" {
		params.Set("title", upu.title)
	}
	if upu.notes != nil {
		params.Set("notes", *upu.notes)
	}
	if upu.prependNotes != "" {
		params.Set("prepend-notes", upu.prependNotes)
	}
	if upu.appendNotes != "" {
		params.Set("append-notes", upu.appendNotes)
	}
	if upu.when != "" {
		params.Set("when", upu.when)
	}
	if upu.deadline != nil {
		params.Set("deadline", *upu.deadline)
	}
	if len(upu.tags) > 0 {
		params.Set("tags", strings.Join(upu.tags, ","))
	}
	if len(upu.addTags) > 0 {
		params.Set("add-tags", strings.Join(upu.addTags, ","))
	}
	if upu.areaID != "" {
		params.Set("area-id", upu.areaID)
	} else if upu.area != "" {
		params.Set("area", upu.area)
	}
	if upu.completed != nil {
		params.Set("completed", boolToString(*upu.completed))
	}
	if upu.canceled != nil {
		params.Set("canceled", boolToString(*upu.canceled))
	}
	if upu.reveal {
		params.Set("reveal", "true")
	}
	if upu.duplicate {
		params.Set("duplicate", "true")
	}
	if upu.creationDate != "" {
		params.Set("creation-date", upu.creationDate)
	}
	if upu.completionDate != "" {
		params.Set("completion-date", upu.completionDate)
	}
	return "things:///update-project?" + encodeParams(params)
}
