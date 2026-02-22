package thingsurl_test

import (
	"testing"

	"github.com/IlyasYOY/telethings/internal/thingsurl"
)

func TestAdd_Title(t *testing.T) {
	add := thingsurl.New("auth-token").Add("this is test task")

	got := add.String()

	want := "things:///add?auth-token=auth-token&title=this%20is%20test%20task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_Completed(t *testing.T) {
	add := thingsurl.New("auth-token").Add("this is test task").Completed()

	got := add.String()

	want := "things:///add?auth-token=auth-token&completed=true&title=this%20is%20test%20task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_Canceled(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").Canceled()

	got := add.String()

	want := "things:///add?auth-token=auth-token&canceled=true&title=task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_Notes(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").WithNotes("buy low fat")

	got := add.String()

	want := "things:///add?auth-token=auth-token&notes=buy%20low%20fat&title=task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_When(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").WithWhen(thingsurl.WhenToday)

	got := add.String()

	want := "things:///add?auth-token=auth-token&title=task&when=today"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_Deadline(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").WithDeadline("2026-12-31")

	got := add.String()

	want := "things:///add?auth-token=auth-token&deadline=2026-12-31&title=task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_Tags(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").WithTags("Errand", "Home")

	got := add.String()

	want := "things:///add?auth-token=auth-token&tags=Errand%2CHome&title=task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_ChecklistItems(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").WithChecklistItems("Milk", "Bread")

	got := add.String()

	want := "things:///add?auth-token=auth-token&checklist-items=Milk%0ABread&title=task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_Titles(t *testing.T) {
	add := thingsurl.New("auth-token").Add("ignored").WithTitles("Milk", "Beer", "Cheese")

	got := add.String()

	want := "things:///add?auth-token=auth-token&titles=Milk%0ABeer%0ACheese"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_List(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").WithList("Shopping")

	got := add.String()

	want := "things:///add?auth-token=auth-token&list=Shopping&title=task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_ListID(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").WithListID("abc123")

	got := add.String()

	want := "things:///add?auth-token=auth-token&list-id=abc123&title=task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_Heading(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").WithList("Project").WithHeading("Phase 1")

	got := add.String()

	want := "things:///add?auth-token=auth-token&heading=Phase%201&list=Project&title=task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_Reveal(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").Reveal()

	got := add.String()

	want := "things:///add?auth-token=auth-token&reveal=true&title=task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_ShowQuickEntry(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").ShowQuickEntry()

	got := add.String()

	want := "things:///add?auth-token=auth-token&show-quick-entry=true&title=task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAdd_CreationDate(t *testing.T) {
	add := thingsurl.New("auth-token").Add("task").WithCreationDate("2026-01-01T00:00:00Z")

	got := add.String()

	want := "things:///add?auth-token=auth-token&creation-date=2026-01-01T00%3A00%3A00Z&title=task"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

