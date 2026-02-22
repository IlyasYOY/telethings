package thingsurl_test

import (
	"testing"

	"github.com/IlyasYOY/telethings/internal/thingsurl"
)

func TestUpdateProject_Basic(t *testing.T) {
	got := thingsurl.New("mytoken").UpdateProject("Jvj7EW1fLoScPhaw2JomCT").String()

	want := "things:///update-project?auth-token=mytoken&id=Jvj7EW1fLoScPhaw2JomCT"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdateProject_When(t *testing.T) {
	got := thingsurl.New("tok").UpdateProject("abc").WithWhen(thingsurl.WhenTomorrow).String()

	want := "things:///update-project?auth-token=tok&id=abc&when=tomorrow"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdateProject_AddTags(t *testing.T) {
	got := thingsurl.New("tok").UpdateProject("abc").WithAddTags("Important").String()

	want := "things:///update-project?add-tags=Important&auth-token=tok&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdateProject_PrependNotes(t *testing.T) {
	got := thingsurl.New("tok").UpdateProject("abc").WithPrependNotes("SFO to JFK.").String()

	want := "things:///update-project?auth-token=tok&id=abc&prepend-notes=SFO%20to%20JFK."
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdateProject_ClearDeadline(t *testing.T) {
	got := thingsurl.New("tok").UpdateProject("abc").WithDeadline("").String()

	want := "things:///update-project?auth-token=tok&deadline=&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdateProject_Area(t *testing.T) {
	got := thingsurl.New("tok").UpdateProject("abc").WithArea("Family").String()

	want := "things:///update-project?area=Family&auth-token=tok&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdateProject_Duplicate(t *testing.T) {
	got := thingsurl.New("tok").UpdateProject("abc").Duplicate().String()

	want := "things:///update-project?auth-token=tok&duplicate=true&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdateProject_Uncompleted(t *testing.T) {
	got := thingsurl.New("tok").UpdateProject("abc").Uncompleted().String()

	want := "things:///update-project?auth-token=tok&completed=false&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
