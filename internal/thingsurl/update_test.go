package thingsurl_test

import (
	"testing"

	"github.com/IlyasYOY/telethings/internal/thingsurl"
)

func TestUpdate_Basic(t *testing.T) {
	got := thingsurl.New("mytoken").Update("SyJEz273ceSkabUbciM73A").String()

	want := "things:///update?auth-token=mytoken&id=SyJEz273ceSkabUbciM73A"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdate_Title(t *testing.T) {
	got := thingsurl.New("tok").Update("abc").WithTitle("Buy bread").String()

	want := "things:///update?auth-token=tok&id=abc&title=Buy%20bread"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdate_When(t *testing.T) {
	got := thingsurl.New("tok").Update("abc").WithWhen(thingsurl.WhenToday).String()

	want := "things:///update?auth-token=tok&id=abc&when=today"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdate_AppendNotes(t *testing.T) {
	got := thingsurl.New("tok").Update("abc").WithAppendNotes("Wholemeal bread").String()

	want := "things:///update?append-notes=Wholemeal%20bread&auth-token=tok&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdate_ClearDeadline(t *testing.T) {
	got := thingsurl.New("tok").Update("abc").WithDeadline("").String()

	want := "things:///update?auth-token=tok&deadline=&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdate_AppendChecklistItems(t *testing.T) {
	got := thingsurl.New("tok").Update("abc").WithAppendChecklistItems("Cheese", "Bread", "Eggplant").String()

	want := "things:///update?append-checklist-items=Cheese%0ABread%0AEggplant&auth-token=tok&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdate_Completed(t *testing.T) {
	got := thingsurl.New("tok").Update("abc").Completed().String()

	want := "things:///update?auth-token=tok&completed=true&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdate_Uncompleted(t *testing.T) {
	got := thingsurl.New("tok").Update("abc").Uncompleted().String()

	want := "things:///update?auth-token=tok&completed=false&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdate_Canceled(t *testing.T) {
	got := thingsurl.New("tok").Update("abc").Canceled().String()

	want := "things:///update?auth-token=tok&canceled=true&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdate_Duplicate(t *testing.T) {
	got := thingsurl.New("tok").Update("abc").Duplicate().String()

	want := "things:///update?auth-token=tok&duplicate=true&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdate_AddTags(t *testing.T) {
	got := thingsurl.New("tok").Update("abc").WithAddTags("Work", "Urgent").String()

	want := "things:///update?add-tags=Work%2CUrgent&auth-token=tok&id=abc"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestUpdate_Reveal(t *testing.T) {
	got := thingsurl.New("tok").Update("abc").Reveal().String()

	want := "things:///update?auth-token=tok&id=abc&reveal=true"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
