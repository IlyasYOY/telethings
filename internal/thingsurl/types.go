package thingsurl

import (
	"net/url"
	"strings"
)

// WhenValue is a string value for the "when" parameter.
// Use predefined constants or a date string (yyyy-mm-dd)
// or date-time string (yyyy-mm-dd@HH:MM).
type WhenValue = string

const (
	WhenToday    WhenValue = "today"
	WhenTomorrow WhenValue = "tomorrow"
	WhenEvening  WhenValue = "evening"
	WhenAnytime  WhenValue = "anytime"
	WhenSomeday  WhenValue = "someday"
)

// ShowListID represents a built-in Things3 list identifier for the show command.
type ShowListID = string

const (
	ShowListInbox          ShowListID = "inbox"
	ShowListToday          ShowListID = "today"
	ShowListAnytime        ShowListID = "anytime"
	ShowListUpcoming       ShowListID = "upcoming"
	ShowListSomeday        ShowListID = "someday"
	ShowListLogbook        ShowListID = "logbook"
	ShowListTomorrow       ShowListID = "tomorrow"
	ShowListDeadlines      ShowListID = "deadlines"
	ShowListRepeating      ShowListID = "repeating"
	ShowListAllProjects    ShowListID = "all-projects"
	ShowListLoggedProjects ShowListID = "logged-projects"
)

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// encodeParams encodes url.Values using %20 for spaces (Things3 uses %20, not +).
func encodeParams(params url.Values) string {
	return strings.ReplaceAll(params.Encode(), "+", "%20")
}
