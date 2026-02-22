package thingsurl

import (
	"net/url"
	"strings"
)

type showThingsURL struct {
	*thingsURL

	id     string
	query  string
	filter []string
}

func (tu *thingsURL) Show(id string) *showThingsURL {
	return &showThingsURL{
		thingsURL: tu,
		id:        id,
	}
}

// ShowByQuery creates a show command that uses the query parameter instead of an ID.
func (tu *thingsURL) ShowByQuery(query string) *showThingsURL {
	return &showThingsURL{
		thingsURL: tu,
		query:     query,
	}
}

func (stu *showThingsURL) WithFilter(tags ...string) *showThingsURL {
	stu.filter = tags
	return stu
}

func (stu *showThingsURL) String() string {
	params := url.Values{}
	if stu.id != "" {
		params.Set("id", stu.id)
	} else if stu.query != "" {
		params.Set("query", stu.query)
	}
	if len(stu.filter) > 0 {
		params.Set("filter", strings.Join(stu.filter, ","))
	}
	return "things:///show?" + encodeParams(params)
}
