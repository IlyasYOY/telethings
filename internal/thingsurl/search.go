package thingsurl

import "net/url"

type searchThingsURL struct {
	*thingsURL

	query string
}

func (tu *thingsURL) Search(query string) *searchThingsURL {
	return &searchThingsURL{
		thingsURL: tu,
		query:     query,
	}
}

func (stu *searchThingsURL) String() string {
	if stu.query == "" {
		return "things:///search"
	}
	params := url.Values{}
	params.Set("query", stu.query)
	return "things:///search?" + encodeParams(params)
}
