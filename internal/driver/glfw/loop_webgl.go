// +build js

package glfw

import (
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/painter"

	"github.com/goxjs/glfw"
	"github.com/goxjs/gl"
)

func (d *gLDriver) initGLFW() {
        err := glfw.Init(gl.ContextWatcher)
        if err != nil {
                fyne.LogError("failed to initialise GLFW", err)
                return
        }
}

func (d *gLDriver) runGL() {
	fps := time.NewTicker(time.Second / 60)
	runMutex.Lock()
	runFlag = true
	runMutex.Unlock()

	settingsChange := make(chan fyne.Settings)
	fyne.CurrentApp().Settings().AddChangeListener(settingsChange)
	d.initGLFW()

	for {
		select {
		case <-d.done:
			fps.Stop()
			glfw.Terminate()
			return
		case f := <-funcQueue:
			f.f()
			if f.done != nil {
				f.done <- true
			}
		case <-settingsChange:
			painter.ClearFontCache()
		case <-fps.C:
			glfw.PollEvents()
			newWindows := []fyne.Window{}
			reassign := false
			for _, win := range d.windows {
				w := win.(*window)
				viewport := w.viewport

				if viewport.ShouldClose() {
					reassign = true
					// remove window from window list
					viewport.Destroy()

					go w.destroy(d)
					continue
				} else {
					newWindows = append(newWindows, win)
				}

				canvas := w.canvas
				if !canvas.isDirty() || !w.visible {
					continue
				}

				d.repaintWindow(w)
			}
			if reassign {
				d.windows = newWindows
			}
		}
	}
}
