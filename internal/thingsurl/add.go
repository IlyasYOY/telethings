package thingsurl

import (
	"net/url"
	"strings"
)

type addThingsURL struct {
	*thingsURL

	title          string
	titles         []string
	notes          string
	when           WhenValue
	deadline       string
	tags           []string
	checklistItems []string
	useClipboard   string
	listID         string
	list           string
	headingID      string
	heading        string
	completed      bool
	canceled       bool
	showQuickEntry bool
	reveal         bool
	creationDate   string
	completionDate string
}

// TODO: make title varargs and remove WithTitles method.
func (tu *thingsURL) Add(title string) *addThingsURL {
	return &addThingsURL{
		thingsURL: tu,
		title:     title,
	}
}

func (atu *addThingsURL) Completed() *addThingsURL {
	atu.completed = true
	return atu
}

func (atu *addThingsURL) Canceled() *addThingsURL {
	atu.canceled = true
	return atu
}

func (atu *addThingsURL) WithNotes(notes string) *addThingsURL {
	atu.notes = notes
	return atu
}

func (atu *addThingsURL) WithWhen(when WhenValue) *addThingsURL {
	atu.when = when
	return atu
}

func (atu *addThingsURL) WithDeadline(deadline string) *addThingsURL {
	atu.deadline = deadline
	return atu
}

func (atu *addThingsURL) WithTags(tags ...string) *addThingsURL {
	atu.tags = tags
	return atu
}

func (atu *addThingsURL) WithChecklistItems(items ...string) *addThingsURL {
	atu.checklistItems = items
	return atu
}

func (atu *addThingsURL) WithTitles(titles ...string) *addThingsURL {
	atu.titles = titles
	return atu
}

func (atu *addThingsURL) WithListID(id string) *addThingsURL {
	atu.listID = id
	return atu
}

func (atu *addThingsURL) WithList(list string) *addThingsURL {
	atu.list = list
	return atu
}

func (atu *addThingsURL) WithHeadingID(id string) *addThingsURL {
	atu.headingID = id
	return atu
}

func (atu *addThingsURL) WithHeading(heading string) *addThingsURL {
	atu.heading = heading
	return atu
}

func (atu *addThingsURL) ShowQuickEntry() *addThingsURL {
	atu.showQuickEntry = true
	return atu
}

func (atu *addThingsURL) Reveal() *addThingsURL {
	atu.reveal = true
	return atu
}

func (atu *addThingsURL) WithCreationDate(date string) *addThingsURL {
	atu.creationDate = date
	return atu
}

func (atu *addThingsURL) WithCompletionDate(date string) *addThingsURL {
	atu.completionDate = date
	return atu
}

func (atu *addThingsURL) WithUseClipboard(mode string) *addThingsURL {
	atu.useClipboard = mode
	return atu
}

func (atu *addThingsURL) String() string {
	params := url.Values{}
	params.Set("auth-token", atu.authToken)
	if len(atu.titles) > 0 {
		params.Set("titles", strings.Join(atu.titles, "\n"))
	} else {
		params.Set("title", atu.title)
	}
	if atu.notes != "" {
		params.Set("notes", atu.notes)
	}
	if atu.when != "" {
		params.Set("when", atu.when)
	}
	if atu.deadline != "" {
		params.Set("deadline", atu.deadline)
	}
	if len(atu.tags) > 0 {
		params.Set("tags", strings.Join(atu.tags, ","))
	}
	if len(atu.checklistItems) > 0 {
		params.Set("checklist-items", strings.Join(atu.checklistItems, "\n"))
	}
	if atu.useClipboard != "" {
		params.Set("use-clipboard", atu.useClipboard)
	}
	if atu.listID != "" {
		params.Set("list-id", atu.listID)
	} else if atu.list != "" {
		params.Set("list", atu.list)
	}
	if atu.headingID != "" {
		params.Set("heading-id", atu.headingID)
	} else if atu.heading != "" {
		params.Set("heading", atu.heading)
	}
	if atu.completed {
		params.Set("completed", "true")
	}
	if atu.canceled {
		params.Set("canceled", "true")
	}
	if atu.showQuickEntry {
		params.Set("show-quick-entry", "true")
	}
	if atu.reveal {
		params.Set("reveal", "true")
	}
	if atu.creationDate != "" {
		params.Set("creation-date", atu.creationDate)
	}
	if atu.completionDate != "" {
		params.Set("completion-date", atu.completionDate)
	}
	return "things:///add?" + encodeParams(params)
}
