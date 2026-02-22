package thingsurl_test

import (
	"testing"

	"github.com/IlyasYOY/telethings/internal/thingsurl"
)

func TestAddProject_Title(t *testing.T) {
	got := thingsurl.New("tok").AddProject("Build treehouse").String()

	want := "things:///add-project?title=Build%20treehouse"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAddProject_When(t *testing.T) {
	got := thingsurl.New("tok").AddProject("Build treehouse").WithWhen(thingsurl.WhenToday).String()

	want := "things:///add-project?title=Build%20treehouse&when=today"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAddProject_Notes(t *testing.T) {
	got := thingsurl.New("tok").AddProject("Project").WithNotes("some notes").String()

	want := "things:///add-project?notes=some%20notes&title=Project"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAddProject_Deadline(t *testing.T) {
	got := thingsurl.New("tok").AddProject("Submit Tax").WithDeadline("2026-12-31").String()

	want := "things:///add-project?deadline=2026-12-31&title=Submit%20Tax"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAddProject_Tags(t *testing.T) {
	got := thingsurl.New("tok").AddProject("Project").WithTags("Work", "Urgent").String()

	want := "things:///add-project?tags=Work%2CUrgent&title=Project"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAddProject_Area(t *testing.T) {
	got := thingsurl.New("tok").AddProject("Plan Birthday Party").WithArea("Family").String()

	want := "things:///add-project?area=Family&title=Plan%20Birthday%20Party"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAddProject_AreaID(t *testing.T) {
	got := thingsurl.New("tok").AddProject("Submit Tax").WithAreaID("Lg8UqVPXo2SbJNiBpDBBQ").String()

	want := "things:///add-project?area-id=Lg8UqVPXo2SbJNiBpDBBQ&title=Submit%20Tax"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAddProject_Todos(t *testing.T) {
	got := thingsurl.New("tok").AddProject("Shopping").WithTodos("Milk", "Bread").String()

	want := "things:///add-project?title=Shopping&to-dos=Milk%0ABread"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAddProject_Completed(t *testing.T) {
	got := thingsurl.New("tok").AddProject("Done project").Completed().String()

	want := "things:///add-project?completed=true&title=Done%20project"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestAddProject_Reveal(t *testing.T) {
	got := thingsurl.New("tok").AddProject("Project").Reveal().String()

	want := "things:///add-project?reveal=true&title=Project"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
