package bot

import "testing"

func TestParseAddCommand(t *testing.T) {
	cases := []struct {
		name         string
		text         string
		wantNil      bool
		wantTitle    string
		wantWhen     string
		wantDeadline string
		wantTags     []string
		wantNotes    string
	}{
		{
			name:      "title only",
			text:      "Buy milk",
			wantTitle: "Buy milk",
		},
		{
			name:      "title with when",
			text:      "Buy milk when:today",
			wantTitle: "Buy milk",
			wantWhen:  "today",
		},
		{
			name:      "title with tags",
			text:      "Buy milk tags:errands,personal",
			wantTitle: "Buy milk",
			wantTags:  []string{"errands", "personal"},
		},
		{
			name:         "title with deadline",
			text:         "Buy milk deadline:2026-12-31",
			wantTitle:    "Buy milk",
			wantDeadline: "2026-12-31",
		},
		{
			name:      "title with unquoted notes",
			text:      "Buy milk notes:important",
			wantTitle: "Buy milk",
			wantNotes: "important",
		},
		{
			name:      "title with quoted notes",
			text:      `Buy milk notes:"pick up oat milk"`,
			wantTitle: "Buy milk",
			wantNotes: "pick up oat milk",
		},
		{
			name:         "title with all modifiers",
			text:         `Buy milk when:today deadline:2026-12-31 tags:errands notes:"any brand"`,
			wantTitle:    "Buy milk",
			wantWhen:     "today",
			wantDeadline: "2026-12-31",
			wantTags:     []string{"errands"},
			wantNotes:    "any brand",
		},
		{
			name:    "empty text",
			text:    "",
			wantNil: true,
		},
		{
			name:    "only whitespace",
			text:    "   ",
			wantNil: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseAddCommandInput(tc.text)
			if tc.wantNil {
				if got != nil {
					t.Fatalf("expected nil result, got %#v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("expected non-nil result")
			}
			if got.title != tc.wantTitle || got.when != tc.wantWhen || got.deadline != tc.wantDeadline || got.notes != tc.wantNotes {
				t.Fatalf("unexpected parsed input: %#v", got)
			}
			if len(got.tags) != len(tc.wantTags) {
				t.Fatalf("unexpected tags count: got=%d want=%d", len(got.tags), len(tc.wantTags))
			}
			for i := range got.tags {
				if got.tags[i] != tc.wantTags[i] {
					t.Fatalf("unexpected tag at %d: got=%q want=%q", i, got.tags[i], tc.wantTags[i])
				}
			}
		})
	}
}
