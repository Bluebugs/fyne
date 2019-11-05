// +build !ci

// +build !android !ios !mobile
// +build js

package app

import (
	"net/url"
	"sync"
	"errors"

        "fyne.io/fyne"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/internal/driver/glfw"
)

// NewWithID returns a new app instance using the OpenGL driver.
// The ID string should be globally unique to this app.
func NewWithID(id string) fyne.App {
        return NewAppWithDriver(glfw.NewGLDriver(), id)
}

func (app *fyneApp) OpenURL(url *url.URL) error {
	return errors.New("OpenURL is not supported yet with GopherJS backend.")
}

func rootConfigDir() string {
	return "/data/"
}

// SettingsSchema is used for loading and storing global settings
type SettingsSchema struct {
        // these items are used for global settings load
        ThemeName string  `json:"theme"`
        Scale     float32 `json:"scale"`
}

type settings struct {
	themeLock sync.RWMutex
	theme     fyne.Theme

	listenerLock    sync.Mutex
	changeListeners []chan fyne.Settings

	schema SettingsSchema
}

// Declare conformity with Settings interface
var _ fyne.Settings = (*settings)(nil)

func (s *settings) Theme() fyne.Theme {
        s.themeLock.RLock()
        defer s.themeLock.RUnlock()
        return s.theme
}
func (s *settings) SetTheme(theme fyne.Theme) {
        s.themeLock.Lock()
        defer s.themeLock.Unlock()
        s.theme = theme
        s.apply()
}
func (s *settings) Scale() float32 {
        s.themeLock.RLock()
        defer s.themeLock.RUnlock()
        return s.schema.Scale
}

func (s *settings) AddChangeListener(listener chan fyne.Settings) {
        s.listenerLock.Lock()
        defer s.listenerLock.Unlock()
        s.changeListeners = append(s.changeListeners, listener)
}
func (s *settings) apply() {
        s.listenerLock.Lock()
        defer s.listenerLock.Unlock()
        for _, listener := range s.changeListeners {
                select {
                case listener <- s:
                default:
                        l := listener
                        go func() { l <- s }()
                }
        }
}

func (s *settings) load() {
	s.setupTheme()
}

func (s *settings) loadFromFile(path string) error {
	return nil
}

func (s *settings) setupTheme() {
        name := s.schema.ThemeName
        if name == "light" {
                s.SetTheme(theme.LightTheme())
        } else {
                s.SetTheme(theme.DarkTheme())
        }
}

func watchFile(path string, callback func()) {
}

func (s *settings) watchSettings() {
}

func (s *settings) stopWatching() {
}

func loadSettings() *settings {
        s := &settings{}
        s.load()
	return s
}
