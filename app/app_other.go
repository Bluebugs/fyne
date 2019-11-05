// +build ci !linux,!darwin,!windows,!freebsd,!openbsd,!netbsd,!js

package app

import (
	"errors"
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func defaultTheme() fyne.Theme {
	return theme.DarkTheme()
}

func rootConfigDir() string {
	return "/tmp/fyne-test/"
}

func (app *fyneApp) OpenURL(url *url.URL) error {
	return errors.New("Unable to open url for unknown operating system")
}
