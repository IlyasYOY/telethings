package thingsurl

import (
	"net/url"
	"strings"
)

type addProjectThingsURL struct {
	*thingsURL

	title          string
	notes          string
	when           WhenValue
	deadline       string
	tags           []string
	areaID         string
	area           string
	todos          []string
	completed      bool
	canceled       bool
	reveal         bool
	creationDate   string
	completionDate string
}

func (tu *thingsURL) AddProject(title string) *addProjectThingsURL {
	return &addProjectThingsURL{
		thingsURL: tu,
		title:     title,
	}
}

func (apu *addProjectThingsURL) WithNotes(notes string) *addProjectThingsURL {
	apu.notes = notes
	return apu
}

func (apu *addProjectThingsURL) WithWhen(when WhenValue) *addProjectThingsURL {
	apu.when = when
	return apu
}

func (apu *addProjectThingsURL) WithDeadline(deadline string) *addProjectThingsURL {
	apu.deadline = deadline
	return apu
}

func (apu *addProjectThingsURL) WithTags(tags ...string) *addProjectThingsURL {
	apu.tags = tags
	return apu
}

func (apu *addProjectThingsURL) WithAreaID(id string) *addProjectThingsURL {
	apu.areaID = id
	return apu
}

func (apu *addProjectThingsURL) WithArea(area string) *addProjectThingsURL {
	apu.area = area
	return apu
}

func (apu *addProjectThingsURL) WithTodos(todos ...string) *addProjectThingsURL {
	apu.todos = todos
	return apu
}

func (apu *addProjectThingsURL) Completed() *addProjectThingsURL {
	apu.completed = true
	return apu
}

func (apu *addProjectThingsURL) Canceled() *addProjectThingsURL {
	apu.canceled = true
	return apu
}

func (apu *addProjectThingsURL) Reveal() *addProjectThingsURL {
	apu.reveal = true
	return apu
}

func (apu *addProjectThingsURL) WithCreationDate(date string) *addProjectThingsURL {
	apu.creationDate = date
	return apu
}

func (apu *addProjectThingsURL) WithCompletionDate(date string) *addProjectThingsURL {
	apu.completionDate = date
	return apu
}

func (apu *addProjectThingsURL) String() string {
	params := url.Values{}
	params.Set("title", apu.title)
	if apu.notes != "" {
		params.Set("notes", apu.notes)
	}
	if apu.when != "" {
		params.Set("when", apu.when)
	}
	if apu.deadline != "" {
		params.Set("deadline", apu.deadline)
	}
	if len(apu.tags) > 0 {
		params.Set("tags", strings.Join(apu.tags, ","))
	}
	if apu.areaID != "" {
		params.Set("area-id", apu.areaID)
	} else if apu.area != "" {
		params.Set("area", apu.area)
	}
	if len(apu.todos) > 0 {
		params.Set("to-dos", strings.Join(apu.todos, "\n"))
	}
	if apu.completed {
		params.Set("completed", "true")
	}
	if apu.canceled {
		params.Set("canceled", "true")
	}
	if apu.reveal {
		params.Set("reveal", "true")
	}
	if apu.creationDate != "" {
		params.Set("creation-date", apu.creationDate)
	}
	if apu.completionDate != "" {
		params.Set("completion-date", apu.completionDate)
	}
	return "things:///add-project?" + encodeParams(params)
}
