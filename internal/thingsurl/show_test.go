package thingsurl_test

import (
	"testing"

	"github.com/IlyasYOY/telethings/internal/thingsurl"
)

func TestShow_BuiltInList(t *testing.T) {
	got := thingsurl.New("tok").Show(thingsurl.ShowListToday).String()

	want := "things:///show?id=today"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestShow_ID(t *testing.T) {
	got := thingsurl.New("tok").Show("GJJVZHE7SNu7xcVuH2xDDh").String()

	want := "things:///show?id=GJJVZHE7SNu7xcVuH2xDDh"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestShow_Query(t *testing.T) {
	got := thingsurl.New("tok").ShowByQuery("vacation").String()

	want := "things:///show?query=vacation"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestShow_QueryWithFilter(t *testing.T) {
	got := thingsurl.New("tok").ShowByQuery("vacation").WithFilter("Errand").String()

	want := "things:///show?filter=Errand&query=vacation"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestShow_AllBuiltInLists(t *testing.T) {
	for _, id := range []string{
		thingsurl.ShowListInbox,
		thingsurl.ShowListToday,
		thingsurl.ShowListAnytime,
		thingsurl.ShowListUpcoming,
		thingsurl.ShowListSomeday,
		thingsurl.ShowListLogbook,
		thingsurl.ShowListTomorrow,
		thingsurl.ShowListDeadlines,
		thingsurl.ShowListRepeating,
		thingsurl.ShowListAllProjects,
		thingsurl.ShowListLoggedProjects,
	} {
		got := thingsurl.New("tok").Show(id).String()
		if got == "" {
			t.Errorf("Show(%s) returned empty string", id)
		}
	}
}
