package thingsurl_test

import (
	"testing"

	"github.com/IlyasYOY/telethings/internal/thingsurl"
)

func TestSearch_WithQuery(t *testing.T) {
	got := thingsurl.New("tok").Search("vacation").String()

	want := "things:///search?query=vacation"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestSearch_NoQuery(t *testing.T) {
	got := thingsurl.New("tok").Search("").String()

	want := "things:///search"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
