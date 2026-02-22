// Package opener provides an interface for opening URLs on the host OS.
package opener

import "os/exec"

type MacOSOpener struct{}

func (MacOSOpener) Open(url string) error {
	return exec.Command("open", url).Run()
}
