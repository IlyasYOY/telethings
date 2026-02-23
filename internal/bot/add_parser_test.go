package bot

import "testing"

func TestParseAddCommand(t *testing.T) {
	const authToken = "secret"

	cases := []struct {
		name    string
		text    string
		wantURL string // empty means we expect empty result
	}{
		{
			name:    "title only",
			text:    "Buy milk",
			wantURL: "things:///add?auth-token=secret&title=Buy%20milk",
		},
		{
			name:    "title with when",
			text:    "Buy milk when:today",
			wantURL: "things:///add?auth-token=secret&title=Buy%20milk&when=today",
		},
		{
			name:    "title with tags",
			text:    "Buy milk tags:errands,personal",
			wantURL: "things:///add?auth-token=secret&tags=errands%2Cpersonal&title=Buy%20milk",
		},
		{
			name:    "title with deadline",
			text:    "Buy milk deadline:2026-12-31",
			wantURL: "things:///add?auth-token=secret&deadline=2026-12-31&title=Buy%20milk",
		},
		{
			name:    "title with unquoted notes",
			text:    "Buy milk notes:important",
			wantURL: "things:///add?auth-token=secret&notes=important&title=Buy%20milk",
		},
		{
			name:    "title with quoted notes",
			text:    `Buy milk notes:"pick up oat milk"`,
			wantURL: "things:///add?auth-token=secret&notes=pick%20up%20oat%20milk&title=Buy%20milk",
		},
		{
			name:    "title with all modifiers",
			text:    `Buy milk when:today deadline:2026-12-31 tags:errands notes:"any brand"`,
			wantURL: "things:///add?auth-token=secret&deadline=2026-12-31&notes=any%20brand&tags=errands&title=Buy%20milk&when=today",
		},
		{
			name:    "empty text",
			text:    "",
			wantURL: "",
		},
		{
			name:    "only whitespace",
			text:    "   ",
			wantURL: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseAddCommand(authToken, tc.text)
			if got != tc.wantURL {
				t.Errorf("parseAddCommand(%q)\n got  %q\n want %q", tc.text, got, tc.wantURL)
			}
		})
	}
}
