package thingsurl_test

import (
	"testing"

	"github.com/IlyasYOY/telethings/internal/thingsurl"
)

func TestVersion(t *testing.T) {
	got := thingsurl.New("tok").Version().String()

	want := "things:///version"
	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
